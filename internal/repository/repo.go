package repository

type Store interface {
	Init() error
	Create(user User) error
	GenerateID() (int, error)
	ReadAll() ([]User, error)
	GetUserByUsername(username string) (*User, error)
	AlreadyExistCheck(username, email string) error
}
