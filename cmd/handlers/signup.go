package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/mail"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/web/templates/pages"

	"github.com/a-h/templ"
)

type SignupPage struct {
	Title       string
	Descrtipion string
	Redirect    template.URL
	Values      map[string]string
	Errors      map[string]string
}

/*
SignupPageHandler handles the GET request to /signup.
*/
func SignupPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	redirect, found := r.Context().Value("redirect").(string)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use redirect middleware")
	}

	data := pages.SignupProps{
		Title:         "pengoe - Sign in",
		Description:   "Sign in to pengoe",
		RedirectUrl:   redirect,
		Username:      "",
		Email:         "",
		Firstname:     "",
		Lastname:      "",
		UsernameCheck: "",
		UsernameError: "",
		EmailCheck:    "",
		EmailError:    "",
	}

	component := pages.Signup(data)
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return nil
}

/*
SignupHandler handles the POST request to /signup.
*/
func SignupHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	redirect, found := r.Context().Value("redirect").(string)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use redirect middleware")
	}
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}

	err := r.ParseForm()
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

  form := r.Form

  username := html.EscapeString(form.Get("username"))
  if username == "" {
    router.BadRequest(w, r, p)
    return errors.New("Username is required")
  }

  email := html.EscapeString(form.Get("email"))
  if email == "" {
    router.BadRequest(w, r, p)
    return errors.New("Email is required")
  }

  firstname := html.EscapeString(form.Get("firstname"))
  if firstname == "" {
    router.BadRequest(w, r, p)
    return errors.New("Firstname is required")
  }

  lastname := html.EscapeString(form.Get("lastname"))
  if lastname == "" {
    router.BadRequest(w, r, p)
    return errors.New("Lastname is required")
  }

  password := html.EscapeString(form.Get("password"))
  if password == "" {
    router.BadRequest(w, r, p)
    return errors.New("Password is required")
  }

	// create user service
	userService := services.NewUserService(db)

	// add user
	newUser := &services.User{
		Username: username,
		Email:    email,
		Fistname: firstname,
		Lastname: lastname,
		Password: password,
	}

	_, err = userService.Signup(newUser)

	// unsuccessful signup, render signup page with error message
	if err != nil {
		_, parseErr := mail.ParseAddress(email)
		_, usernameQueryErr := userService.GetByUsername(username)
		_, emailQueryErr := userService.GetByEmail(email)

		emailInvalid := parseErr != nil
		usernameExists := usernameQueryErr == nil
		emailExists := emailQueryErr == nil

		usernameError := ""
		usernameCheck := ""
		emailError := ""
		emailCheck := ""

		if emailInvalid {
			emailError = "Invalid email"
			emailCheck = "incorrect"
		}

		if usernameExists {
			usernameError = "Username already exists"
			usernameCheck = "incorrect"
		} else {
			usernameCheck = "correct"
		}

		if emailExists {
			emailError = "Email already exists"
			emailCheck = "incorrect"
		} else if !emailInvalid {
			emailCheck = "correct"
		}

		w.WriteHeader(http.StatusConflict)

		data := pages.SignupProps{
			Title:         "pengoe - Sign in",
			Description:   "Sign in to pengoe",
			RedirectUrl:   redirect,
			Firstname:     firstname,
			Lastname:      lastname,
			Username:      username,
			UsernameCheck: usernameCheck,
			UsernameError: usernameError,
			Email:         email,
			EmailCheck:    emailCheck,
			EmailError:    emailError,
		}

		component := pages.Signup(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		return nil
	}

	// successful signup, redirect to signin page
	http.Redirect(
		w,
		r,
		fmt.Sprintf("/signin?redirect=%s", redirect),
		http.StatusSeeOther,
	)

	return nil
}
