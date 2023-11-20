package handlers

import (
	"fmt"
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

type dashboardPage struct {
	Title                string
	Descrtipion          string
	Session              utils.Session
	SelectedAccountId    int
	AccountSelectItems   []utils.AccountSelectItem
	ShowNewAccountButton bool
}

/*
getDashboardTmpl helper function to parse the dashboard template.
*/
func getDashboardTmpl() (*template.Template, error) {
	tmpl, tmplErr := template.ParseFiles(
		layouts.Base,
		pages.Dashboard,
		components.LeftPanel,
		components.TopBar,
		components.Icon,
		components.AccountSelectItem,
		components.Spinner,
	)
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
		logger.Log(logger.ERROR, "dashboard/db", dbErr.Error())
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
		logger.Log(logger.INFO, "dashboard/checkSession", logMsg)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := dashboardPage{
			Title:       "pengoe - Dashboard",
			Descrtipion: "Dashboard for pengoe",
			Session: utils.Session{
				LoggedIn: true,
				User:     *user,
			},
			SelectedAccountId:    1,
			ShowNewAccountButton: true,
			AccountSelectItems: []utils.AccountSelectItem{
				{
					Id:   1,
					Text: "Account 1",
				},
				{
					Id:   2,
					Text: "Account 2",
				},
				{
					Id:   3,
					Text: "Account 3",
				},
			},
		}

		tmpl, tmplErr := getDashboardTmpl()
		if tmplErr != nil {
			logger.Log(logger.ERROR, "dashboard/loggedin/tmpl", tmplErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		logger.Log(logger.INFO, "dashboard/loggedin/tmpl", "Template parsed successfully")

		resErr := tmpl.Execute(w, data)
		if resErr != nil {
			logger.Log(logger.ERROR, "dashboard/loggedin/res", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		logger.Log(logger.INFO, "dashboard/loggedin/res", "Template rendered successfully")
		return
	}

	// not logged in user
	logger.Log(logger.INFO, "dashboard/checkSession", sessionErr.Error())

	tmpl, tmplErr := getDashboardTmpl()
	if tmplErr != nil {
		logger.Log(logger.ERROR, "dashboard/notloggedin/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	logger.Log(logger.INFO, "dashboad/notloggedin/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := utils.Page{
		Title:       "pengoe - Dashboard",
		Descrtipion: "Dashboard for pengoe",
		Session: utils.Session{
			LoggedIn: false,
		},
	}

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		logger.Log(logger.ERROR, "dashboard/notloggedin/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	logger.Log(logger.INFO, "dashboard/notloggedin/res", "Template rendered successfully")
	return
}
