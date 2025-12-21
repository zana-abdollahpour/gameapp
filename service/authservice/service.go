package authservice

import (
	"fmt"
	"gameapp/entity"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	signKey               string
	accessSubject         string
	refreshSubject        string
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
}

func New(
	signKey,
	accessSubject,
	refreshSubject string,
	accessExpirationTime,
	refreshExpirationTime time.Duration,
) Service {
	return Service{
		signKey:               signKey,
		accessSubject:         accessSubject,
		refreshSubject:        refreshSubject,
		accessExpirationTime:  accessExpirationTime,
		refreshExpirationTime: refreshExpirationTime,
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
	signedString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", fmt.Errorf("unexpected error: %w", err)
	}

	return signedString, nil
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.accessSubject, s.accessExpirationTime)
}

func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.refreshSubject, s.refreshExpirationTime)

}

func (s Service) ParseToken(bearerToken string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		return []byte(s.signKey), nil
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
