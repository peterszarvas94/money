package handlers

import (
	"net/http"
	"pengoe/internal/logger"
	"pengoe/web/templates/pages"
	"github.com/a-h/templ"
)

/*
Handler for home page "/".
*/
func HomePageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	logger.Log(logger.INFO, "index/path", r.URL.Path)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := pages.IndexProps{
		Title:       "pengoe - Home",
		Description: "Home page for pengoe",
	}

	component := pages.Index(data);
	handler := templ.Handler(component);
	handler.ServeHTTP(w, r);

	logger.Log(logger.INFO, "index/res", "Template rendered successfully")

	return nil
}
