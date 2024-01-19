package handlers

import (
	"database/sql"
	"errors"
	"fmt"
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
	token, tokenFound := r.Context().Value("token").(*t.Token)
	if !tokenFound {
		router.RedirectToSignin(w, r, p)
		return errors.New("Should use token middleware")
	}
	db, dbFound := r.Context().Value("db").(*sql.DB)
	if !dbFound {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}
	session, sessionFound := r.Context().Value("session").(*services.Session)
	if !sessionFound {
		router.InternalError(w, r, p)
		fmt.Println("Should use session middleware")
	}

	id, idFound := p["id"]
	if !idFound {
		router.NotFound(w, r, p)
		return errors.New("Path variable \"id\" not found")
	}

	// id to int
	accountId, errParse := strconv.Atoi(id)
	if errParse != nil {
		router.NotFound(w, r, p)
		return errParse
	}

	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)

	// get account
	account, accountErr := accountService.GetByID(accountId)
	if accountErr != nil {
		router.NotFound(w, r, p)
		return accountErr
	}

	// check if the user has access to the account
	accessErr := accessService.Check(session.UserId, accountId)
	if accessErr != nil {
		http.Redirect(w, r, "/dashboard", http.StatusUnauthorized)
		return accessErr
	}

	// get accounts
	accounts, accountsErr := accountService.GetByUserId(session.UserId)
	if accountsErr != nil {
		router.InternalError(w, r, p)
		return accountsErr
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := pages.AccountProps{
		Title:                fmt.Sprintf("pengoe - %s", account.Name),
		Description:          fmt.Sprintf("Account page for %s", account.Name),
		Accounts:             accounts,
		ShowNewAccountButton: true,
		Account:              account,
		Token:                token,
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
	token, tokenFound := r.Context().Value("token").(*t.Token)
	if !tokenFound {
		router.RedirectToSignin(w, r, p)
		return errors.New("Should use token middleware")
	}
	db, dbFound := r.Context().Value("db").(*sql.DB)
	if !dbFound {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}
	session, sessionFound := r.Context().Value("session").(*services.Session)
	if !sessionFound {
		router.InternalError(w, r, p)
		fmt.Println("Should use session middleware")
	}

	id, found := p["id"]
	if !found {
		router.NotFound(w, r, p)
		return errors.New("Path variable \"id\" not found")
	}

	// id to int
	accountId, errParse := strconv.Atoi(id)
	if errParse != nil {
		router.NotFound(w, r, p)
		return errParse
	}

	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)

	// check if the user has access to the account
	accessErr := accessService.Check(session.UserId, accountId)
	if accessErr != nil {
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
	body, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		router.InternalError(w, r, p)
		return bodyErr
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
	tokenFromReq, decodeErr := url.QueryUnescape(escapedTokenFromReq)
	if decodeErr != nil {
		router.InternalError(w, r, p)
		return decodeErr
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
		newCsrfToken, tokenErr := t.Manager.RenewToken(session.Id)
		if tokenErr != nil {
			router.InternalError(w, r, p)
			return tokenErr
		}

		// get account
		account, accountErr := accountService.GetByID(accountId)
		if accountErr != nil {
			router.InternalError(w, r, p)
			return accountErr
		}

		// get accounts
		accounts, accountsErr := accountService.GetByUserId(session.UserId)
		if accountsErr != nil {
			router.InternalError(w, r, p)
			return accountsErr
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
	accountErr := accountService.Delete(accountId)
	if accountErr != nil {
		router.NotFound(w, r, p)
		return accountErr
	}

	// redirect to dashboard
	w.Header().Set("HX-Redirect", "/dashboard")

	return nil

}
