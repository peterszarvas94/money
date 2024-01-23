package handlers

import (
	"database/sql"
	"errors"
	"html"
	"net/http"
	"net/mail"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/web/templates/components"

	"github.com/a-h/templ"
)

/*
CheckUserHandler checks if the username or email isaftaken.
Sends icons.
*/
func CheckUserHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	db, dbFound := r.Context().Value("db").(*sql.DB)
	if !dbFound {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}

	userService := services.NewUserService(db)

	// Parse the form
	parseErr := r.ParseForm()
	if parseErr != nil {
		router.InternalError(w, r, p)
		return parseErr
	}

	// Check if username is taken
	username := html.EscapeString(r.FormValue("username"))

	if username != "" {
		_, userErr := userService.GetByUsername(username)
		if userErr != nil {
			component := components.Correct()
			handler := templ.Handler(component)
			handler.ServeHTTP(w, r)
			return nil
		}
		
		component := components.Incorrect()
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)
		return nil
	}

	// Check if email is taken
	email := html.EscapeString(r.FormValue("email"))

	if email != "" {
		_, emailParseErr := mail.ParseAddress(email)
		_, userErr := userService.GetByEmail(email)
		if userErr != nil && emailParseErr == nil {
			component := components.Correct()
			handler := templ.Handler(component)
			handler.ServeHTTP(w, r)
			return nil
		}

		component := components.Incorrect()
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)
		return nil
	}

	return nil
}
