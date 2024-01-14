package handlers

import (
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/serversession"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"time"
)

/*
SignoutHandler signs the user out by deleting the refresh token.
Access token is cleared by the client.
*/
func SignoutHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	secure := utils.Env.Environment == "production"
	var sameSite http.SameSite
	if utils.Env.Environment == "production" {
		sameSite = http.SameSiteLaxMode
	} else {
		sameSite = http.SameSiteDefaultMode
	}

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "signout/db", dbErr.Error())
		return dbErr
	}
	defer db.Close()

	// delete the session from the database
	sessionService := services.NewSessionService(db)

	// get old session form cookie
	session, sessionErr := sessionService.CheckCookie(r)
	if sessionErr != nil {
		logger.Log(logger.ERROR, "signout/session/cookie", sessionErr.Error())
		return sessionErr
	}

	// delete the session from the database
	deleteErr := sessionService.Delete(session.Id)
	if deleteErr != nil {
		logger.Log(logger.ERROR, "signout/session/db", deleteErr.Error())
		return deleteErr
	}

	// dete the server session
	serverSessionErr := serversession.Manager.Delete(session.Id)
	if serverSessionErr != nil {
		logger.Log(logger.ERROR, "signout/session/server", serverSessionErr.Error())
		return serverSessionErr
	}

	// delete the session cookie from the client
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour).UTC(),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})

	w.Header().Set("HX-Redirect", "/signin")

	logger.Log(
		logger.INFO,
		"signout/method",
		"User logged out and redirected to home page",
	)

	return nil
}
