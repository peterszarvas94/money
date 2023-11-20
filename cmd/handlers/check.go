package handlers

import (
	"html"
	"html/template"
	"net/http"
	"net/mail"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/web/templates/components"
)

/*
CheckUserHandler checks if the username or email isaftaken.
Sends icons.
*/
func CheckUserHandler(w http.ResponseWriter, r *http.Request, pattern string) {
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

	// Parse the "correct" and "incorrect" templates
	correctTmpl, correctTmplErr := template.ParseFiles(components.Correct)
	if correctTmplErr != nil {
		logger.Log(logger.ERROR, "check/correct", correctTmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "check/correct", "Template parsed successfully")

	incorrectTmpl, incorrectTmplErr := template.ParseFiles(components.Incorrect)
	if incorrectTmplErr != nil {
		logger.Log(logger.ERROR, "check/incorrect", incorrectTmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "check/incorrect", "Template parsed successfully")

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
			resErr := correctTmpl.Execute(w, nil)
			if resErr != nil {
				logger.Log(logger.ERROR, "check/username/correctres", resErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			
			logger.Log(logger.INFO, "check/username/correctres", "Template rendered successfully")
			return
		}

		logger.Log(logger.WARNING, "check/username", "User exists")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := incorrectTmpl.Execute(w, nil)
		if resErr != nil {
			logger.Log(logger.ERROR, "check/username/incorrectres", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

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
			resErr := correctTmpl.Execute(w, nil)
			if resErr != nil {
				logger.Log(logger.ERROR, "check/email/correct", resErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			return
		}

		if emailParseErr != nil {
			logger.Log(logger.WARNING, "check/email", emailParseErr.Error())
		}

		if userErr == nil {
			logger.Log(logger.WARNING, "check/email", "User already exists")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := incorrectTmpl.Execute(w, nil)
		if resErr != nil {
			logger.Log(logger.ERROR, "check/email/incorrect", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		logger.Log(logger.INFO, "check/email/incorrect", "Template rendered successfully")
		return
	}
}
