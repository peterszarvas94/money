package utils

// db
type User struct {
	Id        int
	Username  string
	Email     string
	Fistname  string
	Lastname  string
	Password  string
	CreatedAt string
	UpdatedAt string
}

type Account struct {
	Id          int
	Name        string
	Description string
	Currency    string
	CreatedAt   string
	UpdatedAt   string
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
	CreatedAt string
	UpdatedAt string
}

// page
type Session struct {
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
	Session     Session
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
