package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"pengoe/internal/router"
	"pengoe/internal/services"
	t "pengoe/internal/token"
	"pengoe/web/templates/pages"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

/*
AccountPageHandler handles the GET request to /account/:id
*/
func AccountPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, found := r.Context().Value("token").(*t.Token)
	if !found {
		router.RedirectToSignin(w, r, p)
		return errors.New("Should use token middleware")
	}
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}
	session, found := r.Context().Value("session").(*services.Session)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use session middleware")
	}

	id, found := p["id"]
	if !found {
		router.NotFound(w, r, p)
		return errors.New("Path variable \"id\" not found")
	}

	// id to int
	accountId, err := strconv.Atoi(id)
	if err != nil {
		router.NotFound(w, r, p)
		return err
	}

	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)
	eventService := services.NewEventService(db)

	// get account
	account, err := accountService.GetByID(accountId)
	if err != nil {
		router.NotFound(w, r, p)
		return err
	}

	// check if the user has access to the account
	err = accessService.Check(session.UserId, accountId)
	if err != nil {
		http.Redirect(w, r, "/dashboard", http.StatusUnauthorized)
		return err
	}

	// get accounts
	accounts, err := accountService.GetByUserId(session.UserId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	// get events
	events, err := eventService.GetByAccountId(accountId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	data := pages.AccountProps{
		Title:                fmt.Sprintf("pengoe - %s", account.Name),
		Description:          fmt.Sprintf("Account page for %s", account.Name),
		Accounts:             accounts,
		ShowNewAccountButton: true,
		Account:              account,
		Token:                token,
		Events:               events,
	}

	component := pages.Account(data)
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return nil
}

/*
DeleteAccountHandler handles the DELETE request to /account/:id
*/
func DeleteAccountHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, found := r.Context().Value("token").(*t.Token)
	if !found {
		router.RedirectToSignin(w, r, p)
		return errors.New("Should use token middleware")
	}
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}
	session, found := r.Context().Value("session").(*services.Session)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use session middleware")
	}

	id, found := p["id"]
	if !found {
		router.NotFound(w, r, p)
		return errors.New("Path variable \"id\" not found")
	}

	// id to int
	accountId, err := strconv.Atoi(id)
	if err != nil {
		router.NotFound(w, r, p)
		return err
	}

	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)

	// check if the user has access to the account
	err = accessService.Check(session.UserId, accountId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New(
			fmt.Sprintf(
				"User %d does not have access to account %d",
				session.UserId,
				accountId,
			),
		)
	}

	// manually parse body, (because DELETE request, go btw)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	// csrf=asd -> [csrf asd]
	splittedBody := strings.Split(string(body), "=")
	if len(splittedBody) != 2 || splittedBody[0] != "csrf" {
		router.Unauthorized(w, r, p)
		return errors.New("CSRF token is missing")
	}

	// [csrf asd] -> asd
	escapedTokenFromReq := splittedBody[1]

	// decode token
	tokenFromReq, err := url.QueryUnescape(escapedTokenFromReq)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	// check if the two tokens are the same
	if token.Value != tokenFromReq {
		router.Unauthorized(w, r, p)
		return errors.New("CSRF token is invalid")
	}

	// check if the csrf token is expired
	if token.Valid.Before(time.Now().UTC()) {
		// csrf token is expired

		// renew csrf token
		newCsrfToken, err := t.Manager.RenewToken(session.Id)
		if err != nil {
			router.InternalError(w, r, p)
			return err
		}

		// get account
		account, err := accountService.GetByID(accountId)
		if err != nil {
			router.InternalError(w, r, p)
			return err
		}

		// get accounts
		accounts, err := accountService.GetByUserId(session.UserId)
		if err != nil {
			router.InternalError(w, r, p)
			return err
		}

		data := pages.AccountProps{
			Title:                fmt.Sprintf("pengoe - %s", account.Name),
			Description:          fmt.Sprintf("Account page for %s", account.Name),
			Accounts:             accounts,
			ShowNewAccountButton: true,
			Account:              account,
			Token:                newCsrfToken,
		}

		component := pages.Account(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		return nil
	}

	// token is not expired, all good

	// delete account
	err = accountService.Delete(accountId)
	if err != nil {
		router.NotFound(w, r, p)
		return err
	}

	// redirect to dashboard
	w.Header().Set("HX-Redirect", "/dashboard")

	return nil
}

/*
NewAccountPageHandler handles the GET request to /account/new
*/
func NewAccountPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, found := r.Context().Value("token").(*t.Token)
	if !found {
		router.RedirectToSignin(w, r, p)
		return errors.New("Should use token middleware")
	}
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}
	session, found := r.Context().Value("session").(*services.Session)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use session middleware")
	}

	accountService := services.NewAccountService(db)

	// get accounts
	accounts, err := accountService.GetByUserId(session.UserId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
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
	token, found := r.Context().Value("token").(*t.Token)
	if !found {
		router.RedirectToSignin(w, r, p)
		return errors.New("Should use token middleware")
	}
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}
	session, found := r.Context().Value("session").(*services.Session)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use session middleware")
	}

	err := r.ParseForm()
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	form := r.Form

	formToken := html.EscapeString(form.Get("csrf"))
	if formToken == "" {
		router.Unauthorized(w, r, p)
		return errors.New("CSRF token is missing")
	}

	name := html.EscapeString(form.Get("name"))
	if name == "" {
		router.BadRequest(w, r, p)
		return errors.New("Name is required")
	}

	description := html.EscapeString(form.Get("description"))

	currency := html.EscapeString(form.Get("currency"))
	if currency == "" {
		router.BadRequest(w, r, p)
		return errors.New("Currency is required")
	}

	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)

	// check if the tokens match
	if token.Value != formToken {
		return router.Unauthorized(w, r, p)
	}

	// token is expired
	if token.Valid.Before(time.Now().UTC()) {
		newCsrfToken, err := t.Manager.RenewToken(session.Id)
		if err != nil {
			return router.InternalError(w, r, p)
		}

		accounts, err := accountService.GetByUserId(session.UserId)
		if err != nil {
			return router.InternalError(w, r, p)
		}

		data := pages.NewAccountProps{
			Title:                "pengoe - New Account",
			Description:          "New Account for pengoe",
			SelectedAccountId:    0,
			Accounts:             accounts,
			ShowNewAccountButton: false,
			Token:                newCsrfToken,
			Account: &services.Account{
				Name:        name,
				Description: description,
				Currency:    currency,
			},
			Refetch: true,
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

	newAccount, err := accountService.New(account)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	// create new access
	access := &services.Access{
		Role:      services.Admin,
		UserId:    session.UserId,
		AccountId: newAccount.Id,
	}

	_, err = accessService.New(access)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/account/%d", newAccount.Id))
	return nil
}
