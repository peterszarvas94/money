package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/web/templates/pages"

	"github.com/a-h/templ"
)

/*
DashboardPage handles the GET request to /dashboard.
*/
func DashboardPage(w http.ResponseWriter, r *http.Request, p map[string]string) error {
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
	accounts, err := accountService.GetByUserId(session.UserId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
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
