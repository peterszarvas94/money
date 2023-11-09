package handlers

import (
	"pengoe/db"
	"pengoe/services"
	"pengoe/types"
	"pengoe/utils"
	"net/http"
	"strconv"
	"time"
)

/*
RefreshTokenHandler gets a new access token and refresh token
if the refresh token is valid.
*/
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, pattern string) {
	token, rokenErr := utils.GetRefreshToken(r)
	if rokenErr != nil {
		utils.Log(utils.ERROR, "refresh/token", rokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.Log(utils.INFO, "refresh/token", "Refresh token received")

	// check if refresh token is valid
	claims, jwtErr := utils.ValidateToken(token)
	if jwtErr != nil {
		utils.Log(utils.ERROR, "refresh/token", jwtErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get user id from claims
	userIdStr, subErr := claims.GetSubject()
	if subErr != nil {
		utils.Log(utils.ERROR, "refresh/sub", subErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// convert user id to int
	userId, userIdErr := strconv.Atoi(userIdStr)
	if userIdErr != nil {
		utils.Log(utils.ERROR, "refresh/userid", userIdErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		utils.Log(utils.ERROR, "signup/post/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// check if user exists
	_, dbErr = userService.GetById(userId)
	if dbErr != nil {
		utils.Log(utils.ERROR, "refresh/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get new access token
	accessToken, accessTokenErr := utils.NewToken(userId, types.Access)
	if accessTokenErr != nil {
		utils.Log(utils.ERROR, "refresh/token", accessTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get new refresh token as well
	refreshToken, refreshTokenErr := utils.NewToken(userId, types.Refresh)
	if refreshTokenErr != nil {
		utils.Log(utils.ERROR, "refresh/token", refreshTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	expires := time.Unix(refreshToken.Expires, 0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// set the new refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    refreshToken.Token,
		Path:     "/refresh",
		Expires:  expires,
		HttpOnly: true,
		// Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	})

	// send access token in header
	w.Header().Set("Authorization", "Bearer "+accessToken.Token)
	w.WriteHeader(http.StatusOK)
}
