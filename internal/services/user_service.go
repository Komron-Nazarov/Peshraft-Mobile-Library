package services

import (
	"errors"
	"fmt"
	"mobile-library/internal/models"
	"mobile-library/internal/repositories"
	"mobile-library/pkg"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) UpdateProfile(id uint, req models.UpdateProfileRequest) (*models.UserResponse, error) {
	// Проверяем существование пользователя
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Валидация (минимум логики)
	if req.Name != "" {
		if len(req.Name) < 2 {
			return nil, errors.New("name is too short")
		}
		user.Name = req.Name
	}

	if req.Phone != "" {
		// Здесь можно добавить регулярку для проверки телефона
		user.Phone = req.Phone
	}

	// Сохраняем только измененные поля
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) ChangePassword(id uint, req models.ChangePasswordRequest) error {
	// 1. Проверка длины пароля
	if len(req.NewPassword) < 6 {
		return errors.New("new password must be at least 6 characters long")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// 2. Проверка старого пароля
	if !pkg.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("incorrect old password")
	}

	// 3. Проверка: не совпадает ли новый пароль со старым
	if req.OldPassword == req.NewPassword {
		return errors.New("new password cannot be the same as the old one")
	}

	// 4. Хеширование
	hash, err := pkg.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("security error: failed to process password")
	}

	// 5. Обновление в БД
	if err := s.userRepo.UpdatePassword(id, hash); err != nil {
		return fmt.Errorf("database error: could not save new password")
	}

	return nil
}
