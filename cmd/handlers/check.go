package handlers

import (
	"html"
	"net/http"
	"net/mail"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/web/templates/components"

	"github.com/a-h/templ"
)

/*
CheckUserHandler checks if the username or email isaftaken.
Sends icons.
*/
func CheckUserHandler(w http.ResponseWriter, r *http.Request) {
	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "check/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	userService := services.NewUserService(db)

	// Parse the form
	parseErr := r.ParseForm()
	if parseErr != nil {
		logger.Log(logger.ERROR, "check/parse", parseErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "check/parse", "Form parsed successfully")

	// Check if username is taken
	username := html.EscapeString(r.FormValue("username"))

	if username != "" {
		_, userErr := userService.GetByUsername(username)
		if userErr != nil {
			logger.Log(logger.INFO, "check/usename", userErr.Error())

			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			component := components.Correct()
			handler := templ.Handler(component)
			handler.ServeHTTP(w, r)

			logger.Log(logger.INFO, "check/username/correctres", "Template rendered successfully")
			return
		}

		logger.Log(logger.WARNING, "check/username", "User exists")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		component := components.Incorrect()
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		logger.Log(logger.INFO, "check/username/incorrectres", "Template rendered successfully")
		return
	}

	// Check if email is taken
	email := html.EscapeString(r.FormValue("email"))

	if email != "" {
		_, emailParseErr := mail.ParseAddress(email)
		_, userErr := userService.GetByEmail(email)
		if userErr != nil && emailParseErr == nil {
			logger.Log(logger.INFO, "check/email", userErr.Error())

			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			component := components.Correct()
			handler := templ.Handler(component)
			handler.ServeHTTP(w, r)

			return
		}

		if emailParseErr != nil {
			logger.Log(logger.WARNING, "check/email", emailParseErr.Error())
		}

		if userErr == nil {
			logger.Log(logger.WARNING, "check/email", "User already exists")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		component := components.Incorrect()
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		logger.Log(logger.INFO, "check/email/incorrect", "Template rendered successfully")
		return
	}
}
