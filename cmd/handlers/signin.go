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
	"pengoe/web/templates/pages"
	"time"

	"github.com/a-h/templ"
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
SigninPageHandler handles the GET request to /signin.
*/
func SigninPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	params := utils.GetQueryParams(r)

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return dbErr
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
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return decodeErr
		}

		if redirect == "" {
			redirect = "/dashboard"
		}

		w.Header().Set("HX-Redirect", redirect)

		logger.Log(logger.INFO, "signin/get/decode", "Redirecting to "+redirect)

		return nil
	}

	// not logged in
	logger.Log(logger.INFO, "signin/checkSession", accessTokenErr.Error())

	logger.Log(logger.INFO, "Template parsed successfully", "signin/signinTmpl")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	redirect := params["redirect"]
	if redirect == "" {
		redirect = "%2Fdashboard"
	}

	data := pages.SigninProps{
		Title:       "pengoe - Sign in",
		Descrtipion: "Sign in to pengoe",
		Session: utils.Session{
			LoggedIn: false,
		},
		RedirectUrl: redirect,
		UsernameOrEmail: "",
		LoginError: "",
	}

	component := pages.Signin(data);
	handler := templ.Handler(component);
	handler.ServeHTTP(w, r);

	logger.Log(logger.INFO, "signin/get/res", "Template rendered successfully")
	return nil
}

/*
SigninHandler handles the POST request to /signin.
*/
func SigninHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {

	formErr := r.ParseForm()
	if formErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return formErr
	}

	logger.Log(logger.INFO, "signin/post/parse", "Form parsed successfully")

	usernameOrEmail := html.EscapeString(r.FormValue("user"))
	password := html.EscapeString(r.FormValue("password"))

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return dbErr
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// login the user
	user, loginErr := userService.Login(usernameOrEmail, password)

	// if the login was unsuccessful
	if loginErr != nil {
		logger.Log(logger.WARNING, "signin/post/login", loginErr.Error())

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)

		logger.Log(logger.INFO, "signin/post/errorTmpl", "Template parsed successfully")

		params := utils.GetQueryParams(r)

		redirect := params["redirect"]
		if redirect == "" {
			redirect = "%2Fdashboard"
		}

		data := pages.SigninProps{
			Title:       "pengoe - Sign in",
			Descrtipion: "Sign in to pengoe",
			Session: utils.Session{
				LoggedIn: false,
			},
			RedirectUrl: redirect,
			UsernameOrEmail: usernameOrEmail,
			LoginError: "Incorrect username or password",
		}

		component := pages.Signin(data);
		handler := templ.Handler(component);
		handler.ServeHTTP(w, r);

		logger.Log(logger.INFO, "signin/post/errorRes", "Template rendered successfully")
		return nil
	}

	userId := user.Id

	// if the login was successful
	accessToken, accessTokenErr := utils.NewToken(userId, utils.AccessToken)
	if accessTokenErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return accessTokenErr
	}

	logger.Log(logger.INFO, "signin/post/access", "Access token created successfully")

	refreshToken, refreshTokenErr := utils.NewToken(userId, utils.RefreshToken)
	if refreshTokenErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return refreshTokenErr
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return decodeErr
	}

	http.Redirect(w, r, redirect, http.StatusSeeOther)

	logger.Log(logger.INFO, "signin/post/res", "Redirected to "+redirect)
	return nil
}
