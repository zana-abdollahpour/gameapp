package authservice

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	subject string
	UserID  uint `json:"user_id"`
}

func (c Claims) Valid() error {
	validator := jwt.NewValidator(jwt.WithLeeway(0))

	if err := validator.Validate(c.RegisteredClaims); err != nil {
		return err
	}

	if c.UserID == 0 {
		return fmt.Errorf("invalid user ID")
	}

	return nil
}
