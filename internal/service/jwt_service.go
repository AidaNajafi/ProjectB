package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	JwtSecret []byte
}

func NewJwtService(j []byte) *JwtService {
	return &JwtService{JwtSecret: j}
}

type UserClaim struct {
	ID    int
	Email string
}

func (j *JwtService) GenerateToken(claim UserClaim) (string, error) {

	claims := jwt.MapClaims{
		"user_id": claim.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
