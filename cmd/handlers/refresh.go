package handlers

import (
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/router"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"strconv"
	"time"
)

/*
RefreshTokenHandler gets a new access token and refresh token
if the refresh token is valid.
*/
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, rokenErr := utils.GetRefreshToken(r)
	if rokenErr != nil {
		router.InternalError(w, r)
		// ignore missing token error, it runs every page load
		return nil
	}

	logger.Log(logger.INFO, "refresh/token", "Refresh token received")

	// check if refresh token is valid
	claims, jwtErr := utils.ValidateToken(token)
	if jwtErr != nil {
		router.InternalError(w, r)
		return jwtErr
	}

	// get user id from claims
	userIdStr, subErr := claims.GetSubject()
	if subErr != nil {
		router.InternalError(w, r)
		return subErr
	}

	// convert user id to int
	userId, userIdErr := strconv.Atoi(userIdStr)
	if userIdErr != nil {
		router.InternalError(w, r)
		return userIdErr
	}

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		router.InternalError(w, r)
		return dbErr
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// check if user exists
	_, dbErr = userService.GetById(userId)
	if dbErr != nil {
		router.InternalError(w, r)
		return dbErr
	}

	// get new access token
	accessToken, accessTokenErr := utils.NewToken(userId, utils.AccessToken)
	if accessTokenErr != nil {
		router.InternalError(w, r)
		return accessTokenErr
	}

	// get new refresh token as well
	refreshToken, refreshTokenErr := utils.NewToken(userId, utils.RefreshToken)
	if refreshTokenErr != nil {
		router.InternalError(w, r)
		return refreshTokenErr
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

	logger.Log(logger.INFO, "refresh/res", "New tokens sent")

	return nil
}
