package handlers

import (
	"fmt"
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"pengoe/web/templates/pages"

	"github.com/a-h/templ"
)

/*
DashboardPageHandler handles the GET request to /dashboard.
*/
func DashboardPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		router.InternalError(w, r)
		return dbErr
	}
	defer db.Close()

	userService := services.NewUserService(db)
	accountService := services.NewAccountService(db)

	// check if the user is logged in, protected route
	user, sessionErr := userService.CheckAccessToken(r)
	if user != nil {
		// logged in user
		logMsg := fmt.Sprintf("Logged in as %d", user.Id)
		logger.Log(logger.INFO, "dashboard/checkSession", logMsg)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// get accounts
		accounts, accountsErr := accountService.GetByUserId(user.Id)
		if accountsErr != nil {
			router.InternalError(w, r)
			return accountsErr
		}

		// create account select items
		accountSelectItems := []utils.AccountSelectItem{}
		for _, account := range accounts {
			accountSelectItems = append(accountSelectItems, utils.AccountSelectItem{
				Id:   account.Id,
				Text: account.Name,
			})
		}

		data := pages.DashboardProps{
			Title:       "pengoe - Dashboard",
			Description: "Dashboard for pengoe",
			Session: utils.Session{
				LoggedIn: true,
				User:     *user,
			},
			AccountSelectItems:   accountSelectItems,
			SelectedAccountId:    0,
			ShowNewAccountButton: true,
		}

		component := pages.Dashboard(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		logger.Log(logger.INFO, "dashboard/loggedin/tmpl", "Template parsed successfully")


		return nil
	}

	// not logged in user
	logger.Log(logger.WARNING, "dashboard/checkSession", sessionErr.Error())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := pages.DashboardProps{
		Title:       "pengoe - Dashboard",
		Description: "Dashboard for pengoe",
		Session: utils.Session{
			LoggedIn: false,
		},
	}

	component := pages.Dashboard(data)
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)


	return nil
}
