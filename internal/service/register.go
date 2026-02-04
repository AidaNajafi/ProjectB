package service

import (
	"authentication/internal/helpers"
	"authentication/models"
	"errors"
)

type SignUpInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInfo struct {
	Username string
	Password string
}

func SignUp(info SignUpInfo) (models.User, error) {
	pass, err := helpers.HashPassword(info.Password)
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		Name:     info.Name,
		Username: info.Username,
		Email:    info.Email,
		Password: pass}

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
