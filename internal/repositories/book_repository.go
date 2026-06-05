package repositories

import (
	"fmt"
	"mobile-library/internal/models"
	"strings"
)

type BookRepository struct {
	db *DB
}

func NewBookRepository(db *DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *models.Book) error {
	return r.db.conn.QueryRow(
		`INSERT INTO books (title, author, description, category, year, available_copies, image_url, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		 RETURNING id, created_at`,
		book.Title, book.Author, book.Description, book.Category,
		book.Year, book.AvailableCopies, book.ImageURL,
	).Scan(&book.ID, &book.CreatedAt)
}

func (r *BookRepository) FindAll(page, pageSize int) ([]models.Book, int, error) {
	offset := (page - 1) * pageSize
	var total int
	r.db.conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&total)

	rows, err := r.db.conn.Query(
		`SELECT id, title, author, description, category, year, available_copies, image_url, created_at
		 FROM books ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.Category,
			&b.Year, &b.AvailableCopies, &b.ImageURL, &b.CreatedAt)
		books = append(books, b)
	}
	return books, total, nil
}

func (r *BookRepository) FindByID(id uint) (*models.Book, error) {
	book := &models.Book{}
	err := r.db.conn.QueryRow(
		`SELECT id, title, author, description, category, year, available_copies, image_url, created_at
		 FROM books WHERE id = $1`, id,
	).Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Category,
		&book.Year, &book.AvailableCopies, &book.ImageURL, &book.CreatedAt)

	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) Search(params models.BookSearchParams) ([]models.Book, int, error) {
	var total int
	var conditions []string
	var args []interface{}
	argCount := 0

	if params.Query != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf(
			"(title ILIKE $%d OR author ILIKE $%d OR description ILIKE $%d)",
			argCount, argCount, argCount,
		))
		args = append(args, "%"+params.Query+"%")
	}
	if params.Category != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("category ILIKE $%d", argCount))
		args = append(args, "%"+params.Category+"%")
	}
	if params.Author != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("author ILIKE $%d", argCount))
		args = append(args, "%"+params.Author+"%")
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	r.db.conn.QueryRow(
		fmt.Sprintf("SELECT COUNT(*) FROM books %s", where),
		args...,
	).Scan(&total)

	offset := (params.Page - 1) * params.PageSize
	argCount++
	args = append(args, params.PageSize)
	argCount++
	args = append(args, offset)

	rows, err := r.db.conn.Query(
		fmt.Sprintf(
			`SELECT id, title, author, description, category, year, available_copies, image_url, created_at
			 FROM books %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
			where, argCount-1, argCount,
		),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description, &b.Category,
			&b.Year, &b.AvailableCopies, &b.ImageURL, &b.CreatedAt)
		books = append(books, b)
	}
	return books, total, nil
}

func (r *BookRepository) UpdateAvailableCopies(id uint, delta int) error {
	_, err := r.db.conn.Exec(
		`UPDATE books SET available_copies = available_copies + $1 WHERE id = $2 AND available_copies + $1 >= 0`,
		delta, id,
	)
	return err
}


func (r *BookRepository) Update(id uint, b *models.Book) error {
	// Используем твой метод r.db.Exec, который мы видели в db.go
	query := `
		UPDATE books 
		SET title = $1, author = $2, description = $3, category = $4, 
		    year = $5, available_copies = $6, image_url = $7, status = $8
		WHERE id = $9`
	
	return r.db.Exec(query, 
		b.Title, b.Author, b.Description, b.Category, 
		b.Year, b.AvailableCopies, b.ImageURL, b.Status, id)
}

func (r *BookRepository) Delete(id uint) error {
	query := `DELETE FROM books WHERE id = $1`
	return r.db.Exec(query, id)
}

func (r *BookRepository) GetCategories() ([]string, error) {
    // Предполагаю, что у тебя в репозитории есть доступ к объекту базы (например, r.db)
    // Если у тебя используется *sql.DB, то метод будет таким:
    
    query := "SELECT DISTINCT category FROM books WHERE category IS NOT NULL"
    rows, err := r.db.GetConn().Query(query) // Или r.db.GetConn().Query(...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []string
    for rows.Next() {
        var cat string
        if err := rows.Scan(&cat); err != nil {
            return nil, err
        }
        categories = append(categories, cat)
    }
    return categories, nil
}

func (r *BookRepository) CreateCategory(name string) error {
    return r.db.Exec("INSERT INTO categories (name) VALUES ($1)", name)
    
}

func (r *BookRepository) UpdateCategory(id, name string) error {
    return r.db.Exec("UPDATE categories SET name = $1 WHERE id = $2", name, id)
    
}

func (r *BookRepository) DeleteCategory(id string) error {
   return r.db.Exec("DELETE FROM categories WHERE id = $1", id)
   
}