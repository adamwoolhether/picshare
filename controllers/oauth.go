package controllers

import (
	"context"
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"net/http"
	llctx "picshare/context"
	"picshare/models"
	"time"
)

type OAuth struct {
	os      models.OAuthService
	configs map[string]*oauth2.Config
}

func NewOAuth(os models.OAuthService, configs map[string]*oauth2.Config) *OAuth {
	return &OAuth{
		os:      os,
		configs: configs,
	}
}

func (o *OAuth) Connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
	}

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
	token, err := oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := llctx.User(r.Context())
	oldToken, err := o.os.Find(user.ID, service)
	if err == models.ErrNotFound {
		//no-op
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		o.os.Delete(oldToken.ID)
	}
	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: service,
	}
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%+v", token)
}

func (o *OAuth) DropboxTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]

	r.ParseForm()
	path := r.FormValue("path")

	user := llctx.User(r.Context())
	userOAuth, err := o.os.Find(user.ID, service)
	if err != nil {
		panic(err)
	}
	token := userOAuth.Token

	config := dropbox.Config{
		Token:    token.AccessToken,
		LogLevel: dropbox.LogInfo,
	}
	dbx := files.New(config)
	res, err := dbx.ListFolder(&files.ListFolderArg{
		Path: path,
	})
	if err != nil {
		panic(err)
	}
	for _, entry := range res.Entries {
		switch meta := entry.(type) {
		case *files.FolderMetadata:
			fmt.Fprintln(w, "FolderMetadata=", meta)
		case *files.FileMetadata:
			fmt.Fprintln(w, "FileMetadata=", meta)

		}
	}
}
