package router

import (
	"net/http"
	"pengoe/internal/logger"
	"pengoe/web/templates/pages"

	"github.com/a-h/templ"
)

/*
NotFound handles the 404 error.
*/
func NotFound(w http.ResponseWriter, r *http.Request) {
	logger.Log(logger.INFO, "notfound/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	component := pages.NotFound()
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)
}

/*
MethodNotAllowed handles the 405 error.
*/
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	logger.Log(logger.INFO, "notallowed/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusMethodNotAllowed)

	component := pages.NotAllowed()
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)
}

/*
InternalError handles the 500 error.
*/
func InternalError(w http.ResponseWriter, r *http.Request) {
	logger.Log(logger.INFO, "internalservererror/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	component := pages.InternalError()
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)
}
