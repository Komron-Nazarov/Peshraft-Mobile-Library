package models

type DashboardStats struct {
	TotalMembers  int `json:"total_members"`
	TotalBooks    int `json:"total_books"`
	ActiveBorrows int `json:"active_borrows"`
	OverdueBooks  int `json:"overdue_books"`
}

type ChartData struct {
	Overdue  int    `json:"overdue"`
	Borrowed int    `json:"borrowed"`
	Date     string `json:"date"`
	Month    int    `json:"month"`
}

type OverdueMember struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	BookTitle   string `json:"book_title"`
	BorrowDate  string `json:"borrow_date"`
	DueDate     string `json:"due_date"`
	DaysOverdue int    `json:"days_overdue"`
}

type UserBookshelfItem struct {
	ID         string `json:"id"`
	ImageURL   string `json:"image_url"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	BorrowDate string `json:"borrow_date"`
	DueDate    string `json:"due_date"`
	Status     string `json:"status"` // Поле для разделения на "активные" и "историю"
}

type BookRequest struct {
	ID             string `json:"id"`
	MemberImageURL string `json:"member_image_url"`
	ReceiverName   string `json:"receiver_name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	RequestDate    string `json:"request_date"`
	BorrowDate     string `json:"borrow_date"`
	DueDate        string `json:"due_date"`
	BookTitle      string `json:"book_title"`
	Author         string `json:"author"`
}

