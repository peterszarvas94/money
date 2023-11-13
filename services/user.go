package services

import (
	"database/sql"
	"errors"
	"net/http"
	"pengoe/types"
	"pengoe/utils"
	"strconv"
	"time"
)

type UserService interface {
	New(user *types.User) error
	Login(usernameOrEmail, password string) (int, error)
	GetById(id int) (*types.User, error)
	GetByUsername(username string) (*types.User, error)
	GetByEmail(email string) (*types.User, error)
	CheckRefreshToken(r *http.Request) (*types.User, error)
	CheckAccessToken(r *http.Request) (*types.User, error)
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
func (s *userService) New(user *types.User) error {
	hashedPassword, hashErr := utils.HashPassword(user.Password)
	if hashErr != nil {
		return hashErr
	}

	now := time.Now()

	_, mutationErr := s.db.Exec(
		`INSERT INTO user (
			username, email, firstname, lastname, password, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Username, user.Email, user.Fistname, user.Lastname, hashedPassword, now, now,
	)
	if mutationErr != nil {
		return mutationErr
	}

	return nil
}

/*
Login is a function that gets a user from the database by username or email,
and checks if the passwords match.
If correct, it returns the user's id.
*/
func (s *userService) Login(usernameOrEmail, password string) (int, error) {
	query, queryErr := s.db.Query(
		"SELECT id, password as hash FROM user WHERE username = ? OR email = ?",
		usernameOrEmail, usernameOrEmail,
	)
	if queryErr != nil {
		return 0, queryErr
	}

	var id int
	var hash string

	for query.Next() {
		scanErr := query.Scan(&id, &hash)
		if scanErr != nil {
			return 0, scanErr
		}
	}

	if id == 0 {
		return 0, sql.ErrNoRows
	}

	match := utils.CheckPasswordHash(hash, password)
	if match != nil {
		return 0, sql.ErrNoRows
	}

	return id, nil
}

/*
GetByUsername is a function that gets a user from the database by username.
*/
func (s *userService) GetById(id int) (*types.User, error) {
	var user types.User

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

	user = types.User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		Password:  password,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return &user, nil
}

/*
GetByUsername is a function that gets a user from the database by username.
*/
func (s *userService) GetByUsername(username string) (*types.User, error) {
	var user types.User

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

	user = types.User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		Password:  password,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return &user, nil
}

/*
GetByEmail is a function that gets a user from the database by email.
*/
func (s *userService) GetByEmail(email string) (*types.User, error) {
	var user types.User

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

	user = types.User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		Password:  password,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return &user, nil
}

/*
CheckRefreshToken is a function that checks if the refreshtoken in cookie is valid.
*/
func (s *userService) CheckRefreshToken(r *http.Request) (*types.User, error) {
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
func (s *userService) CheckAccessToken(r *http.Request) (*types.User, error) {
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
