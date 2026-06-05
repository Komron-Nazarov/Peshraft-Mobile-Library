package models

import "time"

type Book struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Title           string    `json:"title" gorm:"not null;index"`
	Author          string    `json:"author" gorm:"not null;index"`
	Description     string    `json:"description"`
	Category        string    `json:"category" gorm:"index"`
	Year            int       `json:"year"`
	PageCount       int       `json:"page_count"`      // Добавил для админки
    Language        string    `json:"language"`        // Добавил для админки
    AvailableCopies int       `json:"available_copies" gorm:"default:1"`
    Status          string    `json:"status" gorm:"default:'Available'"` // "Available", "Borrowed"
	ImageURL        string    `json:"image_url"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateBookRequest struct {
    Title           string `json:"title"`
    Author          string `json:"author"`
    Description     string `json:"description"`
    Category        string `json:"category"`
    Year            int    `json:"year"`
    AvailableCopies int    `json:"available_copies"`
    ImageURL        string `json:"image_url"`
    Language        string `json:"language"`
    PageCount       int    `json:"page_count"`
}

type BookSearchParams struct {
	Query    string `form:"q"`
	Author   string `form:"author"`
	Category string `form:"category"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

type BookResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Year            int       `json:"year"`
	AvailableCopies int       `json:"available_copies"`
	ImageURL        string    `json:"image_url"`
	CreatedAt       time.Time `json:"created_at"`
}

func (b *Book) ToResponse() BookResponse {
	return BookResponse{
		ID:              b.ID,
		Title:           b.Title,
		Author:          b.Author,
		Description:     b.Description,
		Category:        b.Category,
		Year:            b.Year,
		AvailableCopies: b.AvailableCopies,
		ImageURL:        b.ImageURL,
		CreatedAt:       b.CreatedAt,
	}
}
