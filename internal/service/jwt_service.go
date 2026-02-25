package service

import (
	"errors"
	"fmt"
	"strconv"
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
	ID    int64
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

var (
	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New("invalid token")
)

type VerifiedUser struct {
	ID int64
}

func (j *JwtService) VerifyToken(tokenString string) (VerifiedUser, error) {

	if tokenString == "" {
		return VerifiedUser{}, ErrMissingToken
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w, unexpected signing method: %v", ErrInvalidToken, t.Header["alg"])
		}
		return j.JwtSecret, nil
	})

	if err != nil {
		return VerifiedUser{}, fmt.Errorf("%w : %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return VerifiedUser{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return VerifiedUser{}, ErrInvalidToken
	}
	rawId, ok := claims["user_id"]
	if !ok {
		return VerifiedUser{}, fmt.Errorf("%w missin user id", ErrInvalidToken)
	}

	var id int64
	switch v := rawId.(type) {
	case float64:
		id = int64(v)
	case int64:
		id = v
	case int:
		id = int64(v)
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return VerifiedUser{}, fmt.Errorf("invalid user id, %w", ErrInvalidToken)
		}
		id = n
	default:
		return VerifiedUser{}, fmt.Errorf("invalid user id type, %w", ErrInvalidToken)
	}

	if id <= 0 {
		return VerifiedUser{}, fmt.Errorf("%w: invalid user id", ErrInvalidToken)
	}
	return VerifiedUser{ID: id}, nil

}
