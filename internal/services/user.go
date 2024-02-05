package services

import (
	"database/sql"
	"errors"
	"pengoe/internal/utils"
	"time"
)

type User struct {
	Id        string
	Username  string
	Email     string
	Fistname  string
	Lastname  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserServiceInterface interface {
	Signup(id, username, email, firstname, lastname, password string) error
	Signin(usernameOrEmail, password string) (string, error)
	GetById(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
}

type userService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserServiceInterface {
	return &userService{db: db}
}

/*
Signup is a function that adds a new user to the database.
*/
func (s *userService) Signup(id, username, email, firstname, lastname, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	_, err = s.db.Exec(
		`INSERT INTO user (
			id,
			username,
			email,
			firstname,
			lastname,
			password,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		username,
		email,
		firstname,
		lastname,
		hashedPassword,
		now,
		now,
	)
	if err != nil {
		return err
	}

	return nil
}

/*
Signin is a function that gets a user from the database by username or email,
and checks if the passwords match.
If correct, it returns the user's id.
*/
func (s *userService) Signin(usernameOrEmail, password string) (id string, err error) {
	query, err := s.db.Query(
		`SELECT
			id,
			password
		FROM user
		WHERE username = ?
		OR email = ?`,
		usernameOrEmail,
		usernameOrEmail,
	)
	if err != nil {
		return "", err
	}

	var hash string

	for query.Next() {
		err = query.Scan(
			&id,
			&hash,
		)
		if err != nil {
			return "", err
		}
	}

	match := utils.CheckPasswordHash(hash, password)
	if !match {
		return "", errors.New("Invalid password")
	}

	return id, nil
}

/*
GetById is a function that gets a user from the database by username.
*/
func (s *userService) GetById(id string) (*User, error) {
	query, err := s.db.Query(
		"SELECT * FROM user WHERE id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}

	user := User{}

	var createdAtStr string
	var updatedAtStr string

	for query.Next() {
		err := query.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Fistname,
			&user.Lastname,
			&user.Password,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return &user, err
		}
	}

	createdAt, err := utils.ConvertToTime(createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAt, err := utils.ConvertToTime(updatedAtStr)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return &user, nil
}

/*
GetByUsername is a function that gets a user from the database by username.
*/
func (s *userService) GetByUsername(username string) (*User, error) {
	query, err := s.db.Query(
		"SELECT * FROM user WHERE username = ?",
		username,
	)
	if err != nil {
		return nil, err
	}

	user := User{}

	var createdAtStr string
	var updatedAtStr string

	for query.Next() {
		err := query.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Fistname,
			&user.Lastname,
			&user.Password,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return &user, err
		}
	}

	createdAt, err := utils.ConvertToTime(createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAt, err := utils.ConvertToTime(updatedAtStr)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return &user, nil
}

/*
GetByEmail is a function that gets a user from the database by email.
*/
func (s *userService) GetByEmail(email string) (*User, error) {
	query, err := s.db.Query("SELECT * FROM user WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	user := User{}

	var createdAtStr string
	var updatedAtStr string

	for query.Next() {
		err := query.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Fistname,
			&user.Lastname,
			&user.Password,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return &user, err
		}
	}

	createdAt, err := utils.ConvertToTime(createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAt, err := utils.ConvertToTime(updatedAtStr)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return &user, nil
}
