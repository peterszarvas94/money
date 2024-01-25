package router

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"pengoe/web/templates/pages"

	"github.com/a-h/templ"
)

/*
NotFound handles the 404 error.
*/
func NotFound(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	w.WriteHeader(http.StatusNotFound)

	component := pages.NotFound()
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return errors.New("Page not found")
}

/*
MethodNotAllowed handles the 405 error.
*/
func MethodNotAllowed(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	w.WriteHeader(http.StatusMethodNotAllowed)

	component := pages.NotAllowed()
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return errors.New("Method not allowed")
}

/*
InternalError handles the 500 error.
*/
func InternalError(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	w.WriteHeader(http.StatusInternalServerError)

	component := pages.InternalError()
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return errors.New("Internal server error")
}

/*
Unauthorized handles the 401 error.
*/
func Unauthorized(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	w.WriteHeader(http.StatusUnauthorized)

	return errors.New("Unauthorized")
}

/*
RedirectToSignin to signin
*/
func RedirectToSignin(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	path := r.URL.Path
	escapedPath := url.QueryEscape(path)

	http.Redirect(
		w,
		r,
		fmt.Sprintf("/signin?redirect=%s", escapedPath),
		http.StatusSeeOther,
	)

	return nil
}

/*
Bad Request
*/
func BadRequest(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	w.WriteHeader(http.StatusBadRequest)

	return errors.New("Bad request")
}
