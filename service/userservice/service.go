package userservice

import (
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
}

type Service struct {
	repo Repository
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO: implement OTP

	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
		}

		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number is used before")
		}
	}

	if len(req.Name) <= 2 {
		return RegisterResponse{}, fmt.Errorf("name length should be at least 3 characters long")
	}

	// TODO: provide better password check with regexp
	if len(req.Password) < 7 {
		return RegisterResponse{}, fmt.Errorf("name length should be at least 8 characters long")
	}

	hashedPassword, err := getBcryptHash(req.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return RegisterResponse{}, fmt.Errorf("unexpected error: failure in registration process")
	}

	user := entity.User{ID: 0, Name: req.Name, PhoneNumber: req.PhoneNumber, Password: hashedPassword}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return RegisterResponse{User: createdUser}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct{}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if !exist {
		return LoginResponse{}, fmt.Errorf("user or password is incorrect")
	}

	hashedPasswordFromRequest, err := getBcryptHash(req.Password)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if user.Password != hashedPasswordFromRequest {
		return LoginResponse{}, fmt.Errorf("user or password is incorrect")
	}

	return LoginResponse{}, nil
}

func getBcryptHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
