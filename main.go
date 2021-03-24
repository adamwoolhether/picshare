package main

import (
	"fmt"
	"net/http"
	"picapp/controllers"
	"picapp/middleware"
	"picapp/models"

	"github.com/gorilla/mux"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "adam"
	dbname = "picapp"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	services, err := models.NewServices(psqlInfo)
	must(err)
	//services.DestructiveReset()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)
	requireUserMW := middleware.RequireUser{services.User}

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	// Gallery routes
	r.Handle("/galleries/new", requireUserMW.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMW.ApplyFn(galleriesC.Create)).Methods("POST")


	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe("localhost:3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
