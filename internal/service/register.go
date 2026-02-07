package service

import (
	"authentication/internal/helpers"
	"authentication/internal/repository"
	"authentication/models"
	"errors"
	"fmt"
	"strings"
)

type SignUpInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService struct {
	repo repository.Store
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

func (a *AuthService) Login(info LoginInfo) (models.User, error) {

	user, err := a.repo.GetUserByUsername(info.Username)
	if err != nil {
		return models.User{}, errors.New("Wrong username!")
	}

	user.Password = strings.TrimSpace(user.Password)
	if !helpers.CheckPassword(info.Password, user.Password) {
		fmt.Println(info.Password)
		fmt.Println(user.Password)
		return models.User{}, errors.New("Wrong password!")
	}

	return user, nil
}
