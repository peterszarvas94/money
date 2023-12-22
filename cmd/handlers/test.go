package handlers

import (
	"net/http"
	"pengoe/web/templates/pages"
	"github.com/a-h/templ"
)

func TestPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	path := r.URL.Path;

	data := pages.TestProps{
		Path: path,
		Varibales: p,
	}

	component := pages.Test(data);
	handler := templ.Handler(component);
	handler.ServeHTTP(w, r);

	return nil
}
