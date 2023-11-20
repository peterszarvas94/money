package router

import (
	"html/template"
	"net/http"
	"pengoe/internal/logger"
	"pengoe/web/templates/layouts"
	"pengoe/web/templates/pages"
)

/*
Notfound handles the 404 error.
*/
func Notfound(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles(layouts.Base, pages.NotFound)
	if tmplErr != nil {
		logger.Log(logger.ERROR, "notfound/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "notfound/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	resErr := tmpl.Execute(w, nil)
	if resErr != nil {
		logger.Log(logger.ERROR, "notfound/res", resErr.Error())
		http.Error(w, "Internal server  error", http.StatusInternalServerError)
	}

	logger.Log(logger.INFO, "notfound/res", "Template rendered successfully")
}

/*
MethodNotAllowed handles the 405 error.
*/
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles(layouts.Base, pages.MethodNotAllowed)
	if tmplErr != nil {
		logger.Log(logger.ERROR, "notallowed/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "notallowed/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
		
	resErr := tmpl.Execute(w, nil)
	if resErr != nil {
		logger.Log(logger.ERROR, "notallowed/res", resErr.Error())
		http.Error(w, "Internal server  error", http.StatusInternalServerError)
	}

	logger.Log(logger.INFO, "notallowed/res", "Template rendered successfully")
}
