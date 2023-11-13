package handlers

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"pengoe/db"
	"pengoe/services"
	"pengoe/types"
	"pengoe/utils"
	"time"
)

/*
getSigninTmpl helper function to parse the signin template.
*/
func getSigninTmpl() (*template.Template, error) {
	baseHtml := "templates/layouts/base.html"
	signinHtml := "templates/pages/signin.html"
	iconHtml := "templates/components/icon.html"
	errorHtml := "templates/components/error.html"

	tmpl, tmplErr := template.ParseFiles(baseHtml, signinHtml, iconHtml, errorHtml)
	if tmplErr != nil {
		return nil, tmplErr
	}

	return tmpl, nil
}

/*
SigninPageHandler handles the GET request to /signin.
*/
func SigninPageHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	params := utils.GetQueryParams(r)

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		utils.Log(utils.ERROR, "signin/get/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// check if the user is already logged in, redirect to dashboard
	user, accessTokenErr := userService.CheckAccessToken(r)
	if user != nil {
		logMsg := fmt.Sprintf("Logged in as %d, redirecting to dashboard", user.Id)
		utils.Log(utils.INFO, "signin/checkSession", logMsg)

		redirect := params["redirect"]
		if redirect == "" {
			redirect = "dashboard"
		}

		w.Header().Set("HX-Redirect", "/"+redirect)
		return
	}

	utils.Log(utils.INFO, "signin/checkSession", accessTokenErr.Error())

	tmpl, tmplErr := getSigninTmpl()
	if tmplErr != nil {
		utils.Log(utils.ERROR, "signin/signinTmpl", tmplErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "Template parsed successfully", "signin/signinTmpl")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := types.Page{
		Title:       "pengoe - Sign in",
		Descrtipion: "Sign in to pengoe",
		Data: map[string]string{
			"redirect": params["redirect"],
		},
	}

	resErr := tmpl.Execute(w, data)
	if resErr != nil {
		utils.Log(utils.ERROR, "signin/get/res", resErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	utils.Log(utils.INFO, "signin/get/res", "Template rendered successfully")
	return
}

/*
SigninHandler handles the POST request to /signin.
*/
func SigninHandler(w http.ResponseWriter, r *http.Request, pattern string) {

	fmt.Println(r.URL.Path)

	formErr := r.ParseForm()
	if formErr != nil {
		utils.Log(utils.ERROR, "signin/post/parse", formErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "signin/post/parse", "Form parsed successfully")

	usernameOrEmail := html.EscapeString(r.FormValue("user"))
	password := html.EscapeString(r.FormValue("password"))

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		utils.Log(utils.ERROR, "signin/post/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// login the user
	userId, loginErr := userService.Login(usernameOrEmail, password)

	// if the login was unsuccessful
	if loginErr != nil {
		utils.Log(utils.ERROR, "signin/post/login", loginErr.Error())

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)

		signinData := types.Page{
			Title:       "pengoe - Sign in",
			Descrtipion: "Sign in to pengoe",
			Data: map[string]string{
				"usernameOrEmail": usernameOrEmail,
				"loginError": "Invalid username or password",
			},
		}

		tmpl, tmplErr := getSigninTmpl()
		if tmplErr != nil {
			utils.Log(utils.ERROR, "signin/post/errorTmpl", tmplErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "signin/post/errorTmpl", "Template parsed successfully")

		resErr := tmpl.Execute(w, signinData)
		if resErr != nil {
			utils.Log(utils.ERROR, "signin/post/errorRes", resErr.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		utils.Log(utils.INFO, "signin/post/errorRes", "Template rendered successfully")
		return
	}

	// if the login was successful
	accessToken, accessTokenErr := utils.NewToken(userId, types.Access)
	if accessTokenErr != nil {
		utils.Log(utils.ERROR, "signin/post/access", accessTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "signin/post/access", "Access token created successfully")

	refreshToken, refreshTokenErr := utils.NewToken(userId, types.Refresh)
	if refreshTokenErr != nil {
		utils.Log(utils.ERROR, "signin/post/refresh", refreshTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "signin/post/refresh", "Refresh token created successfully")

	expires := time.Unix(refreshToken.Expires, 0)

	// set the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    refreshToken.Token,
		Path:     "/refresh",
		Expires:  expires,
		HttpOnly: true,
		// Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// set the access token header
	w.Header().Set("Authorization", "Bearer "+accessToken.Token)

	params := utils.GetQueryParams(r)

	redirect := params["redirect"]
	url := "/"+redirect

	http.Redirect(w, r, url, http.StatusSeeOther)

	utils.Log(utils.INFO, "signin/post/res", "Redirected to "+url)
	return
}
