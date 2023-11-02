package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	_ "github.com/libsql/libsql-client-go/libsql"
)

/*
db is a function that returns a connection to the database.
It uses the DB_URL and DB_TOKEN environment variables to connect to the database.
*/
func db() (*sql.DB, error) {
	url, urlFound := os.LookupEnv("DB_URL")
	if !urlFound {
		return nil, errors.New("DB_URL not found")
	}

	token, token_found := os.LookupEnv("DB_TOKEN")
	if !token_found {
		return nil, errors.New("DB_TOKEN not found")
	}

	connectionStr := fmt.Sprintf("%s?authToken=%s", url, token)

	db, dbErr := sql.Open("libsql", connectionStr)
	if dbErr != nil {
		return nil, dbErr
	}

	return db, nil
}

func AddUser(newUsername string, newEmail string, newPassword string) (UserData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return UserData{}, dbErr
	}

	hashedPassword, hashErr := HashPassword(newPassword)
	if hashErr != nil {
		return UserData{}, hashErr
	}
	_, mutationErr := db.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		newUsername, newEmail, hashedPassword,
	)
	if mutationErr != nil {
		return UserData{}, mutationErr
	}

	query, queryErr := db.Query("SELECT * FROM users WHERE username = ?", newUsername)
	if queryErr != nil {
		return UserData{}, queryErr
	}

	var id int
	var username string
	var email string
	var password string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &password)
		if scanErr != nil {
			return UserData{}, scanErr
		}
	}

	if id == 0 {
		return UserData{}, errors.New("User not found")
	}

	user := UserData{
		Id:       id,
		Username: username,
		Email:    email,
	}

	return user, nil
}

/*
LoginUser is a function that returns a user from the database by username or email,
and checks if the passwords match.
*/
func LoginUser(usernameOrEmail string, password string) (UserData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return UserData{}, dbErr
	}

	query, queryErr := db.Query(
		"SELECT id, username, email, password as hash FROM users WHERE username = ? OR email = ?",
		usernameOrEmail, usernameOrEmail,
	)
	if queryErr != nil {
		return UserData{}, queryErr
	}

	var id int
	var username string
	var email string
	var hash string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &hash)
		if scanErr != nil {
			return UserData{}, scanErr
		}
	}

	if id == 0 {
		return UserData{}, errors.New("User not found")
	}

	matchErr := CheckPasswordHash(hash, password)
	if matchErr != nil {
		return UserData{}, matchErr
	}

	user := UserData{
		Id:       id,
		Username: username,
		Email:    email,
	}

	return user, nil

}

/*
GetUserById is a function that returns a user from the database by id.
*/
func GetUserById(id string) (UserData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return UserData{}, dbErr
	}

	query, queryErr := db.Query("SELECT * FROM users WHERE id = ?", id)
	if queryErr != nil {
		return UserData{}, queryErr
	}

	var username string
	var email string
	var password string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &password)
		if scanErr != nil {
			return UserData{}, scanErr
		}
	}

	idInt, idErr := strconv.Atoi(id)
	if idErr != nil {
		return UserData{}, idErr
	}

	user := UserData{
		Id:       idInt,
		Username: username,
		Email:    email,
	}

	return user, nil
}

/*
GetUserByUsername is a function that returns a user from the database by username.
*/
func GetUserByUsername(username string) (UserData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return UserData{}, dbErr
	}

	query, queryErr := db.Query("SELECT * FROM users WHERE username = ?", username)
	if queryErr != nil {
		return UserData{}, queryErr
	}

	var id int
	var email string
	var password string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &password)
		if scanErr != nil {
			return UserData{}, scanErr
		}
	}

	if id == 0 {
		return UserData{}, errors.New("User not found")
	}

	user := UserData{
		Id:       id,
		Username: username,
		Email:    email,
	}

	return user, nil
}

/*
GetUserByEmail is a function that returns a user from the database by email.
*/
func GetUserByEmail(email string) (UserData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return UserData{}, dbErr
	}

	query, queryErr := db.Query("SELECT * FROM users WHERE email = ?", email)
	if queryErr != nil {
		return UserData{}, queryErr
	}

	var id int
	var username string
	var password string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &password)
		if scanErr != nil {
			return UserData{}, scanErr
		}
	}

	if id == 0 {
		return UserData{}, errors.New("User not found")
	}

	user := UserData{
		Id:       id,
		Username: username,
		Email:    email,
	}

	return user, nil
}
