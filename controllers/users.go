package controllers

import (
	"net/http"
	"picapp/context"
	"picapp/email"
	"picapp/models"
	"picapp/rand"
	"picapp/views"
	"time"
)

// NewUsers creates a new Users controller. To be used during initial setup.
// If templates are incorrectly parsed, a panic will occur.
func NewUsers(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
		emailer:   emailer,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService

}

// New renders the form allowing users to create a new account
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

type SignupForm struct {
	Name     string `scheme:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create process the signup form after user submission
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	if err := u.signIn(w, &user); err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	u.emailer.Welcome(user.Name, user.Email)	//consider using a go routine here
	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcome to the site!",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login verifies the provided email-addy & password, logging in the user if correct.
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		u.NewView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Email address not found")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcome back!",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

// Logout will delete a users session cooki(remember_token) and
// update the user resource with a new remember token.
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	// Remove remember token adds a bit more security.
	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)

}

// signIn signs in a given user after account creation and sets cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

//// CookieTest displays the cookies set on the current user
//func (g *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
//	cookie, err := r.Cookie("remember_token")
//	if err != nil {
//		http.Redirect(w, r, "/login", http.StatusFound)
//		return
//	}
//	user, err := g.us.ByRemember(cookie.Value)
//	if err != nil {
//		http.Redirect(w, r, "/login", http.StatusFound)
//		return
//	}
//	fmt.Fprintln(w, user)
//}
