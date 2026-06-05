package repositories

import "mobile-library/internal/models"

type FavoriteRepository struct {
	db *DB
}

func NewFavoriteRepository(db *DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Add(userID, bookID uint) error {
	_, err := r.db.conn.Exec(
		`INSERT INTO favorites (user_id, book_id, created_at) VALUES ($1, $2, NOW()) 
		 ON CONFLICT (user_id, book_id) DO NOTHING`, userID, bookID,
	)
	return err
}

func (r *FavoriteRepository) Remove(userID, bookID uint) error {
	_, err := r.db.conn.Exec(
		`DELETE FROM favorites WHERE user_id = $1 AND book_id = $2`, userID, bookID,
	)
	return err
}

func (r *FavoriteRepository) GetByUser(userID uint) ([]struct {
	Favorite models.Favorite
	Book     models.Book
}, error) {
	rows, err := r.db.conn.Query(
		`SELECT f.id, f.user_id, f.book_id, f.created_at,
				b.id, b.title, b.author, b.description, b.category, b.year, b.available_copies, b.image_url, b.created_at
		 FROM favorites f
		 JOIN books b ON f.book_id = b.id
		 WHERE f.user_id = $1
		 ORDER BY f.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Favorite models.Favorite
		Book     models.Book
	}
	for rows.Next() {
		var r struct {
			Favorite models.Favorite
			Book     models.Book
		}
		rows.Scan(
			&r.Favorite.ID, &r.Favorite.UserID, &r.Favorite.BookID, &r.Favorite.CreatedAt,
			&r.Book.ID, &r.Book.Title, &r.Book.Author, &r.Book.Description,
			&r.Book.Category, &r.Book.Year, &r.Book.AvailableCopies, &r.Book.ImageURL, &r.Book.CreatedAt,
		)
		results = append(results, r)
	}
	return results, nil
}

func (r *FavoriteRepository) IsFavorite(userID, bookID uint) (bool, error) {
	var exists bool
	err := r.db.conn.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND book_id = $2)`,
		userID, bookID,
	).Scan(&exists)
	return exists, err
}
