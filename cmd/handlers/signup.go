package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"pengoe/internal/logger"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/internal/utils"
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
	redirect, redirectFound := r.Context().Value("redirect").(string)
	if !redirectFound {
		router.InternalError(w, r, p)
		return errors.New("Should use redirect middleware")
	}
	db, dbFound := r.Context().Value("db").(*sql.DB)
	if !dbFound {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}

	formErr := r.ParseForm()
	if formErr != nil {
		router.InternalError(w, r, p)
		return formErr
	}

  formValues := r.Form

  values, err := utils.GetFormValues(
    formValues,
    "username",
    "email",
    "firstname",
    "lastname",
    "password",
  )
  if err != nil {
    router.BadRequest(w, r, p)
    return err
  }

  username := values["username"]
  email := values["email"]
  firstname := values["firstname"]
  lastname := values["lastname"]
  password := values["password"]

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

	_, signupErr := userService.Signup(newUser)

	// unsuccessful signup, render signup page with error message
	if signupErr != nil {
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

		if usernameError != "" {
			logger.Log(logger.WARNING, "signup/post", usernameError)
		}
		if emailError != "" {
			logger.Log(logger.WARNING, "signup/post", emailError)
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
