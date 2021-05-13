package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"picshare/conf"
	llctx "picshare/context"
	"picshare/controllers"
	"picshare/middleware"
	"picshare/models"
	"picshare/rand"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	boolPtr := flag.Bool("prod", false, "For production, set to true if providing a config file "+
		"before application initialization")
	flag.Parse()

	cfg := conf.LoadConfig(*boolPtr)
	dbConf := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbConf.PsqlConnInfo(), !cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
		models.WithOAuth(),
	)
	must(err)
	///*
	//// WARNING: Uncommenting this will destroy database
	////services.DestructiveReset()
	//*/
	services.AutoMigrate()

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	b, err := rand.Bytes(32)
	must(err)
	CSRF := csrf.Protect(b, csrf.Secure(cfg.IsProd()))
	userMW := middleware.User{
		UserService: services.User,
	}
	requireUserMW := middleware.RequireUser{
		User: userMW,
	}

	dbxOAuth := &oauth2.Config{
		ClientID:     cfg.DropBox.ID,
		ClientSecret: cfg.DropBox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.DropBox.AuthURL,
			TokenURL: cfg.DropBox.TokenURL,
		},
		RedirectURL: "http://localhost:3000/oauth/dropbox/callback",
	}

	dbxRedirect := func(w http.ResponseWriter, r *http.Request) {
		state := csrf.Token(r)
		cookie := http.Cookie{
			Name: "oauth_state",
			Value: state,
			HttpOnly: true,
			//Expires: time.Now().Local().Add(time.Minute * time.Duration(5)),
		}
		http.SetCookie(w, &cookie)
		url := dbxOAuth.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)
	}
	r.HandleFunc("/oauth/dropbox/connect", requireUserMW.ApplyFn(dbxRedirect))
	dbxCallback := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.FormValue("state")
		cookie, err := r.Cookie("oauth_state")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if cookie == nil || cookie.Value != state {
			http.Error(w, "Invalid state provided", http.StatusBadRequest)
			return
		}
		cookie.Value = ""
		cookie.Expires = time.Now()
		http.SetCookie(w, cookie)
		code := r.FormValue("code")
		token, err := dbxOAuth.Exchange(context.TODO(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		user := llctx.User(r.Context())
		oldToken, err := services.OAuth.Find(user.ID, models.OAuthDropbox)
		if err == models.ErrNotFound {
			//noop
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			services.OAuth.Delete(oldToken.ID)
		}
		userOAuth := models.OAuth{
			UserID: user.ID,
			Token: *token,
			Service: models.OAuthDropbox,
		}
		err = services.OAuth.Create(&userOAuth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%+v", token)
		//fmt.Fprintln(w, "code: ", r.FormValue("code"), "\nstate: ", r.FormValue("state"))
	}
	r.HandleFunc("/oauth/dropbox/callback", requireUserMW.ApplyFn(dbxCallback))

	dbxQuery := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		path := r.FormValue("path")

		user := llctx.User(r.Context())
		userOAuth, err := services.OAuth.Find(user.ID, models.OAuthDropbox)
		if err != nil {
			panic(err)
		}
		token := userOAuth.Token

		data := struct{
			Path string `json:"path"`
		}{
			Path: path,
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		client := dbxOAuth.Client(context.TODO(), &token)
		req, err := http.NewRequest(http.MethodPost, "https://api.dropboxapi.com/2/files/list_folder", bytes.NewReader(dataBytes))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	}
	r.HandleFunc("/oauth/dropbox/test", requireUserMW.ApplyFn(dbxQuery))

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/logout", requireUserMW.ApplyFn(usersC.Logout)).Methods("POST")
	r.Handle("/forgot", usersC.ForgotPwView).Methods("GET")
	r.HandleFunc("/forgot", usersC.InitiatePwReset).Methods("POST")
	r.HandleFunc("/reset", usersC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")

	// Assets
	assetHandler := http.FileServer(http.Dir("./assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetHandler))

	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// Gallery routes
	r.Handle("/galleries", requireUserMW.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMW.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMW.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMW.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMW.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMW.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMW.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMW.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGalleryName)

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), CSRF(userMW.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
