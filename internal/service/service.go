package service

import (
	"authentication/internal/helpers"
	"authentication/internal/repository"
	"errors"
	"fmt"
	"strings"
)

type AuthService struct {
	repo repository.Store
}

func NewAuthService(r repository.Store) *AuthService {
	return &AuthService{repo: r}
}

func (a *AuthService) SignUp(user SignUpRequest) error {
	
	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}

	id, err := a.repo.GenerateID()
	if err != nil {
		return fmt.Errorf("Failed to generate id: %w", err)
	}
	repoUser := repository.User{
		ID:       id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: hashed,
	}
	if err := a.repo.Create(repoUser); err != nil {
		return err
	}

	return nil

}

func (a *AuthService) Login(user LoginRequest) (int, error) {

	InUser, err := a.repo.GetUserByUsername(user.Username)
	if err != nil {
		return 0, fmt.Errorf("Failed to filter by Username! :%w", err)
	}

	user.Password = strings.TrimSpace(user.Password)
	if !helpers.CheckPassword(InUser.Password, user.Password) {
		return 0, errors.New("Wrong password!")
	}

	return InUser.ID, nil
}
