package service

import (
	"authentication/internal/helpers"
	"authentication/internal/provider"
	"authentication/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type AuthService struct {
	repo     repository.Store
	jwtTok   *JwtService
	provider *provider.Provider
}

func NewAuthService(r repository.Store, jwtTok *JwtService) *AuthService {
	return &AuthService{repo: r, jwtTok: jwtTok}
}

func (a *AuthService) SignUp(user SignUpRequest) error {

	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("Hashing Failed : %w", err)
	}
	u, err := a.repo.GetUserByUsername(user.Username)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Failed to retrieve user %w", err)
	}
	if err := a.repo.CreateUser(repository.UserInfo{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		UserPhone: user.UserPhone,
		Password:  hashed,
	}); err != nil {

		log.Println(err)
		return fmt.Errorf("Creating Failed: %w", err)
	}

	if u != nil {
		return errors.New("user already exists")
	}

	return nil

}

type JwtUserInfo struct {
	ID    int
	Email string
}

func (a *AuthService) Login(userLog LoginRequest) (string, error) {
	inUser, err := a.repo.GetUserByUsername(userLog.Username)
	if err != nil {
		return "", fmt.Errorf("This user doesn't exist! :%w", err)
	}
	if inUser == nil {
		return "", fmt.Errorf("user with this username %s not found", userLog.Username)
	}

	if !helpers.CheckPassword(userLog.Password, inUser.Password) {
		return "", errors.New("Wrong password!")
	}
	jwtStr := UserClaim{
		ID:    inUser.ID,
		Email: inUser.Email,
	}
	jwtToken, err := a.jwtTok.GenerateToken(jwtStr)

	return jwtToken, nil
}
