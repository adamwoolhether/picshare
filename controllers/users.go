package controllers

import (
	"fmt"
	"net/http"
	"picapp/models"
	"picapp/views"
)

// NewUsers creates a new Users controller. To be used during initial setup.
// If templates are incorrectly parsed, a panic will occur.
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

//// New renders the form allowing users to create a new account
//// GET /signup
//func (u *Users) New(w http.ResponseWriter, r *http.Request) {
//	if err := u.NewView.Render(w, nil); err != nil {
//		panic(err)
//	}
//}

type SignupForm struct {
	Name     string `scheme:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create process the signup form after user submission
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	err := u.us.Create(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login verifies the provided email-addy & password, logging in the user if correct.
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err) // Implement better error handling for user
	}
	// Do something
	fmt.Fprintln(w, form)
}
