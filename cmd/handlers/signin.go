package handlers

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"pengoe/web/templates/components"
	"pengoe/web/templates/layouts"
	"pengoe/web/templates/pages"
	"time"
)

type SigninPage struct {
	Title       string
	Descrtipion string
	Session     utils.Session
	Redirect    template.URL
	Values			map[string]string
	Errors			map[string]string
}

/*
getSigninTmpl helper function to parse the signin template.
*/
func getSigninTmpl() (*template.Template, error) {
	tmpl, tmplErr := template.ParseFiles(
		layouts.Base,
		pages.Signin,
		components.Icon,
		components.Error,
	)
	if tmplErr != nil {
		return nil, tmplErr
	}

	return tmpl, nil
}

/*
SigninPageHandler handles the GET request to /signin.
*/
func SigninPageHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	params := utils.GetQueryParams(r)

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "signin/get/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	user, accessTokenErr := userService.CheckAccessToken(r)

	// if the user is logged in
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

		w.Header().Set("HX-Redirect", redirect)

		logger.Log(logger.INFO, "signin/get/decode", "Redirecting to "+redirect)

		return
	}

	// not logged in
	logger.Log(logger.INFO, "signin/checkSession", accessTokenErr.Error())

	tmpl, tmplErr := getSigninTmpl()
	if tmplErr != nil {
		logger.Log(logger.ERROR, "signin/signinTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	logger.Log(logger.INFO, "Template parsed successfully", "signin/signinTmpl")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	redirect := params["redirect"]
	if redirect == "" {
		redirect = "%2Fdashboard"
	}

	data := SigninPage{
		Title:       "pengoe - Sign in",
		Descrtipion: "Sign in to pengoe",
		Session: utils.Session{
			LoggedIn: false,
		},
		Redirect: template.URL(redirect),
	}

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		logger.Log(logger.ERROR, "signin/get/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	logger.Log(logger.INFO, "signin/get/res", "Template rendered successfully")
	return
}

/*
SigninHandler handles the POST request to /signin.
*/
func SigninHandler(w http.ResponseWriter, r *http.Request, pattern string) {

	formErr := r.ParseForm()
	if formErr != nil {
		logger.Log(logger.ERROR, "signin/post/parse", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "signin/post/parse", "Form parsed successfully")

	usernameOrEmail := html.EscapeString(r.FormValue("user"))
	password := html.EscapeString(r.FormValue("password"))

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "signin/post/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// login the user
	userId, loginErr := userService.Login(usernameOrEmail, password)

	// if the login was unsuccessful
	if loginErr != nil {
		logger.Log(logger.ERROR, "signin/post/login", loginErr.Error())

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)

		tmpl, tmplErr := getSigninTmpl()
		if tmplErr != nil {
			logger.Log(logger.ERROR, "signin/post/errorTmpl", tmplErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		logger.Log(logger.INFO, "signin/post/errorTmpl", "Template parsed successfully")

		params := utils.GetQueryParams(r)

		redirect := params["redirect"]
		if redirect == "" {
			redirect = "%2Fdashboard"
		}

		data := SigninPage{
			Title:       "pengoe - Sign in",
			Descrtipion: "Sign in to pengoe",
			Session: utils.Session{
				LoggedIn: false,
			},
		Redirect: template.URL(redirect),
			Values: map[string]string{
				"usernameOrEmail": usernameOrEmail,
			},
			Errors: map[string]string{
				"loginError": "Invalid username or password",
			},
		}

		resErr := tmpl.Execute(w, data)
		if resErr != nil {
			logger.Log(logger.ERROR, "signin/post/errorRes", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		logger.Log(logger.INFO, "signin/post/errorRes", "Template rendered successfully")
		return
	}

	// if the login was successful
	accessToken, accessTokenErr := utils.NewToken(userId, utils.Access)
	if accessTokenErr != nil {
		logger.Log(logger.ERROR, "signin/post/access", accessTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "signin/post/access", "Access token created successfully")

	refreshToken, refreshTokenErr := utils.NewToken(userId, utils.Refresh)
	if refreshTokenErr != nil {
		logger.Log(logger.ERROR, "signin/post/refresh", refreshTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "signin/post/refresh", "Refresh token created successfully")

	expires := time.Unix(refreshToken.Expires, 0)

	// set the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    refreshToken.Token,
		Path:     "/refresh",
		Expires:  expires,
		HttpOnly: true,
		// Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// set the access token header
	w.Header().Set("Authorization", "Bearer "+accessToken.Token)

	params := utils.GetQueryParams(r)

	encoded := params["redirect"]
	redirect, decodeErr := url.QueryUnescape(encoded)
	if decodeErr != nil {
		logger.Log(logger.ERROR, "signin/post/decode", decodeErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirect, http.StatusSeeOther)

	logger.Log(logger.INFO, "signin/post/res", "Redirected to "+redirect)
	return
}
