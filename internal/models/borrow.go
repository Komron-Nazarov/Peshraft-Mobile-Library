package models

import "time"

type BorrowStatus string

const (
	BorrowPending  BorrowStatus = "pending" // Юзер нажал "хочу книгу", ждет одобрения
	BorrowActive   BorrowStatus = "active"
	BorrowReturned BorrowStatus = "returned"
	BorrowOverdue  BorrowStatus = "overdue"
)

type Borrow struct {
	ID         uint         `json:"id" gorm:"primaryKey"`
    UserID     uint         `json:"user_id"`
    User       User         `json:"user" gorm:"foreignKey:UserID"` // Связь с юзером
    BookID     uint         `json:"book_id"`
    Book       Book         `json:"book" gorm:"foreignKey:BookID"` // Связь с книгой
    BorrowDate time.Time    `json:"borrow_date"`
    DueDate    time.Time    `json:"due_date"`
    ReturnDate *time.Time   `json:"return_date"`
    Status     BorrowStatus `json:"status" gorm:"default:'pending'"`
    CreatedAt  time.Time    `json:"created_at"`
}

type OverdueBorrow struct {
	ID     uint
	UserID uint
	Title  string
}

type BorrowResponse struct {
	ID        uint         `json:"id"`
	BookTitle string       `json:"book_title"`
	DueDate   string       `json:"due_date"`
	Status    BorrowStatus `json:"status"`
}

type BorrowHistoryRecord struct {
	Borrow Borrow
	Book   Book
}

func (b *Borrow) ToResponse(title string) BorrowResponse {
	return BorrowResponse{
		ID: b.ID, BookTitle: title,
		DueDate: b.DueDate.Format("2006-01-02"),
		Status:  b.Status,
	}
}
