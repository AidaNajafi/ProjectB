package repository

import "authentication/models"

type Store interface {
	Init() error
	Create(user models.User) error
	GenerateID() (int, error)
	ReadAll() ([]models.User, error)
	GetUserByUsername(username string) (models.User, error)
}
