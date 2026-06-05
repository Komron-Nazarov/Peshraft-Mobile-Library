package repositories

import (
	"mobile-library/internal/models"
	"time"
)

type BorrowRepository struct {
	db *DB
}

func NewBorrowRepository(db *DB) *BorrowRepository {
	return &BorrowRepository{db: db}
}

func (r *BorrowRepository) Create(borrow *models.Borrow) error {
	return r.db.conn.QueryRow(
		`INSERT INTO borrows (user_id, book_id, borrow_date, due_date, status, created_at)
         VALUES ($1, $2, $3, $4, $5, NOW())
         RETURNING id, created_at`,
		borrow.UserID, borrow.BookID, borrow.BorrowDate, borrow.DueDate, borrow.Status,
	).Scan(&borrow.ID, &borrow.CreatedAt)
}

func (r *BorrowRepository) FindByUserAndBook(userID, bookID uint) (*models.Borrow, error) {
	b := &models.Borrow{}
	err := r.db.conn.QueryRow(
		`SELECT id, user_id, book_id, borrow_date, due_date, return_date, status, created_at
         FROM borrows WHERE user_id = $1 AND book_id = $2 AND status = 'active' LIMIT 1`,
		userID, bookID,
	).Scan(&b.ID, &b.UserID, &b.BookID, &b.BorrowDate, &b.DueDate, &b.ReturnDate, &b.Status, &b.CreatedAt)

	if err != nil {
		return nil, err
	}
	return b, nil
}

// Переименован для совместимости с сервисом и расширен статусом
func (r *BorrowRepository) UpdateReturnStatus(borrowID uint, returnDate time.Time, status string) error {
	_, err := r.db.conn.Exec(
		`UPDATE borrows SET return_date = $1, status = $2 WHERE id = $3`,
		returnDate, status, borrowID,
	)
	return err
}

func (r *BorrowRepository) GetUserHistory(userID uint) ([]struct {
	Borrow models.Borrow
	Book   models.Book
}, error) {
	rows, err := r.db.conn.Query(
		`SELECT b.id, b.user_id, b.book_id, b.borrow_date, b.due_date, b.return_date, b.status, b.created_at,
                bk.id, bk.title, bk.author, bk.description, bk.category, bk.year, bk.available_copies, bk.image_url, bk.created_at
         FROM borrows b
         JOIN books bk ON b.book_id = bk.id
         WHERE b.user_id = $1
         ORDER BY b.borrow_date DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Borrow models.Borrow
		Book   models.Book
	}
	for rows.Next() {
		var res struct {
			Borrow models.Borrow
			Book   models.Book
		}
		// Важно: Scan требует указатели на каждое поле
		err := rows.Scan(
			&res.Borrow.ID, &res.Borrow.UserID, &res.Borrow.BookID, &res.Borrow.BorrowDate,
			&res.Borrow.DueDate, &res.Borrow.ReturnDate, &res.Borrow.Status, &res.Borrow.CreatedAt,
			&res.Book.ID, &res.Book.Title, &res.Book.Author, &res.Book.Description,
			&res.Book.Category, &res.Book.Year, &res.Book.AvailableCopies, &res.Book.ImageURL, &res.Book.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

// НОВЫЙ МЕТОД: Поиск просроченных книг для сервиса
func (r *BorrowRepository) GetOverdueBorrows() ([]models.OverdueBorrow, error) {
	query := `
		SELECT b.id, b.user_id, bk.title 
		FROM borrows b
		JOIN books bk ON b.book_id = bk.id
		WHERE b.status = 'active' AND b.due_date < NOW()
	`
	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overdue []models.OverdueBorrow
	for rows.Next() {
		var ob models.OverdueBorrow
		if err := rows.Scan(&ob.ID, &ob.UserID, &ob.Title); err != nil {
			return nil, err
		}
		overdue = append(overdue, ob)
	}
	return overdue, nil
}

// НОВЫЙ МЕТОД: Пометка просрочки
func (r *BorrowRepository) MarkOverdue(id uint) error {
	_, err := r.db.conn.Exec(`UPDATE borrows SET status = 'overdue' WHERE id = $1`, id)
	return err
}

func (r *BorrowRepository) GetActiveBorrows(userID uint) ([]models.Borrow, error) {
	rows, err := r.db.conn.Query(
		`SELECT id, user_id, book_id, borrow_date, due_date, return_date, status, created_at
         FROM borrows WHERE user_id = $1 AND status = 'active'`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrows []models.Borrow
	for rows.Next() {
		var b models.Borrow
		rows.Scan(&b.ID, &b.UserID, &b.BookID, &b.BorrowDate, &b.DueDate, &b.ReturnDate, &b.Status, &b.CreatedAt)
		borrows = append(borrows, b)
	}
	return borrows, nil
}
