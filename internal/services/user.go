package services

import (
	"database/sql"
	"errors"
	"net/http"
	"pengoe/internal/utils"
	"strconv"
	"time"
)

type UserService interface {
	New(user *utils.User) (*utils.User, error)
	Login(usernameOrEmail, password string) (*utils.User, error)
	GetById(id int) (*utils.User, error)
	GetByUsername(username string) (*utils.User, error)
	GetByEmail(email string) (*utils.User, error)
	CheckRefreshToken(r *http.Request) (*utils.User, error)
	CheckAccessToken(r *http.Request) (*utils.User, error)
}

type userService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return &userService{db: db}
}

/*
New is a function that adds a new user to the database.
*/
func (s *userService) New(user *utils.User) (*utils.User, error) {
	hashedPassword, hashErr := utils.HashPassword(user.Password)
	if hashErr != nil {
		return nil, hashErr
	}

	now := time.Now()

	mutation, mutationErr := s.db.Exec(
		`INSERT INTO user (
			username, email, firstname, lastname, password, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Username, user.Email, user.Fistname, user.Lastname, hashedPassword, now, now,
	)
	if mutationErr != nil {
		return nil, mutationErr
	}

	id, idErr := mutation.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}

	newUser := &utils.User{
		Id:        int(id),
		Username:  user.Username,
		Email:     user.Email,
		Fistname:  user.Fistname,
		Lastname:  user.Lastname,
		CreatedAt: now.String(),
		UpdatedAt: now.String(),
	}

	return newUser, nil
}

/*
Login is a function that gets a user from the database by username or email,
and checks if the passwords match.
If correct, it returns the user's id.
*/
func (s *userService) Login(usernameOrEmail, password string) (*utils.User, error) {
	query, queryErr := s.db.Query(
		`SELECT id, username, email, firstname, lastname, password as hash, created_at, updated_at
		FROM user WHERE username = ? OR email = ?`,
		usernameOrEmail, usernameOrEmail,
	)
	if queryErr != nil {
		return nil, queryErr
	}

	var id int
	var username string
	var email string
	var firstname string
	var lastname string
	var hash string
	var created_at string
	var updated_at string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &firstname, &lastname, &hash, &created_at, &updated_at)
			if scanErr != nil {
			return nil, scanErr
		}
	}

	if id == 0 {
		return nil, sql.ErrNoRows
	}

	match := utils.CheckPasswordHash(hash, password)
	if match != nil {
		return nil, sql.ErrNoRows
	}

	user := utils.User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return &user, nil
}

/*
GetByUsername is a function that gets a user from the database by username.
*/
func (s *userService) GetById(id int) (*utils.User, error) {
	var user utils.User

	query, queryErr := s.db.Query("SELECT * FROM user WHERE id = ?", id)
	if queryErr != nil {
		return &user, queryErr
	}

	var username string
	var email string
	var firstname string
	var lastname string
	var password string
	var created_at string
	var updated_at string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &firstname, &lastname, &password, &created_at, &updated_at)
		if scanErr != nil {
			return &user, scanErr
		}
	}


	if id == 0 {
		return &user, sql.ErrNoRows
	}

	user = utils.User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return &user, nil
}

/*
GetByUsername is a function that gets a user from the database by username.
*/
func (s *userService) GetByUsername(username string) (*utils.User, error) {
	var user utils.User

	query, queryErr := s.db.Query("SELECT * FROM user WHERE username = ?", username)
	if queryErr != nil {
		return &user, queryErr
	}

	var id int
	var email string
	var firstname string
	var lastname string
	var password string
	var created_at string
	var updated_at string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &firstname, &lastname, &password, &created_at, &updated_at)
		if scanErr != nil {
			return &user, scanErr
		}
	}

	if id == 0 {
		return &user, sql.ErrNoRows
	}

	user = utils.User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return &user, nil
}

/*
GetByEmail is a function that gets a user from the database by email.
*/
func (s *userService) GetByEmail(email string) (*utils.User, error) {
	var user utils.User

	query, queryErr := s.db.Query("SELECT * FROM user WHERE email = ?", email)
	if queryErr != nil {
		return &user, queryErr
	}

	var id int
	var username string
	var firstname string
	var lastname string
	var password string
	var created_at string
	var updated_at string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &firstname, &lastname, &password, &created_at, &updated_at)
		if scanErr != nil {
			return &user, scanErr
		}
	}

	if id == 0 {
		return &user, sql.ErrNoRows
	}

	user = utils.User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return &user, nil
}

/*
CheckRefreshToken is a function that checks if the refreshtoken in cookie is valid.
*/
func (s *userService) CheckRefreshToken(r *http.Request) (*utils.User, error) {
	refreshToken, cookieErr := r.Cookie("refresh")
	if cookieErr != nil {
		return nil, cookieErr
	}

	claims, jwtErr := utils.ValidateToken(refreshToken.Value)
	if jwtErr != nil {
		return nil, jwtErr
	}

	userIdStr, subErr := claims.GetSubject()
	if subErr != nil {
		return nil, subErr
	}

	userId, parseErr := strconv.Atoi(userIdStr)
	if parseErr != nil {
		return nil, parseErr
	}

	user, dbErr := s.GetById(userId)
	if dbErr != nil {
		return nil, dbErr
	}

	return user, nil
}

/*
CheckAccessToken is a function that checks if the accesstoken in auth header is valid.
*/
func (s *userService) CheckAccessToken(r *http.Request) (*utils.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("No authorization header")
	}

	authToken := authHeader[len("Bearer "):]

	claims, jwtErr := utils.ValidateToken(authToken)
	if jwtErr != nil {
		return nil, jwtErr
	}

	userIdStr, subErr := claims.GetSubject()
	if subErr != nil {
		return nil, subErr
	}

	userId, parseErr := strconv.Atoi(userIdStr)
	if parseErr != nil {
		return nil, parseErr
	}

	user, dbErr := s.GetById(userId)
	if dbErr != nil {
		return nil, dbErr
	}

	return user, nil
}
