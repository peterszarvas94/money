package handlers

import (
	"database/sql"
	"fmt"
	"html"
	"net/http"
	t "pengoe/internal/token"
	"pengoe/internal/logger"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/web/templates/pages"
	"time"

	"github.com/a-h/templ"
)

/*
NewAccountPageHandler handles the GET request to /account/new
*/
func NewAccountPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, tokenFound := r.Context().Value("token").(*t.Token)
	if !tokenFound {
		router.RedirectToSignin(w, r, p)
    fmt.Println("Should use token middleware")
	}
	db, dbFound := r.Context().Value("db").(*sql.DB)
	if !dbFound {
		router.InternalError(w, r, p)
    fmt.Println("Should use db middleware")
	}
	session, sessionFound := r.Context().Value("session").(*services.Session)
	if !sessionFound {
		router.InternalError(w, r, p)
    fmt.Println("Should use session middleware")
	}

	accountService := services.NewAccountService(db)

	// get accounts
	accounts, accountsErr := accountService.GetByUserId(session.UserId)
	if accountsErr != nil {
		logger.Log(
			logger.ERROR,
			"newaccount/session/accounts",
			accountsErr.Error(),
		)
		router.InternalError(w, r, p)
		return accountsErr
	}

	data := pages.NewAccountProps{
		Title:                "pengoe - New Account",
		Description:          "New Account for pengoe",
		SelectedAccountId:    0,
		Accounts:             accounts,
		ShowNewAccountButton: false,
		Token:                token,
	}

	component := pages.NewAccount(data)
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return nil
}

/*
NewAccountHandler handles the POST request to /account
*/
func NewAccountHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, tokenFound := r.Context().Value("token").(*t.Token)
	if !tokenFound {
		router.RedirectToSignin(w, r, p)
    fmt.Println("Should use token middleware")
	}
	db, dbFound := r.Context().Value("db").(*sql.DB)
	if !dbFound {
		router.InternalError(w, r, p)
    fmt.Println("Should use db middleware")
	}
	session, sessionFound := r.Context().Value("session").(*services.Session)
	if !sessionFound {
		router.InternalError(w, r, p)
    fmt.Println("Should use session middleware")
	}

	formErr := r.ParseForm()
	if formErr != nil {
		router.InternalError(w, r, p)
		return formErr
	}

	name := html.EscapeString(r.FormValue("name"))
	description := html.EscapeString(r.FormValue("description"))
	currency := html.EscapeString(r.FormValue("currency"))
	csrfTokenFromForm := html.EscapeString(r.FormValue("csrf"))

	if name == "" ||
		description == "" ||
		currency == "" ||
		csrfTokenFromForm == "" {
		router.InternalError(w, r, p)
		logger.Log(
			logger.ERROR,
			"newaccount/post/form",
			"Some form values are empty",
		)
		return nil
	}

	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)

	// check if the tokens match
	if token.Value != csrfTokenFromForm {
		return router.Unauthorized(w, r, p)
	}

  // token is expired
	if token.Valid.Before(time.Now().UTC()) {
		newCsrfToken, tokenErr := t.Manager.RenewToken(session.Id)
		if tokenErr != nil {
			return router.InternalError(w, r, p)
		}

		accounts, accountsErr := accountService.GetByUserId(session.UserId)
		if accountsErr != nil {
			return router.InternalError(w, r, p)
		}

		data := pages.NewAccountProps{
			Title:                "pengoe - New Account",
			Description:          "New Account for pengoe",
			SelectedAccountId:    0,
			Accounts:             accounts,
			ShowNewAccountButton: false,
			Token:                newCsrfToken,
			AccountName:          name,
			AccountDescription:   description,
			AccountCurrency:      currency,
			Refetch:              true,
		}

		component := pages.NewAccount(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		return nil
	}

	// csrf token is not expired

	// create new account
	account := &services.Account{
		Name:        name,
		Description: description,
		Currency:    currency,
	}

	newAccount, newAccountErr := accountService.New(account)
	if newAccountErr != nil {
		router.InternalError(w, r, p)
		return newAccountErr
	}

	// create new access
	access := &services.Access{
		Role:      services.Admin,
		UserId:    session.UserId,
		AccountId: newAccount.Id,
	}

	access, newAccessErr := accessService.New(access)
	if newAccessErr != nil {
		router.InternalError(w, r, p)
		return newAccessErr
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/account/%d", newAccount.Id))
	return nil
}
