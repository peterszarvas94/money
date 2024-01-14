package handlers

import (
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/router"
	"pengoe/internal/serversession"
	"pengoe/internal/services"
	"pengoe/internal/utils"
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

	accountService := services.NewAccountService(db)
	accessService := services.NewAccessService(db)
	sessionService := services.NewSessionService(db)

	// get account early
	account, accountErr := accountService.GetByID(accountId)
	if accountErr != nil {
		router.NotFound(w, r)
		return accountErr
	}

	// check if the user is already logged in, redirect to dashboard
	session, sessionErr := sessionService.CheckCookie(r)
	if session != nil {
		// logged in user
		logger.Log(
			logger.INFO,
			"accountpage/session",
			fmt.Sprintf("Session found with ID %d", session.Id),
		)

    // check if the user has access to the account
		accessErr := accessService.Check(session.UserId, account.Id)
		if accessErr != nil {
			logger.Log(
				logger.WARNING,
				"accountpage/session/access",
				fmt.Sprintf(
					"User %d does not have access to account %d",
					session.UserId,
					account.Id,
				),
			)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return accessErr
		}

    // get csrf token
    serverSession, serverSessionError := serversession.Manager.Get(session.Id)
    if serverSessionError != nil {
      logger.Log(
        logger.ERROR,
        "accountpage/session/csrf",
        serverSessionError.Error(),
      )
      router.InternalError(w, r)
      return serverSessionError
    }

		// get accounts
		accounts, accountsErr := accountService.GetByUserId(session.UserId)
		if accountsErr != nil {
			logger.Log(
				logger.ERROR,
				"accountpage/session/accounts",
				accountsErr.Error(),
			)
			router.InternalError(w, r)
			return accountsErr
		}

		logger.Log(
			logger.INFO,
			"accountpage/session/accounts",
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

		data := pages.AccountProps{
			Title:                fmt.Sprintf("pengoe - %s", account.Name),
			Description:          fmt.Sprintf("Account page for %s", account.Name),
			AccountSelectItems:   accountSelectItems,
			ShowNewAccountButton: true,
			Account:              *account,
      CSRFToken:            serverSession.CSRFToken,
		}

		component := pages.Account(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		logger.Log(
			logger.INFO,
			"accountpage/loggedin/tmpl",
			"Template parsed successfully",
		)

		return nil
	}

	// not logged in user
	logger.Log(
		logger.INFO,
		"accountpage/nosession",
		fmt.Sprintf("No session found. %s", sessionErr.Error()),
	)

	logger.Log(
		logger.INFO,
		"accountpage/nosession/redirect",
		fmt.Sprintf("Redirecting to /signin?redirect=%%2Faccount%%2F%s", id),
	)

	http.Redirect(w, r, "/signin?redirect=%2Faccount%2F"+id, http.StatusSeeOther)

	return nil
}

/*
DeleteAccountHandler handles the DELETE request to /account/:id
*/
func DeleteAccountHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
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

	accountService := services.NewAccountService(db)
	sessionService := services.NewSessionService(db)
	accessService := services.NewAccessService(db)

	// check if the user is logged in
	session, sessionErr := sessionService.CheckCookie(r)
	if session != nil {
		// logged in user
		logger.Log(
			logger.INFO,
			"accountpage/checkSession",
			fmt.Sprintf("Logged in as %d", session.UserId),
		)

		// check if the user has access to the account
		accessErr := accessService.Check(session.UserId, accountId)
		if accessErr != nil {
			logger.Log(
				logger.WARNING,
				"accountpage/session/access",
				fmt.Sprintf(
					"User %d does not have access to account %d",
					session.UserId,
					accountId,
				),
			)
			w.WriteHeader(http.StatusUnauthorized)
			return accessErr
		}

    // manually parse body, to get csrf token (because DELETE request, thx go)
    body, bodyErr := io.ReadAll(r.Body)
    if bodyErr != nil {
      logger.Log(
        logger.ERROR,
        "accountpage/session/csrf",
        bodyErr.Error(),
      )
      router.InternalError(w, r)
      return bodyErr
    }

    // csrf=asd -> asd
    rawCsrfToken := html.EscapeString(strings.Split(string(body), "=")[1])

    // decode body as url values
    csrfToken, decodeErr := url.QueryUnescape(rawCsrfToken)
    if decodeErr != nil {
      logger.Log(
        logger.ERROR,
        "accountpage/session/csrf",
        decodeErr.Error(),
      )
      router.InternalError(w, r)
      return decodeErr
    }

    fmt.Println(csrfToken)

		// check if the csrf token is valid
		serverSession, serverSessionErr := serversession.Manager.Get(session.Id)
		if serverSessionErr != nil {
			logger.Log(
				logger.ERROR,
				"accountpage/session/csrf/server",
				serverSessionErr.Error(),
			)
			router.InternalError(w, r)
			return serverSessionErr
		}

    // check if the csrf token is the same
    if serverSession.CSRFToken.Token != csrfToken {
      logger.Log(
        logger.WARNING,
        "accountpage/session/csrf",
        "CSRF token is invalid",
      )
      w.WriteHeader(http.StatusUnauthorized)
      return errors.New("CSRF token is invalid")
    }

		// check if the csrf token is expired
		if serverSession.CSRFToken.Valid.Before(time.Now().UTC()) {
			// csrf token is expired
			logger.Log(
				logger.INFO,
				"accountpage/session/csrf/expired",
				"CSRF token is expired",
			)

      // renew csrf token
      newCsrfToken, tokenErr := serversession.Manager.RenewCSRFToken(session.Id)
      if tokenErr != nil {
        logger.Log(
          logger.ERROR,
          "accountpage/session/csrf/new",
          tokenErr.Error(),
        )
        router.InternalError(w, r)
        return tokenErr
      }

      logger.Log(
        logger.INFO,
        "accountpage/session/csrf/new",
        "CSRF token renewed successfully",
      )

      // get account
      account, accountErr := accountService.GetByID(accountId)
      if accountErr != nil {
        logger.Log(
          logger.ERROR,
          "accountpage/session/accounts",
          accountErr.Error(),
        )
        router.InternalError(w, r)
        return accountErr
      }

			// get accounts
			accounts, accountsErr := accountService.GetByUserId(session.UserId)
			if accountsErr != nil {
				logger.Log(
					logger.ERROR,
					"accountpage/session/accounts",
					accountsErr.Error(),
				)
				router.InternalError(w, r)
				return accountsErr
			}

			logger.Log(
				logger.INFO,
				"accountpage/session/accounts",
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

			data := pages.AccountProps{
				Title:                fmt.Sprintf("pengoe - %s", account.Name),
				Description:          fmt.Sprintf("Account page for %s", account.Name),
				AccountSelectItems:   accountSelectItems,
				ShowNewAccountButton: true,
				Account:              *account,
        CSRFToken:            *newCsrfToken,
			}

			component := pages.Account(data)
			handler := templ.Handler(component)
			handler.ServeHTTP(w, r)

			logger.Log(
				logger.INFO,
				"accountpage/loggedin/tmpl",
				"Template parsed successfully",
			)

			return nil
		}

    // token is not expired, all good
    logger.Log(
      logger.INFO,
      "accountpage/session/csrf",
      "CSRF token is valid",
    )

		// delete account
		accountErr := accountService.Delete(accountId)
		if accountErr != nil {
			router.NotFound(w, r)
			return accountErr
		}

		logger.Log(logger.INFO, "accountpage/delete", fmt.Sprintf("Account %d deleted", accountId))

		// redirect to dashboard
		w.Header().Set("HX-Redirect", "/dashboard")

		return nil
	}

	// not logged in user
	logger.Log(logger.WARNING, "accountpage/checkSession", sessionErr.Error())

	// redirect to login
	w.Header().Set("HX-Redirect", "/signin")

	return nil
}
