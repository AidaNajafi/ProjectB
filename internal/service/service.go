package service

import (
	"authentication/internal/helpers"
	"authentication/internal/repository"
	"errors"
	"fmt"
)

type JwtUserInfo struct {
	ID    int
	Email string
}

type AuthService struct {
	repo repository.Store
}

func NewAuthService(r repository.Store) *AuthService {
	return &AuthService{repo: r}
}

func (a *AuthService) SignUp(user SignUpRequest) error {

	if err := a.repo.AlreadyExistCheck(user.Username, user.Email); err != nil {
		return err
	}
	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("Hashing Failed : %w", err)
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
		return fmt.Errorf("Creating Failed: %w", err)
	}
	return nil

}

func (a *AuthService) Login(userLog LoginRequest) (JwtUserInfo, error) {
	InUser, err := a.repo.GetUserByUsername(userLog.Username)
	if err != nil {
		return JwtUserInfo{}, fmt.Errorf("This user doesn't exist! :%w", err)
	}
	if !helpers.CheckPassword(userLog.Password, InUser.Password) {
		return JwtUserInfo{}, errors.New("Wrong password!")
	}
	jwtStr := JwtUserInfo{
		ID:    InUser.ID,
		Email: InUser.Email,
	}
	return jwtStr, nil
}
