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
	"pengoe/internal/serversession"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"pengoe/web/templates/pages"
	"strconv"

	"github.com/a-h/templ"
)

type SigninPage struct {
	Title       string
	Descrtipion string
	Session     utils.CurrentUser
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
	sessionService := services.NewSessionService(db)

	session, sessionErr := sessionService.CheckCookie(r)

	// if the user is logged in
	if session != nil {
		logger.Log(
			logger.INFO,
			"signin/session",
			fmt.Sprintf("Session found with ID %d", session.Id),
		)

		if redirect == "" {
			redirect = "/dashboard"
		}

		if !utils.IsValidRedirect(redirect, false) {
			logger.Log(
				logger.ERROR,
				"signin/session/redirect",
				fmt.Sprintf("Invalid redirect %s", redirect),
			)
			router.InternalError(w, r)
			return nil
		}

		http.Redirect(w, r, redirect, http.StatusSeeOther)

		logger.Log(
			logger.INFO,
			"signin/session/redirect",
			fmt.Sprintf("Redirecting to %s", redirect),
		)

		return nil
	}

	// not logged in
	logger.Log(
		logger.INFO,
		"signin/nosession",
		fmt.Sprintf("No session found, %s", sessionErr.Error()),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if redirect == "" {
		redirect = "%2Fdashboard"
	} else {
		redirect = url.QueryEscape(redirect)
	}

	if !utils.IsValidRedirect(redirect, true) {
		logger.Log(
			logger.ERROR,
			"signin/nosession/redirect",
			fmt.Sprintf("Invalid redirect %s", redirect),
		)
		router.InternalError(w, r)
		return nil
	}

	data := pages.SigninProps{
		Title:           "pengoe - Sign in",
		Descrtipion:     "Sign in to pengoe",
		RedirectUrl:     redirect,
		UsernameOrEmail: "",
		SigninErr:       "",
	}

	logger.Log(
		logger.INFO,
		"signin/nosession/render",
		fmt.Sprintf("Rendering signin page with redirect %s", redirect),
	)

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
	user, signinErr := userService.Signin(usernameOrEmail, password)

	// if the login was unsuccessful
	if signinErr != nil {
		logger.Log(logger.WARNING, "signin/post/login", signinErr.Error())

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
			Title:           "pengoe - Sign in",
			Descrtipion:     "Sign in to pengoe",
			RedirectUrl:     redirect,
			UsernameOrEmail: usernameOrEmail,
			SigninErr:       "Incorrect username or password",
		}

		component := pages.Signin(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		return nil
	}

	// if the login was successful
	sessionService := services.NewSessionService(db)
	session, sessionErr := sessionService.New(user)
	if sessionErr != nil {
		logger.Log(logger.ERROR, "signin/post/session", sessionErr.Error())
		router.InternalError(w, r)
		return sessionErr
	}

	logger.Log(logger.INFO, "signin/post/session", "Session created successfully")

	_, ssErr := serversession.Manager.Create(session.Id)
	if ssErr != nil {
		logger.Log(logger.ERROR, "signin/post/serversession", ssErr.Error())
		router.InternalError(w, r)
		return ssErr
	}

	logger.Log(
		logger.INFO,
		"signin/post/serversession",
		"Server session created successfully",
	)

	secure := utils.Env.Environment == "production"
	var sameSite http.SameSite
	if utils.Env.Environment == "production" {
		sameSite = http.SameSiteLaxMode
	} else {
		sameSite = http.SameSiteDefaultMode
	}

	sessionId := strconv.Itoa(session.Id)

	valid := session.ValidUntil.UTC()

	// set the session id to cookie
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionId,
		Path:     "/",
		Expires:  valid,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if redirect == "" {
		redirect = "/dashboard"
	}

	if !utils.IsValidRedirect(redirect, false) {
		router.InternalError(w, r)
		return nil
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, redirect, http.StatusSeeOther)

	logger.Log(
		logger.INFO,
		"signin/post/res",
		fmt.Sprintf("Redirected to %s", redirect),
	)

	return nil
}
