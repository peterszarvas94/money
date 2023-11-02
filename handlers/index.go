package handlers

import (
	"fmt"
	"go-htmx/utils"
	"html/template"
	"net/http"
)

/*
Handler for home page "/".
*/
func HomePageHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	utils.Log(utils.INFO, "index/path", r.URL.Path)

	session := utils.CheckSession(r)

	if session.LoggedIn {
		sessionId := fmt.Sprint(session.User.Id)
		sessionLog := fmt.Sprintf("Session found: %s", sessionId)
		utils.Log(utils.INFO, "index/checkSession", sessionLog)
	} else {
		utils.Log(utils.INFO, "index/checkSession", "No session")
	}

	baseHtml := "templates/base.html"
	titleHtml := "templates/title.html"
	indexHtml := "templates/index.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, titleHtml, indexHtml)
	if tmplErr != nil {
		utils.Log(utils.ERROR, "index/tmpl", tmplErr.Error())
		http.Error(w, "Intenal server error at tmpl", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "index/tmpl", "Template parsed successfully")

	data := utils.PageData{
		Session: session,
		Title: "pengoe",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		utils.Log(utils.ERROR, "index/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "index/res", "Template rendered successfully")
}
