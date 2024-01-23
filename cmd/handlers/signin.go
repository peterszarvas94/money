package handlers

import (
	"database/sql"
	"errors"
	"html"
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
	redirect, found := r.Context().Value("redirect").(string)
	if !found {
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

  usernameOrEmail := html.EscapeString(form.Get("user"))
  if usernameOrEmail == "" {
    router.BadRequest(w, r, p)
    return errors.New("Username or email is required")
  }

  password := html.EscapeString(form.Get("password"))
  if password == "" {
    router.BadRequest(w, r, p)
    return errors.New("Password is required")
  }

	userService := services.NewUserService(db)

	// login the user
	user, err := userService.Signin(usernameOrEmail, password)

	// if the login was unsuccessful
	if err != nil {
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
	session, err := sessionService.New(user)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	_, err = token.Manager.Create(session.Id)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	secure := utils.Env.ENVIRONMENT == "production"
	var sameSite http.SameSite
	if utils.Env.ENVIRONMENT == "production" {
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
