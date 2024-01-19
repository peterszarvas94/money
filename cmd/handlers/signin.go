package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/internal/token"
	"pengoe/internal/utils"
	"pengoe/web/templates/pages"
	"strconv"

	"github.com/a-h/templ"
)

type SigninPage struct {
	Title       string
	Descrtipion string
	LoggedIn    bool
	User        services.User
	Redirect    template.URL
	Values      map[string]string
	Errors      map[string]string
}

/*
SigninPageHandler handles the GET request to /signin.
*/
func SigninPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	redirect, redirectFound := r.Context().Value("redirect").(string)
	if !redirectFound {
		router.InternalError(w, r, p)
		return errors.New("Should use redirect middleware")
	}

	data := pages.SigninProps{
		Title:           "pengoe - Sign in",
		Descrtipion:     "Sign in to pengoe",
		RedirectUrl:     redirect,
		UsernameOrEmail: "",
		SigninErr:       "",
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
    "user",
		"password",
	)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	usernameOrEmail := values["user"]
	password := values["password"]

	userService := services.NewUserService(db)

	// login the user
	user, signinErr := userService.Signin(usernameOrEmail, password)

	// if the login was unsuccessful
	if signinErr != nil {
		w.WriteHeader(http.StatusUnauthorized)

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
		router.InternalError(w, r, p)
		return sessionErr
	}

	_, tokenErr := token.Manager.Create(session.Id)
	if tokenErr != nil {
		router.InternalError(w, r, p)
		return tokenErr
	}

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
		Name:     "session",
		Value:    sessionId,
		Path:     "/",
		Expires:  valid,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, redirect, http.StatusSeeOther)

	return nil
}
