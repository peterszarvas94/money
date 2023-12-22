package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"pengoe/web/templates/pages"
	"strconv"

	"github.com/a-h/templ"
)

/*
AccountPageHandler handles the GET request to /account/:id
*/
func AccountPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	id, found := p["id"]
	if !found {
		router.NotFound(w, r)
		return errors.New("Path variable \"id\" not found")
	}

	// id to int
	accountId, errParse := strconv.Atoi(id)
	if errParse != nil {
		router.NotFound(w, r)
		return errParse
	}

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

	// get account early
	account, accountErr := accountService.GetById(accountId)
	if accountErr != nil {
		router.NotFound(w, r)
		return accountErr
	}

	// check if the user is logged in, protected route
	user, sessionErr := userService.CheckAccessToken(r)
	if user != nil {
		// logged in user
		logMsg := fmt.Sprintf("Logged in as %d", user.Id)
		logger.Log(logger.INFO, "dashboard/checkSession", logMsg)

		// TODO: check if the user has access to the account

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

		data := pages.AccountProps{
			Title:       "pengoe - Dashboard",
			Description: "Dashboard for pengoe",
			Session: utils.Session{
				LoggedIn: true,
				User:     *user,
			},
			AccountSelectItems:   accountSelectItems,
			ShowNewAccountButton: true,
			Account:              *account,
		}

		component := pages.Account(data)
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
