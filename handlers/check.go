package handlers

import (
	"pengoe/db"
	"pengoe/services"
	"pengoe/utils"
	"html"
	"html/template"
	"net/http"
	"net/mail"
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
		utils.Log(utils.ERROR, "check/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	userService := services.NewUserService(db)

	// Parse the "correct" and "incorrect" templates
	correct := "templates/components/correct.html"
	correctTmpl, correctTmplErr := template.ParseFiles(correct)
	if correctTmplErr != nil {
		utils.Log(utils.ERROR, "check/correct", correctTmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "check/correct", "Template parsed successfully")

	incorrect := "templates/components/incorrect.html"
	incorrectTmpl, incorrectTmplErr := template.ParseFiles(incorrect)
	if incorrectTmplErr != nil {
		utils.Log(utils.ERROR, "check/incorrect", incorrectTmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "check/incorrect", "Template parsed successfully")

	// Parse the form
	parseErr := r.ParseForm()
	if parseErr != nil {
		utils.Log(utils.ERROR, "check/parse", parseErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "check/parse", "Form parsed successfully")

	// Check if username is taken
	username := html.EscapeString(r.FormValue("username"))

	if username != "" {
		_, userErr := userService.GetByUsername(username)
		if userErr != nil {
			utils.Log(utils.INFO, "check/usename", userErr.Error())

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			resErr := correctTmpl.Execute(w, nil)
			if resErr != nil {
				utils.Log(utils.ERROR, "check/username/correctres", resErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			
			utils.Log(utils.INFO, "check/username/correctres", "Template rendered successfully")
			return
		}

		utils.Log(utils.WARNING, "check/username", "User exists")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := incorrectTmpl.Execute(w, nil)
		if resErr != nil {
			utils.Log(utils.ERROR, "check/username/incorrectres", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "check/username/incorrectres", "Template rendered successfully")
		return
	}

	// Check if email is taken
	email := html.EscapeString(r.FormValue("email"))

	if email != "" {
		_, emailParseErr := mail.ParseAddress(email)
		_, userErr := userService.GetByEmail(email)
		if userErr != nil && emailParseErr == nil {
			utils.Log(utils.INFO, "check/email", userErr.Error())

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			resErr := correctTmpl.Execute(w, nil)
			if resErr != nil {
				utils.Log(utils.ERROR, "check/email/correct", resErr.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			return
		}

		if emailParseErr != nil {
			utils.Log(utils.WARNING, "check/email", emailParseErr.Error())
		}

		if userErr == nil {
			utils.Log(utils.WARNING, "check/email", "User already exists")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resErr := incorrectTmpl.Execute(w, nil)
		if resErr != nil {
			utils.Log(utils.ERROR, "check/email/incorrect", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "check/email/incorrect", "Template rendered successfully")
		return
	}
}
