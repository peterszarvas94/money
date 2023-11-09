package handlers

import (
	"pengoe/db"
	"pengoe/services"
	"pengoe/types"
	"pengoe/utils"
	"html/template"
	"net/http"
)

/*
Handler for home page "/".
*/
func HomePageHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	utils.Log(utils.INFO, "index/path", r.URL.Path)

	baseHtml := "templates/layouts/base.html"
	indexHtml := "templates/pages/index.html"
	iconHtml := "templates/components/icon.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, indexHtml, iconHtml)
	if tmplErr != nil {
		utils.Log(utils.ERROR, "index/tmpl", tmplErr.Error())
		http.Error(w, "Intenal server error at tmpl", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "index/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		utils.Log(utils.ERROR, "signup/get/db", dbErr.Error())
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

	data := types.Page{
		Title:       "pengoe - Home",
		Descrtipion: "Home page for pengoe",
		Session: types.Session{
			LoggedIn: loggedIn,
		},
	}

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		utils.Log(utils.ERROR, "index/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "index/res", "Template rendered successfully")
}
