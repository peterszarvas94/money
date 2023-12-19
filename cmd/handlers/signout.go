package handlers

import (
	"net/http"
	"pengoe/internal/logger"
	"time"
)

/*
SignoutHandler signs the user out by deleting the refresh token.
Access token is cleared by the client.
	"pengoe/utils"
*/
func SignoutHandler(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/refresh",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/signin?redirect=%2Fdashboard", http.StatusSeeOther)

	logger.Log(logger.INFO, "signout/method", "User logged out and redirected to home page")
	return nil
}
