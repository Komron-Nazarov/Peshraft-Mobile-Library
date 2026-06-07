// package models

// import "time"

// type Notification struct {
//     ID        uint      `json:"id" gorm:"primaryKey"`
//     UserID    *uint     `json:"user_id"`     // Обязательно указатель
//     Message   string    `json:"message"`     // С большой буквы
//     IsRead    bool      `json:"is_read"`     // С большой буквы
//     CreatedAt time.Time `json:"created_at"`
// }
// // type Notification struct {
// // 	ID                   uint      `json:"id" gorm:"primaryKey"`
// // 	UserID               *uint     `json:"user_id"` 
	
// // 	// Старые поля, которые вызывали ошибку в старом коде
// // 	Message              string    `json:"message"` 
// // 	IsRead               bool      `json:"is_read" gorm:"default:false"`
	
// // 	// Новые поля для админки
// // 	NotificationType     string    `json:"notification_type"` 
// // 	NotificationImageURL string    `json:"notification_image_url"`
// // 	Title                string    `json:"title"`
// // 	Description          string    `json:"description"`
// // 	Date                 string    `json:"date"`
// // 	Time                 string    `json:"time"`
// // 	CreatedAt            time.Time `json:"created_at"`
// // }

// type NotificationResponse struct {
// 	ID        uint      `json:"id"`
// 	Message   string    `json:"message"`
// 	IsRead    bool      `json:"is_read"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// // func (n *Notification) ToResponse() NotificationResponse {
// // 	return NotificationResponse{
// // 		ID:        n.ID,
// // 		Message:   n.Message,
// // 		IsRead:    n.IsRead,
// // 		CreatedAt: n.CreatedAt,
// // 	}
// // }

// func (n *Notification) ToResponse() NotificationResponse {
// 	return NotificationResponse{
// 		ID:        n.ID,
// 		Message:   n.Message,
// 		IsRead:    n.IsRead,
// 		CreatedAt: n.CreatedAt,
// 	}
// }






package models

import "time"

type Notification struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	UserID               *uint     `json:"user_id"` // Указатель, так как для общих уведомлений тут NULL
	Title                string    `json:"title"`   // Соответствует колонке title
	Message              string    `json:"message"`
	Type                 string    `json:"type"`                   // Соответствует колонке type ('news', 'warning', 'system', 'borrow')
	NotificationImageURL string    `json:"notification_image_url"` // Соответствует колонке notification_image_url
	IsRead               bool      `json:"is_read" gorm:"default:false"`
	CreatedAt            time.Time `json:"created_at"`
}

type NotificationResponse struct {
	ID                   uint      `json:"id"`
	Title                string    `json:"title"`
	Message              string    `json:"message"`
	Type                 string    `json:"type"`
	NotificationImageURL string    `json:"notification_image_url"`
	IsRead               bool      `json:"is_read"`
	CreatedAt            time.Time `json:"created_at"`
}

func (n *Notification) ToResponse() NotificationResponse {
	title := n.Title
	if title == "" {
		title = "Уведомление" // Дефолтный заголовок, если в базе пусто
	}
	notificationType := n.Type
	if notificationType == "" {
		notificationType = "news" // Дефолтный тип
	}

	return NotificationResponse{
		ID:                   n.ID,
		Title:                title,
		Message:              n.Message,
		Type:                 notificationType,
		NotificationImageURL: n.NotificationImageURL,
		IsRead:               n.IsRead,
		CreatedAt:            n.CreatedAt,
	}
}