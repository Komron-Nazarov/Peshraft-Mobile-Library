// package repositories

// import (
// 	"mobile-library/internal/models"
// )

// type NotificationRepository struct {
// 	db *DB
// }

// func NewNotificationRepository(db *DB) *NotificationRepository {
// 	return &NotificationRepository{db: db}
// }

// // Create сохраняет новое уведомление (нужно для BorrowService)
// // func (r *NotificationRepository) Create(n *models.Notification) error {
// // 	query := `INSERT INTO notifications (user_id, message, is_read) VALUES ($1, $2, $3) RETURNING id`
// // 	return r.db.conn.QueryRow(query, n.UserID, n.Message, n.IsRead).Scan(&n.ID)
// // }
// // func (r *NotificationRepository) Create(n *models.Notification) error {
// //     query := `INSERT INTO notifications (user_id, message, is_read, created_at) 
// //               VALUES ($1, $2, $3, $4) RETURNING id`
// //     return r.db.conn.QueryRow(query, n.UserID, n.Message, n.IsRead, n.CreatedAt).Scan(&n.ID)
// // }
// func (r *NotificationRepository) Create(n *models.Notification) error {
//     query := `INSERT INTO notifications (user_id, message, is_read, created_at) 
//               VALUES ($1, $2, $3, $4) RETURNING id`
//     // Используй n.UserID, n.Message, n.IsRead
//     return r.db.conn.QueryRow(query, n.UserID, n.Message, n.IsRead, n.CreatedAt).Scan(&n.ID)
// }


// // GetByUserID получает список уведомлений пользователя
// func (r *NotificationRepository) GetByUserID(userID uint) ([]models.NotificationResponse, int, error) {
// 	query := `SELECT id, message, is_read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
// 	rows, err := r.db.conn.Query(query, userID)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	defer rows.Close()

// 	var notifs []models.NotificationResponse
// 	unread := 0
// 	for rows.Next() {
// 		var n models.NotificationResponse
// 		var createdAt interface{}
// 		if err := rows.Scan(&n.ID, &n.Message, &n.IsRead, &createdAt); err != nil {
// 			return nil, 0, err
// 		}
// 		if !n.IsRead {
// 			unread++
// 		}
// 		notifs = append(notifs, n)
// 	}
// 	return notifs, unread, nil
// }

// // MarkAsRead отмечает уведомление прочитанным с проверкой владельца
// func (r *NotificationRepository) MarkAsRead(id uint, userID uint) error {
// 	_, err := r.db.conn.Exec("UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2", id, userID)
// 	return err
// }







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

// Create сохраняет новое уведомление (используется в BorrowService и Admin логике)
func (r *NotificationRepository) Create(n *models.Notification) error {
	// Подставляем дефолтные значения, если поля не переданы
	if n.Type == "" {
		n.Type = "news"
	}
	if n.Title == "" {
		n.Title = "Уведомление"
	}

	query := `
		INSERT INTO notifications (user_id, title, message, type, notification_image_url, is_read, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`

	return r.db.conn.QueryRow(query, 
		n.UserID, n.Title, n.Message, n.Type, n.NotificationImageURL, n.IsRead, n.CreatedAt,
	).Scan(&n.ID)
}

// GetByUserID получает список уведомлений пользователя (включая общие, где user_id IS NULL)
func (r *NotificationRepository) GetByUserID(userID uint) ([]models.NotificationResponse, int, error) {
	// Добавляем новые колонки в SELECT и проверяем общие уведомления (user_id IS NULL)
	query := `
		SELECT id, title, message, type, notification_image_url, is_read, created_at 
		FROM notifications 
		WHERE user_id = $1 OR user_id IS NULL 
		ORDER BY created_at DESC`

	rows, err := r.db.conn.Query(query, userID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifs []models.NotificationResponse
	unread := 0

	for rows.Next() {
		var n models.NotificationResponse
		// Сканируем все новые поля напрямую в структуру ответа для мобилки
		err := rows.Scan(
			&n.ID, 
			&n.Title, 
			&n.Message, 
			&n.Type, 
			&n.NotificationImageURL, 
			&n.IsRead, 
			&n.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Защита: если в базе лежат пустые строки, даем дефолты
		if n.Title == "" {
			n.Title = "Уведомление"
		}
		if n.Type == "" {
			n.Type = "news"
		}

		if !n.IsRead {
			unread++
		}
		notifs = append(notifs, n)
	}

	if notifs == nil {
		notifs = []models.NotificationResponse{} // Возвращаем пустой массив [] вместо null
	}

	return notifs, unread, nil
}

// MarkAsRead отмечает уведомление прочитанным с проверкой владельца
func (r *NotificationRepository) MarkAsRead(id uint, userID uint) error {
	// Добавляем OR user_id IS NULL на случай прочтения общего уведомления
	query := `UPDATE notifications SET is_read = true WHERE id = $1 AND (user_id = $2 OR user_id IS NULL)`
	_, err := r.db.conn.Exec(query, id, userID)
	return err
}