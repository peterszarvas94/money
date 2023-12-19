package handlers

import (
	"fmt"
	"html"
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"pengoe/web/templates/pages"

	"github.com/a-h/templ"
)

type newAccountPage struct {
	Title                string
	Descrtipion          string
	Session              utils.Session
	SelectedAccountId    int
	AccountSelectItems   []utils.AccountSelectItem
	ShowNewAccountButton bool
}

/*
NewAccountPageHandler handles the GET request to /account/new
*/
func NewAccountPageHandler(w http.ResponseWriter, r *http.Request) {

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "newaccount/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	userService := services.NewUserService(db)
	accountService := services.NewAccountService(db)

	// check if the user is logged in, protected route
	user, sessionErr := userService.CheckAccessToken(r)
	if user != nil {
		// logged in user
		logMsg := fmt.Sprintf("Logged in as %d", user.Id)
		logger.Log(logger.INFO, "newaccount/checkSession", logMsg)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// get accounts
		accounts, accountsErr := accountService.GetByUserId(user.Id)
		if accountsErr != nil {
			logger.Log(logger.ERROR, "dashboard/accounts", accountsErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// create account select items
		accountSelectItems := []utils.AccountSelectItem{}
		for _, account := range accounts {
			accountSelectItems = append(accountSelectItems, utils.AccountSelectItem{
				Id:   account.Id,
				Text: account.Name,
			})
		}

		data := pages.NewAccountProps{
			Title:       "pengoe - New Account",
			Description: "New Account for pengoe",
			Session: &utils.Session{
				LoggedIn: true,
			},
			SelectedAccountId:    0,
			AccountSelectItems:   accountSelectItems,
			ShowNewAccountButton: false,
		}

		component := pages.NewAccount(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		logger.Log(logger.INFO, "newaccount/loggedin/tmpl", "Template parsed successfully")

		logger.Log(logger.INFO, "newaccount/loggedin/res", "Template rendered successfully")
		return
	}

	// not logged in user
	logger.Log(logger.INFO, "newaccount/checkSession", sessionErr.Error())

	data := pages.NewAccountProps{
		Title:       "pengoe - New Account",
		Description: "New Account for pengoe",
		Session: &utils.Session{
			LoggedIn: true,
		},
	}

	component := pages.NewAccount(data)
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	logger.Log(logger.INFO, "newaccount/notloggedin/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	logger.Log(logger.INFO, "newaccount/notloggedin/res", "Template rendered successfully")
	return
}

/*
NewAccountHandler handles the POST request to /account
*/
func NewAccountHandler(w http.ResponseWriter, r *http.Request) {
	formErr := r.ParseForm()
	if formErr != nil {
		logger.Log(logger.ERROR, "newaccount/post/form", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "newaccount/post/form", "Form parsed successfully")

	name := html.EscapeString(r.FormValue("name"))
	description := html.EscapeString(r.FormValue("description"))
	currency := html.EscapeString(r.FormValue("currency"))

	// TODO: handle empty values

	logger.Log(logger.INFO, "newaccount/post/form", "Form values escaped successfully")

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "newaccount/post/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "newaccount/post/db", "Connected to db successfully")

	defer db.Close()

	userService := services.NewUserService(db)
	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)

	// check if the user is logged in, protected route
	user, sessionErr := userService.CheckAccessToken(r)
	if user != nil {
		// logged in user
		logMsg := fmt.Sprintf("Logged in as %d", user.Id)
		logger.Log(logger.INFO, "newaccount/post/checkSession", logMsg)

		// create new account
		account := &utils.Account{
			Name:        name,
			Description: description,
			Currency:    currency,
		}

		newAccount, newAccountErr := accountService.New(account)
		if newAccountErr != nil {
			logger.Log(logger.ERROR, "newaccount/post/newAccount", newAccountErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.Log(logger.INFO, "newaccount/post/newAccount", "Account created successfully")

		// create new access
		access := &utils.Access{
			Role:      utils.Admin,
			UserId:    user.Id,
			AccountId: newAccount.Id,
		}

		access, newAccessErr := accessService.New(access)
		if newAccessErr != nil {
			logger.Log(logger.ERROR, "newaccount/post/newAccess", newAccessErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.Log(logger.INFO, "newaccount/post/newAccess", "Access created successfully")

		w.Header().Set("HX-Redirect", "/dashboard")
	}

	// not logged in user
	if sessionErr != nil {
		logger.Log(logger.INFO, "newaccount/post/checkSession", sessionErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
