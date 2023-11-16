package types

// db
type User struct {
	Id       int
	Username string
	Email    string
	Fistname string
	Lastname string
	Password string
	CreatedAt  string
	UpdatedAt  string
}

// page
type Session struct {
	LoggedIn bool
	User     User
}

type AccountSelectItem struct {
	Id       int
	Text     string
}

type Page struct {
	Title         string
	Descrtipion   string
	Session       Session
	Data					map[string]string
}

type TokenVariant string

const (
	Access  TokenVariant = "access"
	Refresh TokenVariant = "refresh"
)

type JWT struct {
	Token   string
	Expires int64
}
