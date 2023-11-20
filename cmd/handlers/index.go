package handlers

import (
	"html/template"
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"pengoe/web/templates/components"
	"pengoe/web/templates/layouts"
	"pengoe/web/templates/pages"
)

/*
Handler for home page "/".
*/
func HomePageHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	logger.Log(logger.INFO, "index/path", r.URL.Path)

	tmpl, tmplErr := template.ParseFiles(layouts.Base, pages.Index, components.Icon)
	if tmplErr != nil {
		logger.Log(logger.ERROR, "index/tmpl", tmplErr.Error())
		http.Error(w, "Intenal server error at tmpl", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "index/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "index/get/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	loggedIn := false

	user, sessionErr := userService.CheckAccessToken(r)
	if sessionErr == nil && user != nil {
		loggedIn = true
	}

	data := utils.Page{
		Title:       "pengoe - Home",
		Descrtipion: "Home page for pengoe",
		Session: utils.Session{
			LoggedIn: loggedIn,
		},
	}

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		logger.Log(logger.ERROR, "index/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "index/res", "Template rendered successfully")
}
