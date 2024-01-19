package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/web/templates/pages"

	"github.com/a-h/templ"
)

/*
DashboardPageHandler handles the GET request to /dashboard.
*/
func DashboardPageHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
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
	accounts, accountsErr := accountService.GetByUserId(session.UserId)
	if accountsErr != nil {
		router.InternalError(w, r, p)
		return accountsErr
	}

	data := pages.DashboardProps{
		Title:                "pengoe - Dashboard",
		Description:          "Dashboard for pengoe",
		Accounts:             accounts,
		SelectedAccountId:    0,
		ShowNewAccountButton: true,
	}

	component := pages.Dashboard(data)
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return nil
}
