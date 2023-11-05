package utils

// Types used in multiple files

type UserData struct {
	Id       int
	Username string
	Email    string
}

type SessionData struct {
	LoggedIn bool
	User     UserData
}

type PageData struct {
	Session SessionData
	Title   string
	Descrtipion string
}

type SigninData struct {
	User  string
	Error string
	Title string
	Descrtipion string
}

type ExistsData struct {
	Username bool
	Email    bool
}

type SignupData struct {
	Firstname string
	Lastname string
	Username string
	Email    string
	Error    string
	Exists   ExistsData
	Title    string
	Descrtipion string
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
