package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"net/http"
	"picshare/conf"
	"picshare/controllers"
	"picshare/middleware"
	"picshare/models"
	"picshare/rand"
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

	oauthConfigs := make(map[string]*oauth2.Config)
	oauthConfigs[models.OAuthDropbox] = &oauth2.Config{
		ClientID:     cfg.DropBox.ID,
		ClientSecret: cfg.DropBox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.DropBox.AuthURL,
			TokenURL: cfg.DropBox.TokenURL,
		},
		RedirectURL: "http://localhost:3000/oauth/dropbox/callback",
	}
	oauthC := controllers.NewOAuth(services.OAuth, oauthConfigs)

	b, err := rand.Bytes(32)
	must(err)
	CSRF := csrf.Protect(b, csrf.Secure(cfg.IsProd()))
	userMW := middleware.User{
		UserService: services.User,
	}
	requireUserMW := middleware.RequireUser{
		User: userMW,
	}

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
	r.HandleFunc("/galleries/{id:[0-9]+}/images/link", requireUserMW.ApplyFn(galleriesC.ImageUploadLink)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMW.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMW.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGalleryName)

	// OAuth
	r.HandleFunc("/oauth/{service:[a-z]+}/connect", requireUserMW.ApplyFn(oauthC.Connect))
	r.HandleFunc("/oauth/{service:[a-z]+}/callback", requireUserMW.ApplyFn(oauthC.Callback))
	r.HandleFunc("/oauth/{service:[a-z]+}/test", requireUserMW.ApplyFn(oauthC.DropboxTest))

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), CSRF(userMW.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
