package handlers

import (
	"net/http"
	"time"
)

func SignoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// clear jwt cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "jwt",
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-1 * time.Hour),
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
