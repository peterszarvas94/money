package utils

import "time"

// db
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

type Account struct {
	Id          int
	Name        string
	Description string
	Currency    string
	CreatedAt   time.Time
	UpdatedAt	  time.Time
}

type Role string

const (
	Admin Role = "admin"
	Viewer Role = "viewer"
)

type Access struct {
	Id        int
	Role			Role
	UserId    int
	AccountId int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	Id         int
	UserId     int
	ValidUntil time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// page
type CurrentUser struct {
	LoggedIn bool
	User     User
}

type AccountSelectItem struct {
	Id   int
	Text string
}

type Page struct {
	Title       string
	Descrtipion string
	Session     CurrentUser
	Data        map[string]string
}

type TokenVariant string

const (
	AccessToken  TokenVariant = "access"
	RefreshToken TokenVariant = "refresh"
)

type JWT struct {
	Token   string
	Expires int64
}
