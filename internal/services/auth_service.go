package services

import (
	"errors"
	"fmt"
	"mobile-library/internal/models"
	"mobile-library/internal/repositories"
	"mobile-library/pkg"
	"time"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	secret   string
}

func NewAuthService(userRepo *repositories.UserRepository, secret string) *AuthService {
	return &AuthService{userRepo: userRepo, secret: secret}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.UserResponse, string, error) {
	// Check if email exists
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, "", errors.New("email already registered")
	}

	hash, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// user := &models.User{
	// 	Name:     req.Name,
	// 	Email:    req.Email,
	// 	Phone:    req.Phone,
	// 	Password: hash,
	// }

	// if err := s.userRepo.Create(user); err != nil {
	// 	return nil, "", fmt.Errorf("failed to create user: %w", err)
	// }

	// token, err := pkg.GenerateToken(user.ID, user.Email, user.Name, s.secret, 24*time.Hour)

	user := &models.User{
    Name:        req.Name,
    Email:       req.Email,
    Phone:       req.Phone,
    Password:    hash,
    Role:        "user",            // Дефолтная роль
    JobPosition: "Student",         // Дефолтная позиция
}


if err := s.userRepo.Create(user); err != nil {
    return nil, "", fmt.Errorf("failed to create user: %w", err)
}

// 2. Передаем user.Role в обновленную функцию
token, err := pkg.GenerateToken(user.ID, user.Email, user.Name, user.Role, s.secret, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	response := user.ToResponse()

	return &response, token, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.UserResponse, string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !pkg.CheckPasswordHash(req.Password, user.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	// token, err := pkg.GenerateToken(user.ID, user.Email, user.Name, s.secret, 24*time.Hour)
	token, err := pkg.GenerateToken(user.ID, user.Email, user.Name, user.Role, s.secret, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	resp := user.ToResponse()

	return &resp, token, nil
}
