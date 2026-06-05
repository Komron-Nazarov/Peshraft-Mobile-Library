package services

import (
	"mobile-library/internal/models"
	"mobile-library/internal/repositories"
)

type NotificationService struct {
	repo *repositories.NotificationRepository
}

func NewNotificationService(r *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: r}
}

func (s *NotificationService) GetUserNotifications(userID uint) ([]models.NotificationResponse, int, error) {
	return s.repo.GetByUserID(userID)
}

// Теперь принимает два аргумента: ID уведомления и ID владельца
func (s *NotificationService) MarkAsRead(id uint, userID uint) error {
	return s.repo.MarkAsRead(id, userID)
}
