package handlers

import (
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/internal/utils"
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
		logger.Log(logger.ERROR, "refresh/token", rokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, "refresh/token", "Refresh token received")

	// check if refresh token is valid
	claims, jwtErr := utils.ValidateToken(token)
	if jwtErr != nil {
		logger.Log(logger.ERROR, "refresh/token", jwtErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get user id from claims
	userIdStr, subErr := claims.GetSubject()
	if subErr != nil {
		logger.Log(logger.ERROR, "refresh/sub", subErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// convert user id to int
	userId, userIdErr := strconv.Atoi(userIdStr)
	if userIdErr != nil {
		logger.Log(logger.ERROR, "refresh/userid", userIdErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.ERROR, "signup/post/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// check if user exists
	_, dbErr = userService.GetById(userId)
	if dbErr != nil {
		logger.Log(logger.ERROR, "refresh/db", dbErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get new access token
	accessToken, accessTokenErr := utils.NewToken(userId, utils.AccessToken)
	if accessTokenErr != nil {
		logger.Log(logger.ERROR, "refresh/token", accessTokenErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// get new refresh token as well
	refreshToken, refreshTokenErr := utils.NewToken(userId, utils.RefreshToken)
	if refreshTokenErr != nil {
		logger.Log(logger.ERROR, "refresh/token", refreshTokenErr.Error())
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
