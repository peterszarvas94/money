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
func SignoutHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/refresh",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	w.Header().Set("HX-Redirect", "/signin?redirect=%2Fdashboard")

	logger.Log(logger.INFO, "signout/method", "User logged out and redirected to home page")
	return nil
}
