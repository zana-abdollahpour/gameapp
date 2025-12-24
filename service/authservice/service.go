package authservice

import (
	"fmt"
	"gameapp/entity"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	SignKey               string
	AccessSubject         string
	RefreshSubject        string
	AccessExpirationTime  time.Duration
	RefreshExpirationTime time.Duration
}

type Service struct {
	config Config
}

func New(
	cfg Config,
) Service {
	return Service{
		config: cfg,
	}
}

func (s Service) createToken(userID uint, subject string, expireDuration time.Duration) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		subject: subject,
		UserID:  userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", fmt.Errorf("unexpected error: %w", err)
	}

	return signedString, nil
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.AccessSubject, s.config.AccessExpirationTime)
}

func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.RefreshSubject, s.config.RefreshExpirationTime)

}

func (s Service) ParseToken(bearerToken string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		return []byte(s.config.SignKey), nil
	}

	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, keyFunc, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	} else if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unknown claims type, cannot proceed")
	}
}
