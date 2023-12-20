package handlers

import (
	"fmt"
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
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, rokenErr := utils.GetRefreshToken(r)
	if rokenErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println("refresh/token", rokenErr)
		return rokenErr
	}

	logger.Log(logger.INFO, "refresh/token", "Refresh token received")

	// check if refresh token is valid
	claims, jwtErr := utils.ValidateToken(token)
	if jwtErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return jwtErr
	}

	// get user id from claims
	userIdStr, subErr := claims.GetSubject()
	if subErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return subErr
	}

	// convert user id to int
	userId, userIdErr := strconv.Atoi(userIdStr)
	if userIdErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return userIdErr
	}

	// connect to db
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return dbErr
	}
	defer db.Close()
	userService := services.NewUserService(db)

	// check if user exists
	_, dbErr = userService.GetById(userId)
	if dbErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return dbErr
	}

	// get new access token
	accessToken, accessTokenErr := utils.NewToken(userId, utils.AccessToken)
	if accessTokenErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return accessTokenErr
	}

	// get new refresh token as well
	refreshToken, refreshTokenErr := utils.NewToken(userId, utils.RefreshToken)
	if refreshTokenErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
