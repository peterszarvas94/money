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
func SignupPageHandler(w http.ResponseWriter, r *http.Request) {

	params := utils.GetQueryParams(r)

	// connet to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "signup/get/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// check if the user is already logged in, redirect to dashboard
	user, accessTokenErr := userService.CheckAccessToken(r)
	if user != nil {
		logMsg := fmt.Sprintf("Logged in as %d, redirecting to dashboard", user.Id)
		logger.Log(logger.INFO, "signin/checkSession", logMsg)

		//decode uri componetns
		encoded := params["redirect"]
		redirect, decodeErr := url.QueryUnescape(encoded)
		if decodeErr != nil {
			logger.Log(logger.ERROR, "signin/get/decode", decodeErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if redirect == "" {
			redirect = "/dashboard"
		}

		logger.Log(logger.INFO, "signin/get/decode", "Redirect to "+redirect)

		// todo: change this to http... with some extension?
		w.Header().Set("HX-Redirect", redirect)
		return
	}

	logger.Log(logger.INFO, "signup/checkSession", accessTokenErr.Error())

	logger.Log(logger.INFO, "signup/signupTmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	redirect := params["redirect"]
	if redirect == "" {
		redirect = "%2Fdashboard"
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

	logger.Log(logger.INFO, "signup/get/res", "Template rendered successfully")
	return
}

/*
NewUserHandler handles the POST request to /signup.
*/
func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	params := utils.GetQueryParams(r)

	formErr := r.ParseForm()
	if formErr != nil {
		logger.Log(logger.ERROR, "signup/post/form", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
		logger.Log(logger.ERROR, "signup/post/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
		logger.Log(logger.ERROR, "signup/post/userservice", userErr.Error())

		// check if email is valid
		_, invalid := mail.ParseAddress(email)
		if invalid != nil {
			logMsg := fmt.Sprintf("Invalid email: %s", email)
			logger.Log(logger.ERROR, "signup/post/emailinvalid", logMsg)
		}

		// check if username already exists
		_, usernameQueryErr := userService.GetByUsername(username)
		if usernameQueryErr != nil {
			logger.Log(logger.ERROR, "signup/post/usernamequery", usernameQueryErr.Error())
		}

		// check if email already exists
		_, emailQueryErr := userService.GetByEmail(email)
		if emailQueryErr != nil {
			logger.Log(logger.ERROR, "signup/post/emailquery", emailQueryErr.Error())
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
			logger.Log(logger.ERROR, "signup/post/usernameexists", logMsg)
			usernameError = "Username already exists"
			usernameCheck = "incorrect"
		} else {
			usernameCheck = "correct"
		}

		if emailExists {
			logMsg := fmt.Sprintf("Email already exists: %s", email)
			logger.Log(logger.ERROR, "signup/post/emailexists", logMsg)
			emailError = "Email already exists"
			emailCheck = "incorrect"
		} else if !emailInvalid {
			emailCheck = "correct"
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusConflict)

		logger.Log(logger.INFO, "signup/post/signupTmpl", "Template parsed successfully")

		redirect := params["redirect"]
		if redirect == "" {
			redirect = "%2Fdashboard"
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

		logger.Log(logger.INFO, "signup/post/res", "Template rendered successfully")
		return
	}

	// successful signup, redirect to signin
	logger.Log(logger.INFO, "signup/post/user", "User added successfully, redirect to "+params["redirect"])

	http.Redirect(w, r, "/signin?redirect="+params["redirect"], http.StatusSeeOther)
	return
}
