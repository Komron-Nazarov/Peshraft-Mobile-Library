package models

import "time"

type Notification struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    *uint     `json:"user_id"`     // Обязательно указатель
    Message   string    `json:"message"`     // С большой буквы
    IsRead    bool      `json:"is_read"`     // С большой буквы
    CreatedAt time.Time `json:"created_at"`
}
// type Notification struct {
// 	ID                   uint      `json:"id" gorm:"primaryKey"`
// 	UserID               *uint     `json:"user_id"` 
	
// 	// Старые поля, которые вызывали ошибку в старом коде
// 	Message              string    `json:"message"` 
// 	IsRead               bool      `json:"is_read" gorm:"default:false"`
	
// 	// Новые поля для админки
// 	NotificationType     string    `json:"notification_type"` 
// 	NotificationImageURL string    `json:"notification_image_url"`
// 	Title                string    `json:"title"`
// 	Description          string    `json:"description"`
// 	Date                 string    `json:"date"`
// 	Time                 string    `json:"time"`
// 	CreatedAt            time.Time `json:"created_at"`
// }

type NotificationResponse struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// func (n *Notification) ToResponse() NotificationResponse {
// 	return NotificationResponse{
// 		ID:        n.ID,
// 		Message:   n.Message,
// 		IsRead:    n.IsRead,
// 		CreatedAt: n.CreatedAt,
// 	}
// }

func (n *Notification) ToResponse() NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		Message:   n.Message,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}