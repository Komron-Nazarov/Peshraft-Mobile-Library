package repositories

import (
	"mobile-library/internal/models"
)

type ReviewRepository struct {
	db *DB
}

func NewReviewRepository(db *DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// Create сохраняет отзыв в базу
func (r *ReviewRepository) Create(review *models.Review) error {
	query := `
		INSERT INTO reviews (book_id, user_id, rating, review, review_category, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at`
	
	return r.db.conn.QueryRow(query, 
		review.BookID, review.UserID, review.Rating, review.Review, review.ReviewCategory,
	).Scan(&review.ID, &review.CreatedAt)
}

// GetByBookID возвращает все отзывы на книгу вместе с именами авторов
func (r *ReviewRepository) GetByBookID(bookID uint) ([]models.Review, error) {
	query := `
		SELECT 
			r.id, r.book_id, r.user_id, u.name as user_name, u.job_position as user_job,
			r.rating, r.review, r.review_category, r.created_at
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.book_id = $1
		ORDER BY r.created_at DESC`

	rows, err := r.db.conn.Query(query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var rev models.Review
		err := rows.Scan(
			&rev.ID, &rev.BookID, &rev.UserID, &rev.UserName, &rev.UserJob,
			&rev.Rating, &rev.Review, &rev.ReviewCategory, &rev.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}