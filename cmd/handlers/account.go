package handlers

import (
	"fmt"
	"html"
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

type newAccountPage struct {
	Title                string
	Descrtipion          string
	Session              utils.Session
	SelectedAccountId    int
	AccountSelectItems   []utils.AccountSelectItem
	ShowNewAccountButton bool
}

/*
getNewAccountTmpl helper function to parse the newaccount template.
*/
func getNewAccountTmpl() (*template.Template, error) {
	tmpl, tmplErr := template.ParseFiles(
		layouts.Base,
		pages.NewAccount,
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
NewAccountPageHandler handles the GET request to /account/new
*/
func NewAccountPageHandler(w http.ResponseWriter, r *http.Request, pattern string) {

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

	// check if the user is logged in, protected route
	user, sessionErr := userService.CheckAccessToken(r)
	if user != nil {
		// logged in user
		logMsg := fmt.Sprintf("Logged in as %d", user.Id)
		logger.Log(logger.INFO, "newaccount/checkSession", logMsg)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := newAccountPage{
			Title:       "pengoe - New Account",
			Descrtipion: "New Account for pengoe",
			Session: utils.Session{
				LoggedIn: true,
				User:     *user,
			},
			SelectedAccountId:    0,
			ShowNewAccountButton: false,
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

		tmpl, tmplErr := getNewAccountTmpl()
		if tmplErr != nil {
			logger.Log(logger.ERROR, "newaccount/loggedin/tmpl", tmplErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		logger.Log(logger.INFO, "newaccount/loggedin/tmpl", "Template parsed successfully")

		resErr := tmpl.Execute(w, data)
		if resErr != nil {
			logger.Log(logger.ERROR, "newaccount/loggedin/res", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		logger.Log(logger.INFO, "newaccount/loggedin/res", "Template rendered successfully")
		return
	}

	// not logged in user
	logger.Log(logger.INFO, "newaccount/checkSession", sessionErr.Error())

	tmpl, tmplErr := getNewAccountTmpl()
	if tmplErr != nil {
		logger.Log(logger.ERROR, "newaccount/notloggedin/tmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	logger.Log(logger.INFO, "newaccount/notloggedin/tmpl", "Template parsed successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := utils.Page{
		Title:       "pengoe - New Account",
		Descrtipion: "New Account for pengoe",
		Session: utils.Session{
			LoggedIn: false,
		},
	}

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		logger.Log(logger.ERROR, "newaccount/notloggedin/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	logger.Log(logger.INFO, "newaccount/notloggedin/res", "Template rendered successfully")
	return
}

/*
NewAccountHandler handles the POST request to /account
*/
func NewAccountHandler(w http.ResponseWriter, r *http.Request, pattern string) {
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
			Role:       utils.Admin,
			UserId:     user.Id,
			AccountId:  newAccount.Id,
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
