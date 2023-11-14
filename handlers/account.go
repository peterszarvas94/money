package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"pengoe/db"
	"pengoe/services"
	"pengoe/types"
	"pengoe/utils"
)

type NewAccountPage struct {
	Title       string
	Descrtipion string
	Session     types.Session
}

/*
getNewAccountTmpl helper function to parse the newaccount template.
*/
func getNewAccountTmpl() (*template.Template, error) {
	baseHtml := "templates/layouts/base.html"
	newaccountHtml := "templates/pages/newaccount.html"
	spinnerHtml := "templates/components/spinner.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, newaccountHtml, spinnerHtml)
	if tmplErr != nil {
		return nil, tmplErr
	}

	return tmpl, nil
}

/*
NewAccountPageHandler handles the GET request to /account/new
*/
func NewAccountPageHandler(w http.ResponseWriter, r *http.Request, pattern string) {

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		utils.Log(utils.ERROR, "newaccount/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	userService := services.NewUserService(db)

	// check if the user is logged in, protected route
	user, sessionErr := userService.CheckAccessToken(r)
	if user != nil {
		// logged in user
		logMsg := fmt.Sprintf("Logged in as %d", user.Id)
		utils.Log(utils.INFO, "newaccount/checkSession", logMsg)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := NewAccountPage{
			Title:       "pengoe - New Account",
			Descrtipion: "New Account for pengoe",
			Session: types.Session{
				LoggedIn: true,
				User:     *user,
			},
		}

		tmpl, tmplErr := getNewAccountTmpl()
		if tmplErr != nil {
			utils.Log(utils.ERROR, "newaccount/loggedin/tmpl", tmplErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "newaccount/loggedin/tmpl", "Template parsed successfully")

		resErr := tmpl.Execute(w, data)
		if resErr != nil {
			utils.Log(utils.ERROR, "newaccount/loggedin/res", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "newaccount/loggedin/res", "Template rendered successfully")
		return
	}

	// not logged in user
	utils.Log(utils.INFO, "newaccount/checkSession", sessionErr.Error())

	tmpl, tmplErr := getNewAccountTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "newaccount/notloggedin/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "newaccount/notloggedin/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := types.Page{
		Title:       "pengoe - New Account",
		Descrtipion: "New Account for pengoe",
		Session: types.Session{
			LoggedIn: false,
		},
	}

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		utils.Log(utils.ERROR, "newaccount/notloggedin/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "newaccount/notloggedin/res", "Template rendered successfully")
	return
}