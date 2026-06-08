// package models

// import "time"

// // type Book struct {
// // 	ID              uint      `json:"id" gorm:"primaryKey"`
// // 	Title           string    `json:"title" gorm:"not null;index"`
// // 	Author          string    `json:"author" gorm:"not null;index"`
// // 	Description     string    `json:"description"`
// // 	Category        string    `json:"category" gorm:"index"`
// // 	Year            int       `json:"year"`
// // 	PageCount       int       `json:"page_count"`      // Добавил для админки
// //     Language        string    `json:"language"`        // Добавил для админки
// //     AvailableCopies int       `json:"available_copies" gorm:"default:1"`
// //     Status          string    `json:"status" gorm:"default:'Available'"` // "Available", "Borrowed"
// // 	ImageURL        string    `json:"image_url"`
// // 	CreatedAt       time.Time `json:"created_at"`
// // }

// type Book struct {
// 	ID              uint      `json:"id" gorm:"primaryKey"`
// 	Title           string    `json:"title" gorm:"not null"`
// 	Author          string    `json:"author" gorm:"not null"`
// 	Description     string    `json:"description"`
// 	Category        string    `json:"category"`
// 	Year            int       `json:"year"`
// 	AvailableCopies int       `json:"available_copies" gorm:"default:0"`
// 	ImageURL        string    `json:"image_url"`
// 	BgImageURL      string    `json:"bg_image_url"` // Новое поле для заднего фона
// 	PageCount       int       `json:"page_count" gorm:"default:0"`
// 	Language        string    `json:"language" gorm:"default:'English'"`
// 	Status          string    `json:"status" gorm:"default:'Available'"`
// 	CreatedAt       time.Time `json:"created_at"`

// 	// Динамические поля, которые требует фронтенд (они не хранятся напрямую в таблице books, а вычисляются при GET запросах)
// 	Rating         float64 `json:"rating" gorm:"-"`           // Средний рейтинг из таблицы reviews
// 	Readers        int     `json:"readers" gorm:"-"`          // Сколько раз книгу брали из таблицы borrows
// 	IsFavoriteBook bool    `json:"isFavoriteBook" gorm:"-"`   // Лайкнул ли её текущий юзер
// }

// type CreateBookRequest struct {
//     Title           string `json:"title"`
//     Author          string `json:"author"`
//     Description     string `json:"description"`
//     Category        string `json:"category"`
//     Year            int    `json:"year"`
//     AvailableCopies int    `json:"available_copies"`
//     ImageURL        string `json:"image_url"`
// 	BgImageURL      string `json:"bg_image_url"`
//     Language        string `json:"language"`
//     PageCount       int    `json:"page_count"`
// }

// type BookSearchParams struct {
// 	Query    string `form:"q"`
// 	Author   string `form:"author"`
// 	Category string `form:"category"`
// 	Page     int    `form:"page,default=1"`
// 	PageSize int    `form:"page_size,default=20"`
// }

// type BookResponse struct {
// 	ID              uint      `json:"id"`
// 	Title           string    `json:"title"`
// 	Author          string    `json:"author"`
// 	Description     string    `json:"description"`
// 	Category        string    `json:"category"`
// 	Year            int       `json:"year"`
// 	AvailableCopies int       `json:"available_copies"`
// 	ImageURL        string    `json:"image_url"`
// 	CreatedAt       time.Time `json:"created_at"`
// }

// func (b *Book) ToResponse() BookResponse {
// 	return BookResponse{
// 		ID:              b.ID,
// 		Title:           b.Title,
// 		Author:          b.Author,
// 		Description:     b.Description,
// 		Category:        b.Category,
// 		Year:            b.Year,
// 		AvailableCopies: b.AvailableCopies,
// 		ImageURL:        b.ImageURL,
// 		CreatedAt:       b.CreatedAt,
// 	}
// }


// type Review struct {
// 	ID             uint      `json:"id"`
// 	BookID         uint      `json:"book_id"`
// 	UserID         uint      `json:"user_id"`
// 	UserName       string    `json:"user_name"`    // Вытащим через JOIN для мобилки
// 	UserJob        string    `json:"user_job"`     // Чтобы на фронте было видно, кто пишет (н-р, Student)
// 	Rating         int       `json:"rating"`
// 	Review         string    `json:"review"`
// 	ReviewCategory string    `json:"review_category"`
// 	CreatedAt      time.Time `json:"created_at"`
// }

// type CreateReviewRequest struct {
// 	Rating         int    `json:"rating" binding:"required,min=1,max=5"`
// 	Review         string `json:"review"`
// 	ReviewCategory string `json:"review_category"`
// }




package models

import "time"

type Book struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Title           string    `json:"title" gorm:"not null"`
	Author          string    `json:"author" gorm:"not null"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Year            int       `json:"year"`
	AvailableCopies int       `json:"available_copies" gorm:"default:0"`
	ImageURL        string    `json:"image_url"`
	BgImageURL      string    `json:"bg_image_url"` // Поле для заднего фона
	PageCount       int       `json:"page_count" gorm:"default:0"`
	Language        string    `json:"language" gorm:"default:'English'"`
	
	// ДОБАВЛЯЕМ omitempty: если статус пустой, он вообще удалится из JSON и не сломает Select на фронте
	Status          string    `json:"status,omitempty" gorm:"default:'Available'"`
	CreatedAt       time.Time `json:"created_at"`

	// Динамические поля для фронтенда (вычисляются при GET запросах)
	Rating         float64 `json:"rating" gorm:"-"`           // Средний рейтинг из таблицы reviews
	Readers        int     `json:"readers" gorm:"-"`          // Сколько раз книгу брали из таблицы borrows
	IsFavoriteBook bool    `json:"isFavoriteBook" gorm:"-"`   // Лайкнул ли её текущий юзер
}

type CreateBookRequest struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	Year            int    `json:"year"`
	AvailableCopies int    `json:"available_copies"`
	ImageURL        string `json:"image_url"`
	BgImageURL      string `json:"bg_image_url"`
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

type Review struct {
	ID             uint      `json:"id"`
	BookID         uint      `json:"book_id"`
	UserID         uint      `json:"user_id"`
	UserName       string    `json:"user_name"`    // Вытащим через JOIN для мобилки
	UserJob        string    `json:"user_job"`     // Чтобы на фронте было видно, кто пишет (н-р, Student)
	Rating         int       `json:"rating"`
	Review         string    `json:"review"`
	ReviewCategory string    `json:"review_category"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateReviewRequest struct {
	Rating         int    `json:"rating" binding:"required,min=1,max=5"`
	Review         string `json:"review"`
	ReviewCategory string `json:"review_category"`
}