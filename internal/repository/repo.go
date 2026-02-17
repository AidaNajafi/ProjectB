package repository

type UserInfo struct {
	ID       int
	Name     string
	Username string
	Email    string
	Password string
}

type Store interface {
	CreateUser(user UserInfo) error
	GetUserByUsername(username string) (*UserInfo, error)
}
