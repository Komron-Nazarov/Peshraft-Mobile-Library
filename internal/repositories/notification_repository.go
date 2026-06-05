package repositories

import (
	"mobile-library/internal/models"
)

type NotificationRepository struct {
	db *DB
}

func NewNotificationRepository(db *DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create сохраняет новое уведомление (нужно для BorrowService)
// func (r *NotificationRepository) Create(n *models.Notification) error {
// 	query := `INSERT INTO notifications (user_id, message, is_read) VALUES ($1, $2, $3) RETURNING id`
// 	return r.db.conn.QueryRow(query, n.UserID, n.Message, n.IsRead).Scan(&n.ID)
// }
// func (r *NotificationRepository) Create(n *models.Notification) error {
//     query := `INSERT INTO notifications (user_id, message, is_read, created_at) 
//               VALUES ($1, $2, $3, $4) RETURNING id`
//     return r.db.conn.QueryRow(query, n.UserID, n.Message, n.IsRead, n.CreatedAt).Scan(&n.ID)
// }
func (r *NotificationRepository) Create(n *models.Notification) error {
    query := `INSERT INTO notifications (user_id, message, is_read, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
    // Используй n.UserID, n.Message, n.IsRead
    return r.db.conn.QueryRow(query, n.UserID, n.Message, n.IsRead, n.CreatedAt).Scan(&n.ID)
}


// GetByUserID получает список уведомлений пользователя
func (r *NotificationRepository) GetByUserID(userID uint) ([]models.NotificationResponse, int, error) {
	query := `SELECT id, message, is_read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.conn.Query(query, userID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifs []models.NotificationResponse
	unread := 0
	for rows.Next() {
		var n models.NotificationResponse
		var createdAt interface{}
		if err := rows.Scan(&n.ID, &n.Message, &n.IsRead, &createdAt); err != nil {
			return nil, 0, err
		}
		if !n.IsRead {
			unread++
		}
		notifs = append(notifs, n)
	}
	return notifs, unread, nil
}

// MarkAsRead отмечает уведомление прочитанным с проверкой владельца
func (r *NotificationRepository) MarkAsRead(id uint, userID uint) error {
	_, err := r.db.conn.Exec("UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
