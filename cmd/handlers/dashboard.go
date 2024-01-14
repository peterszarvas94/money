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

	accountService := services.NewAccountService(db)
	sessionService := services.NewSessionService(db)

	// check if the user is already logged in, redirect to dashboard
	session, sessionErr := sessionService.CheckCookie(r)
	if session != nil {
		// logged in user
		logger.Log(
			logger.INFO,
			"dashboard/session",
			fmt.Sprintf("Session found with ID %d", session.Id),
		)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// get accounts
		accounts, accountsErr := accountService.GetByUserId(session.UserId)
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
			Title:                "pengoe - Dashboard",
			Description:          "Dashboard for pengoe",
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
	logger.Log(
		logger.INFO,
		"dashboard/nosession",
		fmt.Sprintf("No session found. %s", sessionErr.Error()),
	)

	http.Redirect(w, r, "/signin?redirect=%2Fdashboard", http.StatusSeeOther)

	return nil
}
