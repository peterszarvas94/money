package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
NewToken is a function that returns a new JWT.
It takes an id and a tokenvariant as arguments.
Example:
NewToken(1, ACCESS)
*/
func NewToken(id int, variant TokenVariant) (JWT, error) {
	currentTime := time.Now().Unix()
	expirationTime := currentTime + 3600

	secret := Env.JWTSecret

	idStr := strconv.Itoa(id)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": currentTime,
		"exp": expirationTime,
		"sub": idStr,
		"typ": variant,
	})

	signedToken, signErr := token.SignedString([]byte(secret))
	if signErr != nil {
		return JWT{}, signErr
	}

	return JWT{
		Token:   signedToken,
		Expires: expirationTime,
	}, nil
}

func ValidateToken(token string) (jwt.MapClaims, error) {
	secret := Env.JWTSecret

	parsedToken, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Error: Unexpected signing method")
		}

		return []byte(secret), nil
	})

	if parseErr != nil {
		return nil, parseErr
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("Error: Invalid signature")
	}

	return claims, nil
}

/*
GetAccessToken is a function that returns the access token from the authorization header.
*/
func GetAccessToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid authorization header")
	}

	return parts[1], nil
}

/*
GetRefreshToken is a function that returns the refresh token from the refresh cookie.
*/
func GetRefreshToken(r *http.Request) (string, error) {
	cookie, cookieErr := r.Cookie("refresh")
	if cookieErr != nil {
		return "", cookieErr
	}

	return cookie.Value, nil
}
