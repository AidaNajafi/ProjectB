package service

import (
	"authentication/internal/helpers"
	"authentication/internal/repository"
	"errors"
	"fmt"
)

type AuthService struct {
	repo   repository.Store
	jwtTok *JwtService
}

func NewAuthService(r repository.Store, jwtTok *JwtService) *AuthService {
	return &AuthService{repo: r, jwtTok: jwtTok}
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

	if err := a.repo.Create(repository.User{
		ID:       id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: hashed,
	}); err != nil {
		return fmt.Errorf("Creating Failed: %w", err)
	}
	return nil

}

type JwtUserInfo struct {
	ID    int
	Email string
}

func (a *AuthService) Login(userLog LoginRequest) (string, error) {
	InUser, err := a.repo.GetUserByUsername(userLog.Username)
	if err != nil {
		return "", fmt.Errorf("This user doesn't exist! :%w", err)
	}
	if !helpers.CheckPassword(userLog.Password, InUser.Password) {
		return "", errors.New("Wrong password!")
	}
	jwtStr := UserClaim{
		ID:    InUser.ID,
		Email: InUser.Email,
	}
	jwtToken, err := a.jwtTok.GenerateToken(jwtStr)

	return jwtToken, nil
}
