package repositories

import "mobile-library/internal/models"

type AdminRepository struct {
    db *DB
}

func NewAdminRepository(db *DB) *AdminRepository {
    return &AdminRepository{db: db}
}

// GetStats — собирает общую статистику
func (r *AdminRepository) GetStats() (*models.DashboardStats, error) {
    stats := &models.DashboardStats{}
    // Пример запроса: считаем количество через UNION или отдельные запросы
    err := r.db.conn.QueryRow(`
        SELECT 
            (SELECT COUNT(*) FROM users) as total_members,
            (SELECT COUNT(*) FROM books) as total_books,
            (SELECT COUNT(*) FROM borrows WHERE status = 'active') as active_borrows,
            (SELECT COUNT(*) FROM borrows WHERE due_date < NOW() AND status = 'active') as overdue_books
    `).Scan(&stats.TotalMembers, &stats.TotalBooks, &stats.ActiveBorrows, &stats.OverdueBooks)
    
    return stats, err
}