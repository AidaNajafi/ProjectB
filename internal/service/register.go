package service

import (
	"authentication/internal/helpers"
	"authentication/internal/repository"
	"authentication/models"
	"errors"
)

type SignUpInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService struct {
	repo *repository.CsvRepository
}

func NewAuthService(r *repository.CsvRepository) *AuthService {
	return &AuthService{repo: r}
}

type LoginInfo struct {
	Username string
	Password string
}

func (a *AuthService) SignUp(info SignUpInfo) (models.User, error) {
	pass, err := helpers.HashPassword(info.Password)
	if err != nil {
		return models.User{}, err
	}
	id, err := a.repo.GenerateID()
	user := models.User{
		ID:       id,
		Name:     info.Name,
		Username: info.Username,
		Email:    info.Email,
		Password: pass}

	if err := a.repo.Create(user); err != nil {
		return models.User{}, err
	}

	return user, nil

}

func Login(info LoginInfo) (bool, error) {

	var body models.User
	if info.Username == body.Username || info.Password == body.Password {
		return true, nil
	} else {
		return false, errors.New("Wrong credentials")
	}
}
