package handlers

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/router"
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
	Values      map[string]string
	Errors      map[string]string
}

/*
SigninPageHandler handles the GET request to /signin.
*/
func SigninPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	redirect := utils.GetQueryParam(r.URL.Query(), "redirect")

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		router.InternalError(w, r)
		return dbErr
	}
	defer db.Close()
	userService := services.NewUserService(db)

	user, accessTokenErr := userService.CheckAccessToken(r)

	// if the user is logged in
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

		w.Header().Set("HX-Redirect", redirect)

		logger.Log(logger.INFO, "signin/get/decode", fmt.Sprintf("Redirect to %s", redirect))

		return nil
	}

	// not logged in
	logger.Log(logger.INFO, "signin/checkSession", accessTokenErr.Error())

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

	data := pages.SigninProps{
		Title:       "pengoe - Sign in",
		Descrtipion: "Sign in to pengoe",
		Session: utils.Session{
			LoggedIn: false,
		},
		RedirectUrl:     redirect,
		UsernameOrEmail: "",
		LoginError:      "",
	}

	component := pages.Signin(data)
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return nil
}

/*
SigninHandler handles the POST request to /signin.
*/
func SigninHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	redirect := utils.GetQueryParam(r.URL.Query(), "redirect")

	formErr := r.ParseForm()
	if formErr != nil {
		router.InternalError(w, r)
		return formErr
	}

	logger.Log(logger.INFO, "signin/post/parse", "Form parsed successfully")

	usernameOrEmail := html.EscapeString(r.FormValue("user"))
	password := html.EscapeString(r.FormValue("password"))

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		router.InternalError(w, r)
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

		if redirect == "" {
			redirect = "%2Fdashboard"
		} else {
			redirect = url.QueryEscape(redirect)
		}

		if !utils.IsValidRedirect(redirect, true) {
			router.InternalError(w, r)
			return nil
		}

		data := pages.SigninProps{
			Title:       "pengoe - Sign in",
			Descrtipion: "Sign in to pengoe",
			Session: utils.Session{
				LoggedIn: false,
			},
			RedirectUrl:     redirect,
			UsernameOrEmail: usernameOrEmail,
			LoginError:      "Incorrect username or password",
		}

		component := pages.Signin(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		return nil
	}

	// if the login was successful
	userId := user.Id

	accessToken, accessTokenErr := utils.NewToken(userId, utils.AccessToken)
	if accessTokenErr != nil {
		router.InternalError(w, r)
		return accessTokenErr
	}

	logger.Log(logger.INFO, "signin/post/access", "Access token created successfully")

	refreshToken, refreshTokenErr := utils.NewToken(userId, utils.RefreshToken)
	if refreshTokenErr != nil {
		router.InternalError(w, r)
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
		// TODO: uncomment this when https is enabled
		// Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// set the access token header
	w.Header().Set("Authorization", "Bearer "+accessToken.Token)

	if redirect == "" {
		redirect = "/dashboard"
	}

	if !utils.IsValidRedirect(redirect, false) {
		router.InternalError(w, r)
		return nil
	}

	http.Redirect(w, r, redirect, http.StatusSeeOther)

	logger.Log(logger.INFO, "signin/post/res", fmt.Sprintf("Redirected to %s", redirect))
	return nil
}
