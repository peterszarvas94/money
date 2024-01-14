package handlers

import (
	"fmt"
	"html"
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/router"
	"pengoe/internal/serversession"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"pengoe/web/templates/pages"
	"time"

	"github.com/a-h/templ"
)

/*
NewAccountPageHandler handles the GET request to /account/new
*/
func NewAccountPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
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
			"newaccount/session",
			fmt.Sprintf("Session found with ID %d", session.Id),
		)

		// check if there is a server session
		serverSession, serverSessionErr := serversession.Manager.Get(session.Id)
		if serverSessionErr != nil {
			logger.Log(
				logger.ERROR,
				"newaccount/session/server",
				serverSessionErr.Error(),
			)
			router.InternalError(w, r)
			return serverSessionErr
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// get accounts
		accounts, accountsErr := accountService.GetByUserId(session.UserId)
		if accountsErr != nil {
			logger.Log(
				logger.ERROR,
				"newaccount/session/accounts",
				accountsErr.Error(),
			)
			router.InternalError(w, r)
			return accountsErr
		}

		logger.Log(
			logger.INFO,
			"newaccount/session/accounts",
			fmt.Sprintf("Found %d accounts", len(accounts)),
		)

		// create account select items
		accountSelectItems := []utils.AccountSelectItem{}
		for _, account := range accounts {
			accountSelectItems = append(accountSelectItems, utils.AccountSelectItem{
				Id:   account.Id,
				Text: account.Name,
			})
		}

		data := pages.NewAccountProps{
			Title:                "pengoe - New Account",
			Description:          "New Account for pengoe",
			SelectedAccountId:    0,
			AccountSelectItems:   accountSelectItems,
			ShowNewAccountButton: false,
			CSRFToken:            serverSession.CSRFToken,
		}

		component := pages.NewAccount(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		logger.Log(
			logger.INFO,
			"newaccount/session/tmpl",
			"Template rendered successfully",
		)

		return nil
	}

	// not logged in user
	logger.Log(
		logger.INFO,
		"newaccount/nosession",
		fmt.Sprintf("No session found. %s", sessionErr.Error()),
	)

	logger.Log(
		logger.INFO,
		"newaccount/nosession/redirect",
		fmt.Sprintf("Redirecting to /signin?redirect=%%2Faccount%%2Fnew"),
	)

	http.Redirect(w, r, "/signin?redirect=%2Faccount%2Fnew", http.StatusSeeOther)

	return nil
}

/*
NewAccountHandler handles the POST request to /account
*/
func NewAccountHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	formErr := r.ParseForm()
	if formErr != nil {
		router.InternalError(w, r)
		return formErr
	}

	logger.Log(logger.INFO, "newaccount/post/form", "Form parsed successfully")

	name := html.EscapeString(r.FormValue("name"))
	description := html.EscapeString(r.FormValue("description"))
	currency := html.EscapeString(r.FormValue("currency"))
	csrfToken := html.EscapeString(r.FormValue("csrf"))

	if name == "" ||
		description == "" ||
		currency == "" ||
		csrfToken == "" {
		router.InternalError(w, r)
		logger.Log(
			logger.ERROR,
			"newaccount/post/form",
			"Some form values are empty",
		)
		return nil
	}

	logger.Log(
		logger.INFO,
		"newaccount/post/form",
		"Form values escaped successfully",
	)

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		router.InternalError(w, r)
		return dbErr
	}
	defer db.Close()

	logger.Log(logger.INFO, "newaccount/post/db", "Connected to db successfully")

	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)
	sessionService := services.NewSessionService(db)

	// check if the user is already logged in, redirect to dashboard
	session, sessionErr := sessionService.CheckCookie(r)
	if session != nil {
		// logged in user
		logger.Log(
			logger.INFO,
			"newaccount/session",
			fmt.Sprintf("Session found with ID %d", session.Id),
		)

		// check if there is a server session
		serverSession, serverSessionErr := serversession.Manager.Get(session.Id)
		if serverSessionErr != nil {
			logger.Log(
				logger.ERROR,
				"newaccount/session/server",
				serverSessionErr.Error(),
			)
			router.InternalError(w, r)
			return serverSessionErr
		}

		// check if the csrf token is valid
		if serverSession.CSRFToken.Token != csrfToken {
			logger.Log(
				logger.ERROR,
				"newaccount/session/csrf",
				"CSRF token is invalid",
			)
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		// check if the csrf token is expired
		if serverSession.CSRFToken.Valid.Before(time.Now().UTC()) {
			// csrf token is expired
			logger.Log(
				logger.ERROR,
				"newaccount/session/expired",
				"CSRF token is expired",
			)

			// renew csrf token
			newCsrfToken, tokenErr := serversession.Manager.RenewCSRFToken(session.Id)
			if tokenErr != nil {
				logger.Log(
					logger.ERROR,
					"newaccount/session/newcsrf",
					tokenErr.Error(),
				)
				router.InternalError(w, r)
				return tokenErr
			}

			logger.Log(
				logger.INFO,
				"newaccount/session/newcsrf",
				"CSRF token renewed successfully",
			)

			// get accounts
			accounts, accountsErr := accountService.GetByUserId(session.UserId)
			if accountsErr != nil {
				logger.Log(
					logger.ERROR,
					"newaccount/session/accounts",
					accountsErr.Error(),
				)
				router.InternalError(w, r)
				return accountsErr
			}

			logger.Log(
				logger.INFO,
				"newaccount/session/accounts",
				fmt.Sprintf("Found %d accounts", len(accounts)),
			)

			// create account select items
			accountSelectItems := []utils.AccountSelectItem{}
			for _, account := range accounts {
				accountSelectItems = append(accountSelectItems, utils.AccountSelectItem{
					Id:   account.Id,
					Text: account.Name,
				})
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			data := pages.NewAccountProps{
				Title:                "pengoe - New Account",
				Description:          "New Account for pengoe",
				SelectedAccountId:    0,
				AccountSelectItems:   accountSelectItems,
				ShowNewAccountButton: false,
				CSRFToken:            *newCsrfToken,
				AccountName:          name,
				AccountDescription:   description,
				AccountCurrency:      currency,
				Refetch:              true,
			}

			component := pages.NewAccount(data)
			handler := templ.Handler(component)
			handler.ServeHTTP(w, r)

			logger.Log(
				logger.INFO,
				"newaccount/session/tmpl",
				"Template rendered successfully",
			)

			return nil
		}

		// csrf token is not expired

		// create new account
		account := &utils.Account{
			Name:        name,
			Description: description,
			Currency:    currency,
		}

		newAccount, newAccountErr := accountService.New(account)
		if newAccountErr != nil {
			router.InternalError(w, r)
			return newAccountErr
		}

		logger.Log(logger.INFO, "newaccount/post/newAccount", "Account created successfully")

		// create new access
		access := &utils.Access{
			Role:      utils.Admin,
			UserId:    session.UserId,
			AccountId: newAccount.Id,
		}

		access, newAccessErr := accessService.New(access)
		if newAccessErr != nil {
			router.InternalError(w, r)
			return newAccessErr
		}

		logger.Log(logger.INFO, "newaccount/post/newAccess", "Access created successfully")

		w.Header().Set("HX-Redirect", fmt.Sprintf("/account/%d", newAccount.Id))
		return nil
	}

	// not logged in user
	if sessionErr != nil {
		router.InternalError(w, r)
		return sessionErr
	}

	return nil
}
