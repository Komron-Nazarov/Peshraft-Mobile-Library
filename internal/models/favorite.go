package models

import "time"

type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	BookID    uint      `json:"book_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

type FavoriteResponse struct {
	ID        uint      `json:"id"`
	BookID    uint      `json:"book_id"`
	Book      Book      `json:"book"`
	CreatedAt time.Time `json:"created_at"`
}
