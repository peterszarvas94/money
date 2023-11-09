package handlers

import (
	"fmt"
	"pengoe/db"
	"pengoe/services"
	"pengoe/types"
	"pengoe/utils"
	"html"
	"html/template"
	"net/http"
	"net/mail"
)

/*
getSignupTmpl helper function to parse the signup template.
*/
func getSignupTmpl() (*template.Template, error) {
	baseHtml := "templates/layouts/base.html"
	welcomeHtml := "templates/layouts/welcome.html"
	signupHtml := "templates/pages/signup.html"
	iconHtml := "templates/components/icon.html"
	errorHtml := "templates/components/error.html"
	incorrectHtml := "templates/components/incorrect.html"
	correctHtml := "templates/components/correct.html"

	tmpl, tmplErr := template.ParseFiles(
		baseHtml,
		welcomeHtml,
		signupHtml,
		iconHtml,
		errorHtml,
		incorrectHtml,
		correctHtml,
	)
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signup/signupTmpl", tmplErr.Error())
		return nil, tmplErr
	}

	utils.Log(utils.INFO, "signup/signupTmpl", "Template parsed successfully")
	return tmpl, nil
}

/*
SignupPageHandler handles the GET request to /signup.
*/
func SignupPageHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	// connet to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		utils.Log(utils.ERROR, "signup/get/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	user, sessionErr := userService.CheckAccessToken(r)
	if sessionErr == nil && user != nil {
		utils.Log(utils.INFO, "signup/checkSession", "Already logged in, redirecting to index")
		// http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.Log(utils.INFO, "signup/checkSession", "Not logged in")

	tmpl, tmplErr := getSignupTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signup/signupTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signup/signupTmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	pageData := types.Page{
		Title:       "pengoe - Sign up",
		Descrtipion: "Sign up to pengoe",
	}

	resErr := tmpl.Execute(w, pageData)
	if resErr != nil {
		utils.Log(utils.ERROR, "signup/get/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signup/get/res", "Template rendered successfully")
	return
}

/*
NewUserHandler handles the POST request to /signup.
*/
func NewUserHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	formErr := r.ParseForm()
	if formErr != nil {
		utils.Log(utils.ERROR, "signup/post/form", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "signup/post/form", "Form parsed successfully")

	username := html.EscapeString(r.FormValue("username"))
	email := html.EscapeString(r.FormValue("email"))
	firstname := html.EscapeString(r.FormValue("firstname"))
	lastname := html.EscapeString(r.FormValue("lastname"))
	password := html.EscapeString(r.FormValue("password"))

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		utils.Log(utils.ERROR, "signup/post/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// create user service
	userService := services.NewUserService(db)

	// add user
	newUser := &types.User{
		Username: username,
		Email:    email,
		Fistname: firstname,
		Lastname: lastname,
		Password: password,
	}

	userErr := userService.New(newUser)

	// unsuccessful signup, render signup page with error message
	if userErr != nil {
		utils.Log(utils.ERROR, "signup/post/userservice", userErr.Error())

		// check if email is valid
		_, invalid := mail.ParseAddress(email)
		if invalid != nil {
			logMsg := fmt.Sprintf("Invalid email: %s", email)
			utils.Log(utils.ERROR, "signup/post/emailinvalid", logMsg)
		}

		// check if username already exists
		_, usernameQueryErr := userService.GetByUsername(username)
		if usernameQueryErr != nil {
			utils.Log(utils.ERROR, "signup/post/usernamequery", usernameQueryErr.Error())
		}

		// check if email already exists
		_, emailQueryErr := userService.GetByEmail(email)
		if emailQueryErr != nil {
			utils.Log(utils.ERROR, "signup/post/emailquery", emailQueryErr.Error())
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
			utils.Log(utils.ERROR, "signup/post/usernameexists", logMsg)
			usernameError = "Username already exists"
			usernameCheck = "incorrect"
		} else {
			usernameCheck = "correct"
		}

		if emailExists {
			logMsg := fmt.Sprintf("Email already exists: %s", email)
			utils.Log(utils.ERROR, "signup/post/emailexists", logMsg)
			emailError = "Email already exists"
			emailCheck = "incorrect"
		} else if !emailInvalid {
			emailCheck = "correct"
		}

		signupData := types.Page{
			Session: types.Session{
				User: types.User{
					Username: username,
					Email:    email,
					Fistname: firstname,
					Lastname: lastname,
				},
			},
			Title:       "pengoe - Sign up",
			Descrtipion: "Sign up to pengoe",
			Data: map[string]string{
				"usernameValue": username,
				"usernameError": usernameError,
				"usernameCheck": usernameCheck,

				"emailValue": email,
				"emailError": emailError,
				"emailCheck": emailCheck,

				"firstnameValue": firstname,
				"lastnameValue":  lastname,
			},
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusConflict)

		tmpl, tmplErr := getSignupTmpl()
		if tmplErr != nil {
			utils.Log(utils.ERROR, "signup/post/signupTmpl", tmplErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "signup/post/signupTmpl", "Template parsed successfully")

		res_err := tmpl.Execute(w, signupData)
		if res_err != nil {
			utils.Log(utils.ERROR, "signup/post/res", res_err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "signup/post/res", "Template rendered successfully")
		return
	}

	// successful signup, redirect to signin
	utils.Log(utils.INFO, "signup/post/user", "User added successfully")
	// http.Redirect(w, r, "/signin", http.StatusSeeOther)
	return
}
