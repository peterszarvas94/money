package services

import (
	"database/sql"
	"pengoe/internal/utils"
	"time"
)


type User struct {
	Id        int
	Username  string
	Email     string
	Fistname  string
	Lastname  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserServiceInterface interface {
	Signup(user *User) (*User, error)
	Signin(usernameOrEmail, password string) (*User, error)
	GetById(id int) (*User, error)
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
func (s *userService) Signup(user *User) (*User, error) {
	hashedPassword, hashErr := utils.HashPassword(user.Password)
	if hashErr != nil {
		return nil, hashErr
	}

	now := time.Now().UTC()

	mutation, mutationErr := s.db.Exec(
		`INSERT INTO user (
			username,
      email,
      firstname,
      lastname,
      password,
      created_at,
      updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Username,
		user.Email,
		user.Fistname,
		user.Lastname,
		hashedPassword,
		now,
		now,
	)
	if mutationErr != nil {
		return nil, mutationErr
	}

	id, idErr := mutation.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}

	newUser := &User{
		Id:        int(id),
		Username:  user.Username,
		Email:     user.Email,
		Fistname:  user.Fistname,
		Lastname:  user.Lastname,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return newUser, nil
}

/*
Signin is a function that gets a user from the database by username or email,
and checks if the passwords match.
If correct, it returns the user's id.
*/
func (s *userService) Signin(usernameOrEmail, password string) (*User, error) {
	query, queryErr := s.db.Query(
		`SELECT
      id,
      username,
      email,
      firstname,
      lastname,
      password as hash,
      created_at,
      updated_at
		FROM user
    WHERE username = ?
    OR email = ?`,
		usernameOrEmail,
		usernameOrEmail,
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
		scanErr := query.Scan(
			&id,
			&username,
			&email,
			&firstname,
			&lastname,
			&hash,
			&created_at,
			&updated_at,
		)
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

	created, createdErr := utils.ConvertToTime(created_at)
	if createdErr != nil {
		return nil, createdErr
	}

	updated, updatedErr := utils.ConvertToTime(updated_at)
	if updatedErr != nil {
		return nil, updatedErr
	}

	user := User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		CreatedAt: created,
		UpdatedAt: updated,
	}

	return &user, nil
}

/*
GetById is a function that gets a user from the database by username.
*/
func (s *userService) GetById(id int) (*User, error) {
	var user User

	query, queryErr := s.db.Query(
		"SELECT * FROM user WHERE id = ?",
		id,
	)
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
		scanErr := query.Scan(
			&id,
			&username,
			&email,
			&firstname,
			&lastname,
			&password,
			&created_at,
			&updated_at,
		)
		if scanErr != nil {
			return &user, scanErr
		}
	}

	if id == 0 {
		return &user, sql.ErrNoRows
	}

	created, createdErr := utils.ConvertToTime(created_at)
	if createdErr != nil {
		return nil, createdErr
	}

	updated, updatedErr := utils.ConvertToTime(updated_at)
	if updatedErr != nil {
		return nil, updatedErr
	}

	user = User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		CreatedAt: created,
		UpdatedAt: updated,
	}

	return &user, nil
}

/*
GetByUsername is a function that gets a user from the database by username.
*/
func (s *userService) GetByUsername(username string) (*User, error) {
	var user User

	query, queryErr := s.db.Query(
		"SELECT * FROM user WHERE username = ?",
		username,
	)
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
		scanErr := query.Scan(
			&id,
			&username,
			&email,
			&firstname,
			&lastname,
			&password,
			&created_at,
			&updated_at,
		)
		if scanErr != nil {
			return &user, scanErr
		}
	}

	if id == 0 {
		return &user, sql.ErrNoRows
	}

	created, createdErr := utils.ConvertToTime(created_at)
	if createdErr != nil {
		return nil, createdErr
	}

	updated, updatedErr := utils.ConvertToTime(updated_at)
	if updatedErr != nil {
		return nil, updatedErr
	}

	user = User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		CreatedAt: created,
		UpdatedAt: updated,
	}

	return &user, nil
}

/*
GetByEmail is a function that gets a user from the database by email.
*/
func (s *userService) GetByEmail(email string) (*User, error) {
	var user User

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
		scanErr := query.Scan(
			&id,
			&username,
			&email,
			&firstname,
			&lastname,
			&password,
			&created_at,
			&updated_at,
		)
		if scanErr != nil {
			return &user, scanErr
		}
	}

	if id == 0 {
		return &user, sql.ErrNoRows
	}

	created, createdErr := utils.ConvertToTime(created_at)
	if createdErr != nil {
		return nil, createdErr
	}

	updated, updatedErr := utils.ConvertToTime(updated_at)
	if updatedErr != nil {
		return nil, updatedErr
	}

	user = User{
		Id:        id,
		Username:  username,
		Email:     email,
		Fistname:  firstname,
		Lastname:  lastname,
		CreatedAt: created,
		UpdatedAt: updated,
	}

	return &user, nil
}
