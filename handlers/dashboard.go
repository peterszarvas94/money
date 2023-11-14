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

type Account struct {
	Id       int
	Text     string
	Selected bool
	New      bool
}

type DashboardPage struct {
	Title       string
	Descrtipion string
	Session     types.Session
	Accounts    []Account
}

/*
getDashboardTmpl helper function to parse the dashboard template.
*/
func getDashboardTmpl() (*template.Template, error) {
	baseHtml := "templates/layouts/base.html"
	dashboardHtml := "templates/pages/dashboard.html"
	iconHtml := "templates/components/icon.html"
	accountHtml := "templates/components/account.html"
	spinnerHtml := "templates/components/spinner.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, dashboardHtml, iconHtml, accountHtml, spinnerHtml)
	if tmplErr != nil {
		return nil, tmplErr
	}

	return tmpl, nil
}

/*
DashboardPageHandler handles the GET request to /dashboard.
*/
func DashboardPageHandler(w http.ResponseWriter, r *http.Request, pattern string) {

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		utils.Log(utils.ERROR, "check/db", dbErr.Error())
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
		utils.Log(utils.INFO, "dashboard/checkSession", logMsg)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := DashboardPage{
			Title:       "pengoe - Dashboard",
			Descrtipion: "Dashboard for pengoe",
			Session: types.Session{
				LoggedIn: true,
				User:     *user,
			},
			Accounts: []Account{
				{
					Id:   1,
					Text: "Account 1",
				},
				{
					Id:       2,
					Text:     "Account 2",
					Selected: true,
				},
				{
					Id:   3,
					Text: "Account 3",
				},
				{
					Id:   0,
					Text: "New account",
					New:  true,
				},
			},
		}

		tmpl, tmplErr := getDashboardTmpl()
		if tmplErr != nil {
			utils.Log(utils.ERROR, "dashboard/loggedin/tmpl", tmplErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "dashboard/tmpl", "Template parsed successfully")

		resErr := tmpl.Execute(w, data)
		if resErr != nil {
			utils.Log(utils.ERROR, "dashboard/loggedin/res", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "dashboard/loggedin/res", "Template rendered successfully")
		return
	}

	// not logged in user
	utils.Log(utils.INFO, "dashboard/checkSession", sessionErr.Error())

	tmpl, tmplErr := getDashboardTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "dashboard/notloggedin/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "dashboad/notloggedin/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := types.Page{
		Title:       "pengoe - Dashboard",
		Descrtipion: "Dashboard for pengoe",
		Session: types.Session{
			LoggedIn: false,
		},
	}

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		utils.Log(utils.ERROR, "dashboard/notloggedin/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "dashboard/notloggedin/res", "Template rendered successfully")
	return
}
