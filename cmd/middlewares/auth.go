package middlewares

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"pengoe/internal/router"
	"pengoe/internal/services"
	t "pengoe/internal/token"
	"pengoe/internal/utils"
)

/*
AuthPage checks if the user is already logged in.
If logged in, redirects to "redirect" query param.
Otherwise, set context value "redirect" to "redirect" query param.
Used for signup and signin pages.
*/
func AuthPage(next router.HandlerFunc) router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p map[string]string) error {
		redirect := utils.GetQueryParam(r.URL.Query(), "redirect")

		if redirect == "" {
			redirect = "/dashboard"
		}

		if !utils.IsValidRedirect(redirect, false) {
			router.InternalError(w, r, p)
			return nil
		}

		token, tokenErr := t.GetSessionFromCookie(r)
		if tokenErr == nil || token != nil {
			http.Redirect(w, r, redirect, http.StatusSeeOther)
			return nil
		}

		redirect = url.QueryEscape(redirect)

		ctx := context.WithValue(r.Context(), "redirect", redirect)
		r = r.WithContext(ctx)

		return next(w, r, p)
	}
}

/*
Token checks if the user is logged in.
If not logged in, redirects to /signin.
Otherwise, set context value "token" to the token.
*/
func Token(next router.HandlerFunc) router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p map[string]string) error {
		token, tokenErr := t.GetSessionFromCookie(r)
		if tokenErr != nil {
			if r.Method == http.MethodGet {
				router.RedirectToSignin(w, r, p)
			} else {
				router.Unauthorized(w, r, p)
			}
			return tokenErr
		}

		ctx := context.WithValue(r.Context(), "token", token)
		r = r.WithContext(ctx)

		return next(w, r, p)
	}
}

/*
Session injects the session into the request context.
It needs WithToken and WithDB to be called before.
*/
func Session(next router.HandlerFunc) router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p map[string]string) error {
		token, tokenFound := r.Context().Value("token").(*t.Token)
		if !tokenFound {
			return router.RedirectToSignin(w, r, p)
		}

		db, tokenFound := r.Context().Value("db").(*sql.DB)
		if !tokenFound {
			router.InternalError(w, r, p)
		}

		sessionService := services.NewSessionService(db)

		session, sessionErr := sessionService.GetById(token.SessionID)
		if sessionErr != nil {
			router.InternalError(w, r, p)
			return sessionErr
		}

		ctx := context.WithValue(r.Context(), "session", session)
		r = r.WithContext(ctx)

		return next(w, r, p)
	}
}
