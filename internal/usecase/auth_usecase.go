package usecase

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthUsecase interface {
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type AuthService struct {
	secretKey []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{secretKey: []byte(secret)}
}

func (a *AuthService) GenerateToken(userID int, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Токен на 24 часа

	return token.SignedString(a.secretKey)
}

func (a *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return a.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
