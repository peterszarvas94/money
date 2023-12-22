package handlers

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/mail"
	"net/url"
	"pengoe/internal/db"
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
	Session     utils.Session
	Redirect    template.URL
	Values      map[string]string
	Errors      map[string]string
}

/*
SignupPageHandler handles the GET request to /signup.
*/
func SignupPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	redirect := utils.GetQueryParam(r.URL.Query(), "redirect")

	// connet to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		router.InternalError(w, r)
		return dbErr
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// check if the user is already logged in, redirect to dashboard
	user, accessTokenErr := userService.CheckAccessToken(r)
	if user != nil {
		logMsg := fmt.Sprintf("Logged in as %d, redirecting to dashboard", user.Id)
		logger.Log(logger.INFO, "signin/checkSession", logMsg)

		if redirect == "" {
			redirect = "/dashboard"
		}

		if !utils.IsValidRedirect(redirect, false) {
			router.InternalError(w, r)
			return nil
		}

		logger.Log(logger.INFO, "signin/get/decode", fmt.Sprintf("Redirecting to %s", redirect))

		w.Header().Set("HX-Redirect", redirect)
		return nil
	}

	logger.Log(logger.INFO, "signup/checkSession", accessTokenErr.Error())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if redirect == "" {
		redirect = "%2Fdashboard"
	} else {
		redirect = url.QueryEscape(redirect)
	}

	if !utils.IsValidRedirect(redirect, true) {
		router.InternalError(w, r)
		return nil
	}

	data := pages.SignupProps{
		Title:       "pengoe - Sign in",
		Description: "Sign in to pengoe",
		Session: utils.Session{
			LoggedIn: false,
		},
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
NewUserHandler handles the POST request to /signup.
*/
func NewUserHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	redirect := utils.GetQueryParam(r.URL.Query(), "redirect")

	formErr := r.ParseForm()
	if formErr != nil {
		router.InternalError(w, r)
		return formErr
	}

	logger.Log(logger.INFO, "signup/post/form", "Form parsed successfully")

	username := html.EscapeString(r.FormValue("username"))
	email := html.EscapeString(r.FormValue("email"))
	firstname := html.EscapeString(r.FormValue("firstname"))
	lastname := html.EscapeString(r.FormValue("lastname"))
	password := html.EscapeString(r.FormValue("password"))

	// TODO: handle empty values

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		router.InternalError(w, r)
		return dbErr
	}
	defer db.Close()

	// create user service
	userService := services.NewUserService(db)

	// add user
	newUser := &utils.User{
		Username: username,
		Email:    email,
		Fistname: firstname,
		Lastname: lastname,
		Password: password,
	}

	_, userErr := userService.New(newUser)

	// unsuccessful signup, render signup page with error message
	if userErr != nil {
		logger.Log(logger.WARNING, "signup/post/userservice", userErr.Error())

		// check if email is valid
		_, invalid := mail.ParseAddress(email)
		if invalid != nil {
			logMsg := fmt.Sprintf("Invalid email: %s", email)
			logger.Log(logger.WARNING, "signup/post/emailinvalid", logMsg)
		}

		// check if username already exists
		_, usernameQueryErr := userService.GetByUsername(username)
		if usernameQueryErr != nil {
			logger.Log(logger.WARNING, "signup/post/usernamequery", usernameQueryErr.Error())
		}

		// check if email already exists
		_, emailQueryErr := userService.GetByEmail(email)
		if emailQueryErr != nil {
			logger.Log(logger.WARNING, "signup/post/emailquery", emailQueryErr.Error())
		}

		emailInvalid := invalid != nil
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
			logMsg := fmt.Sprintf("Username already exists: %s", username)
			logger.Log(logger.WARNING, "signup/post/usernameexists", logMsg)
			usernameError = "Username already exists"
			usernameCheck = "incorrect"
		} else {
			usernameCheck = "correct"
		}

		if emailExists {
			logMsg := fmt.Sprintf("Email already exists: %s", email)
			logger.Log(logger.WARNING, "signup/post/emailexists", logMsg)
			emailError = "Email already exists"
			emailCheck = "incorrect"
		} else if !emailInvalid {
			emailCheck = "correct"
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusConflict)

		logger.Log(logger.INFO, "signup/post/signupTmpl", "Template parsed successfully")

		if redirect == "" {
			redirect = "%2Fdashboard"
		} else {
			redirect = url.QueryEscape(redirect)
		}

		if !utils.IsValidRedirect(redirect, true) {
			router.InternalError(w, r)
			return nil
		}

		data := pages.SignupProps{
			Title:       "pengoe - Sign in",
			Description: "Sign in to pengoe",
			Session: utils.Session{
				LoggedIn: false,
			},
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

	if redirect == "" {
		redirect = "%2Fdashboard"
	} else {
		redirect = url.QueryEscape(redirect)
	}

	if !utils.IsValidRedirect(redirect, true) {
		router.InternalError(w, r)
		return nil
	}

	// successful signup, redirect to signin
	logger.Log(logger.INFO, "signup/post/user", fmt.Sprintf("User added successfully, redirect to /signin?redirect=%s", redirect))

	fmt.Printf("redirect: %s\n", redirect)

	http.Redirect(w, r, fmt.Sprintf("/signin?redirect=%s", redirect), http.StatusSeeOther)
	return nil
}
