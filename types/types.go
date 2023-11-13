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

type Page struct {
	Title         string
	Descrtipion   string
	Session       Session
	Data					map[string]string
}

// type SigninData struct {
// 	Page            Page
// 	UsernameOrEmail string
// 	Error           string
// }
//
// type Exists struct {
// 	Username bool
// 	Email    bool
// }
//
// type SignupData struct {
// 	Page   Page
// 	User   User
// 	Exists Exists
// 	Error  string
// }

type TokenVariant string

const (
	Access  TokenVariant = "access"
	Refresh TokenVariant = "refresh"
)

type JWT struct {
	Token   string
	Expires int64
}
