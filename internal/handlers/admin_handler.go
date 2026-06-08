// // // // package handlers

// // // // import (
// // // // 	"net/http"
// // // // 	"mobile-library/internal/repositories"
// // // // 	"github.com/gin-gonic/gin"
// // // // )

// // // // type AdminHandler struct {
// // // // 	db *repositories.DB
// // // // }

// // // // func NewAdminHandler(db *repositories.DB) *AdminHandler {
// // // // 	return &AdminHandler{db: db}
// // // // }

// // // // // GetStats возвращает данные для карточек дашборда
// // // // func (h *AdminHandler) GetStats(c *gin.Context) {
// // // // 	var totalBooks int64
// // // // 	var totalMembers int64
// // // // 	var activeBorrows int64
// // // // 	var overdueBooks int64

// // // // 	// 1. Считаем книги
// // // // 	h.db.Model(&models.Book{}).Count(&totalBooks)

// // // // 	// 2. Считаем пользователей
// // // // 	h.db.Model(&models.User{}).Count(&totalMembers)

// // // // 	// 3. Считаем активные аренды
// // // // 	h.db.Model(&models.Borrow{}).Where("status = ?", "active").Count(&activeBorrows)

// // // // 	// 4. Считаем просроченные (упрощенно: статус overdue или дата возврата прошла)
// // // // 	h.db.Model(&models.Borrow{}).Where("status = ? OR (status = ? AND due_date < NOW())", "overdue", "active").Count(&overdueBooks)

// // // // 	c.JSON(http.StatusOK, gin.H{
// // // // 		"total_books":      totalBooks,
// // // // 		"total_members":    totalMembers,
// // // // 		"active_borrows":   activeBorrows,
// // // // 		"overdue_books":    overdueBooks,
// // // // 		"revenue":          0, // Если у тебя библиотека бесплатная, можно оставить 0
// // // // 	})
// // // // }

// // // // // GetAllMembers возвращает список всех пользователей для таблицы
// // // // func (h *AdminHandler) GetAllMembers(c *gin.Context) {
// // // // 	var users []models.User
// // // // 	if err := h.db.Find(&users).Error; err != nil {
// // // // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список пользователей"})
// // // // 		return
// // // // 	}

// // // // 	var response []models.UserResponse
// // // // 	for _, u := range users {
// // // // 		response = append(response, u.ToResponse())
// // // // 	}

// // // // 	c.JSON(http.StatusOK, gin.H{"members": response})
// // // // }







// // // package handlers

// // // import (
// // // 	"net/http"
// // // 	"mobile-library/internal/models" // Импорт моделей обязателен
// // // 	"mobile-library/internal/repositories"
// // // 	"github.com/gin-gonic/gin"
// // // )

// // // type AdminHandler struct {
// // // 	db *repositories.DB
// // // }

// // // func NewAdminHandler(db *repositories.DB) *AdminHandler {
// // // 	return &AdminHandler{db: db}
// // // }

// // // // GetStats возвращает данные для карточек дашборда веб-сайта
// // // func (h *AdminHandler) GetStats(c *gin.Context) {
// // // 	var totalBooks int64
// // // 	var totalMembers int64
// // // 	var activeBorrows int64
// // // 	var overdueBooks int64

// // // 	// Используем GORM методы через h.db.DB (если в репозитории DB это *gorm.DB)
// // // 	h.db.Model(&models.Book{}).Count(&totalBooks)
// // // 	h.db.Model(&models.User{}).Count(&totalMembers)
// // // 	h.db.Model(&models.Borrow{}).Where("status = ?", "active").Count(&activeBorrows)
// // // 	h.db.Model(&models.Borrow{}).Where("status = ? OR (status = ? AND due_date < NOW())", "overdue", "active").Count(&overdueBooks)

// // // 	c.JSON(http.StatusOK, gin.H{
// // // 		"total_books":    totalBooks,
// // // 		"total_members":  totalMembers,
// // // 		"active_borrows": activeBorrows,
// // // 		"overdue_books":  overdueBooks,
// // // 		"revenue":        0, 
// // // 	})
// // // }

// // // // GetAllMembers возвращает список всех пользователей
// // // func (h *AdminHandler) GetAllMembers(c *gin.Context) {
// // // 	var users []models.User
// // // 	if err := h.db.Find(&users).Error; err != nil {
// // // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список пользователей"})
// // // 		return
// // // 	}

// // // 	// Предполагаем, что у модели User есть метод ToResponse
// // // 	var response []models.UserResponse
// // // 	for _, u := range users {
// // // 		response = append(response, u.ToResponse())
// // // 	}

// // // 	c.JSON(http.StatusOK, gin.H{"members": response})
// // // }

// // // // AcceptAdminRequest одобряет роль администратора для пользователя
// // // func (h *AdminHandler) AcceptAdminRequest(c *gin.Context) {
// // // 	id := c.Param("id")
	
// // // 	if err := h.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
// // // 		"role":             "admin",
// // // 		"is_pending_admin": false,
// // // 	}).Error; err != nil {
// // // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении роли"})
// // // 		return
// // // 	}

// // // 	c.JSON(http.StatusOK, gin.H{"message": "Пользователь теперь администратор"})
// // // }

















// // package handlers

// // import (
// // 	"net/http"
// // 	"mobile-library/internal/models"
// // 	"mobile-library/internal/repositories"
// // 	"github.com/gin-gonic/gin"
// // )

// // type AdminHandler struct {
// // 	db *repositories.DB
// // }

// // func NewAdminHandler(db *repositories.DB) *AdminHandler {
// // 	return &AdminHandler{db: db}
// // }

// // // GetStats возвращает данные для дашборда через чистый SQL
// // func (h *AdminHandler) GetStats(c *gin.Context) {
// // 	var totalBooks, totalMembers, activeBorrows, overdueBooks int

// // 	conn := h.db.GetConn()

// // 	// 1. Считаем книги
// // 	conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&totalBooks)

// // 	// 2. Считаем пользователей
// // 	conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalMembers)

// // 	// 3. Считаем активные аренды
// // 	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'active'").Scan(&activeBorrows)

// // 	// 4. Считаем просроченные
// // 	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'overdue' OR (status = 'active' AND due_date < NOW())").Scan(&overdueBooks)

// // 	c.JSON(http.StatusOK, gin.H{
// // 		"total_books":    totalBooks,
// // 		"total_members":  totalMembers,
// // 		"active_borrows": activeBorrows,
// // 		"overdue_books":  overdueBooks,
// // 		"revenue":        0,
// // 	})

// // 	// Внутри GetStats для графиков
// // rows, _ := conn.Query(`
// //     SELECT TO_CHAR(borrow_date, 'Mon'), COUNT(*) 
// //     FROM borrows 
// //     GROUP BY TO_CHAR(borrow_date, 'Mon')
// // `)
// // // ... логика упаковки в JSON для фронта
// // }

// // // GetAllMembers возвращает список всех пользователей
// // func (h *AdminHandler) GetAllMembers(c *gin.Context) {
// // 	rows, err := h.db.GetConn().Query("SELECT id, name, email, phone, role, job_position, date_of_birth, created_at FROM users")
// // 	if err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
// // 		return
// // 	}
// // 	defer rows.Close()

// // 	var response []models.UserResponse
// // 	for rows.Next() {
// // 		var u models.User
// // 		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.JobPosition, &u.DateOfBirth, &u.CreatedAt); err != nil {
// // 			continue
// // 		}
// // 		response = append(response, u.ToResponse())
// // 	}

// // 	c.JSON(http.StatusOK, gin.H{"members": response})
// // }

// // // AcceptAdminRequest подтверждает права админа
// // func (h *AdminHandler) AcceptAdminRequest(c *gin.Context) {
// // 	id := c.Param("id")
	
// // 	query := "UPDATE users SET role = 'admin', is_pending_admin = false WHERE id = $1"
// // 	err := h.db.Exec(query, id)
// // 	if err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить роль"})
// // 		return
// // 	}

// // 	c.JSON(http.StatusOK, gin.H{"message": "Пользователь подтвержден как админ"})
// // }










// // package handlers

// // import (
// // 	"net/http"
// // 	"mobile-library/internal/models"
// // 	"mobile-library/internal/repositories"
// // 	"github.com/gin-gonic/gin"
// // )

// // type AdminHandler struct {
// // 	db *repositories.DB
// // }

// // func NewAdminHandler(db *repositories.DB) *AdminHandler {
// // 	return &AdminHandler{db: db}
// // }

// // // GetStats возвращает данные для дашборда через чистый SQL
// // func (h *AdminHandler) GetStats(c *gin.Context) {
// // 	var totalBooks, totalMembers, activeBorrows, overdueBooks int

// // 	conn := h.db.GetConn()

// // 	// 1. Считаем книги
// // 	conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&totalBooks)

// // 	// 2. Считаем пользователей
// // 	conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalMembers)

// // 	// 3. Считаем активные аренды
// // 	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'active'").Scan(&activeBorrows)

// // 	// 4. Считаем просроченные
// // 	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'overdue' OR (status = 'active' AND due_date < NOW())").Scan(&overdueBooks)

// // 	// Блок для графиков пока закомментируем, чтобы не было ошибки компиляции
// // 	/*
// // 	rows, _ := conn.Query(`
// // 	    SELECT TO_CHAR(borrow_date, 'Mon'), COUNT(*) 
// // 	    FROM borrows 
// // 	    GROUP BY TO_CHAR(borrow_date, 'Mon')
// // 	`)
// // 	defer rows.Close()
// // 	*/

// // 	c.JSON(http.StatusOK, gin.H{
// // 		"total_books":    totalBooks,
// // 		"total_members":  totalMembers,
// // 		"active_borrows": activeBorrows,
// // 		"overdue_books":  overdueBooks,
// // 		"revenue":        0,
// // 	})
// // }

// // // GetAllMembers возвращает список всех пользователей
// // func (h *AdminHandler) GetAllMembers(c *gin.Context) {
// // 	rows, err := h.db.GetConn().Query("SELECT id, name, email, phone, role, job_position, date_of_birth, created_at FROM users")
// // 	if err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
// // 		return
// // 	}
// // 	defer rows.Close()

// // 	var response []models.UserResponse
// // 	for rows.Next() {
// // 		var u models.User
// // 		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.JobPosition, &u.DateOfBirth, &u.CreatedAt); err != nil {
// // 			continue
// // 		}
// // 		response = append(response, u.ToResponse())
// // 	}

// // 	c.JSON(http.StatusOK, gin.H{"members": response})
// // }

// // // // AcceptAdminRequest подтверждает права админа
// // // func (h *AdminHandler) AcceptAdminRequest(c *gin.Context) {
// // // 	id := c.Param("id")
	
// // // 	query := "UPDATE users SET role = 'admin', is_pending_admin = false WHERE id = $1"
// // // 	err := h.db.Exec(query, id)
// // // 	if err != nil {
// // // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить роль"})
// // // 		return
// // // 	}

// // // 	c.JSON(http.StatusOK, gin.H{"message": "Пользователь подтвержден как админ"})
// // // }

// // // AcceptAdminRequest подтверждает права админа
// // func (h *AdminHandler) AcceptAdminRequest(c *gin.Context) {
// //     id := c.Param("id")
    
// //     conn := h.db.GetConn()
// //     query := "UPDATE users SET role = 'admin', is_pending_admin = false WHERE id = $1"
    
// //     // Exec возвращает (Result, error), игнорируем Result через заглушку "_"
// //     _, err := conn.Exec(query, id)
// //     if err != nil {
// //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить роль"})
// //         return
// //     }

// //     c.JSON(http.StatusOK, gin.H{"message": "Пользователь подтвержден как админ"})
// // // }
// // package handlers

// // import (
// //     "net/http"
// //     "mobile-library/internal/models"
// //     "mobile-library/internal/repositories"
// //     "github.com/gin-gonic/gin"
// // )

// // type AdminHandler struct {
// //     db *repositories.DB
// // }

// // func NewAdminHandler(db *repositories.DB) *AdminHandler {
// //     return &AdminHandler{db: db}
// // }

// // // GetStats godoc
// // // @Summary      Получить статистику
// // // @Description  Возвращает количество книг, юзеров и аренды
// // // @Tags         admin
// // // @Security     ApiKeyAuth
// // // @Success      200 {object} map[string]interface{}
// // // @Failure      500 {object} map[string]string
// // // @Router       /admin/stats [get]
// // func (h *AdminHandler) GetStats(c *gin.Context) {
// //     var totalBooks, totalMembers, activeBorrows, overdueBooks int
// //     conn := h.db.GetConn()
// //     conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&totalBooks)
// //     conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalMembers)
// //     conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'active'").Scan(&activeBorrows)
// //     conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'overdue' OR (status = 'active' AND due_date < NOW())").Scan(&overdueBooks)

// //     c.JSON(http.StatusOK, gin.H{
// //         "total_books":    totalBooks,
// //         "total_members":  totalMembers,
// //         "active_borrows": activeBorrows,
// //         "overdue_books":  overdueBooks,
// //         "revenue":        0,
// //     })
// // }

// // // GetAllMembers godoc
// // // @Summary      Список всех пользователей
// // // @Description  Получить список всех зарегистрированных в системе пользователей
// // // @Tags         admin
// // // @Security     ApiKeyAuth
// // // @Success      200 {object} map[string]interface{}
// // // @Failure      500 {object} map[string]string
// // // @Router       /admin/members [get]
// // func (h *AdminHandler) GetAllMembers(c *gin.Context) {
// //     rows, err := h.db.GetConn().Query("SELECT id, name, email, phone, role, job_position, date_of_birth, created_at FROM users")
// //     if err != nil {
// //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
// //         return
// //     }
// //     defer rows.Close()

// //     var response []models.UserResponse
// //     for rows.Next() {
// //         var u models.User
// //         if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.JobPosition, &u.DateOfBirth, &u.CreatedAt); err != nil {
// //             continue
// //         }
// //         response = append(response, u.ToResponse())
// //     }
// //     c.JSON(http.StatusOK, gin.H{"members": response})
// // }

// // // AcceptAdminRequest godoc
// // // @Summary      Подтвердить права админа
// // // @Description  Сделать пользователя администратором
// // // @Tags         admin
// // // @Security     ApiKeyAuth
// // // @Param        id   path      int  true  "ID пользователя"
// // // @Success      200  {object}  map[string]interface{}
// // // @Failure      500  {object}  map[string]string
// // // @Router       /admin/members/{id}/accept [post]
// // func (h *AdminHandler) AcceptAdminRequest(c *gin.Context) {
// //     id := c.Param("id")
// //     _, err := h.db.GetConn().Exec("UPDATE users SET role = 'admin', is_pending_admin = false WHERE id = $1", id)
// //     if err != nil {
// //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить роль"})
// //         return
// //     }
// //     c.JSON(http.StatusOK, gin.H{"message": "Пользователь подтвержден как админ"})
// // }

// // // DeleteMember godoc
// // // @Summary      Удалить пользователя
// // // @Description  Удаляет пользователя по ID
// // // @Tags         admin
// // // @Security     ApiKeyAuth
// // // @Param        id   path      int  true  "ID пользователя"
// // // @Success      200  {object}  map[string]interface{}
// // // @Failure      500  {object}  map[string]string
// // // @Router       /admin/members/{id} [delete]
// // func (h *AdminHandler) DeleteMember(c *gin.Context) {
// //     id := c.Param("id")
// //     _, err := h.db.GetConn().Exec("DELETE FROM users WHERE id = $1", id)
// //     if err != nil {
// //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить пользователя"})
// //         return
// //     }
// //     c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
// // }


// // // 1. Stat (Общая статистика)
// // // GET /admin/dashboard/stats
// // func (h *AdminHandler) GetStats(c *gin.Context) {
// //     // Твой текущий метод GetStats подходит
// //     // (Возвращает total_books, total_members, active_borrows, overdue_books)
// // }

// // // 2. Stat in month (Для графиков)
// // // GET /admin/dashboard/chart
// // func (h *AdminHandler) GetChartData(c *gin.Context) {
// //     // SQL: GROUP BY по месяцу и дате
// //     // Возвращает массив объектов {overdue, borrowed, date, month}
// //     c.JSON(http.StatusOK, stats)
// // }

// // // 3. Overdue Received Members
// // // GET /admin/dashboard/overdue
// // func (h *AdminHandler) GetOverdueMembers(c *gin.Context) {
// //     // SQL: SELECT users.name, books.title, borrows.due_date ... 
// //     // WHERE borrows.status = 'active' AND borrows.due_date < NOW()
// //     c.JSON(http.StatusOK, overdueList)
// // }

// // // GetChartData godoc
// // // @Summary      Статистика по месяцам для графика
// // // @Description  Возвращает количество просроченных и взятых книг по месяцам
// // // @Tags         admin
// // // @Security     ApiKeyAuth
// // // @Success      200 {array} object
// // // @Router       /admin/chart [get]
// // func (h *AdminHandler) GetChartData(c *gin.Context) {
// //     query := `
// //         SELECT 
// //             TO_CHAR(borrow_date, 'MM') as month,
// //             TO_CHAR(borrow_date, 'DD-MM-YYYY') as date,
// //             COUNT(*) FILTER (WHERE status = 'active') as borrowed,
// //             COUNT(*) FILTER (WHERE status = 'overdue' OR (status = 'active' AND due_date < NOW())) as overdue
// //         FROM borrows
// //         GROUP BY 1, 2
// //         ORDER BY 2 ASC
// //     `
// //     rows, err := h.db.GetConn().Query(query)
// //     if err != nil {
// //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных графика"})
// //         return
// //     }
// //     defer rows.Close()

// //     var stats []gin.H
// //     for rows.Next() {
// //         var month, date string
// //         var borrowed, overdue int
// //         rows.Scan(&month, &date, &borrowed, &overdue)
// //         stats = append(stats, gin.H{
// //             "month":    month,
// //             "date":     date,
// //             "borrowed": borrowed,
// //             "overdue":  overdue,
// //         })
// //     }
// //     c.JSON(http.StatusOK, stats)
// // }

// // // GetOverdueMembers godoc
// // // @Summary      Список должников
// // // @Description  Получает список пользователей, у которых просрочен возврат книг
// // // @Tags         admin
// // // @Security     ApiKeyAuth
// // // @Success      200 {array} object
// // // @Router       /admin/overdue [get]
// // func (h *AdminHandler) GetOverdueMembers(c *gin.Context) {
// //     query := `
// //         SELECT 
// //             u.id, u.name, u.phone, b.title, br.borrow_date, br.due_date,
// //             (NOW()::date - br.due_date::date) as days_overdue
// //         FROM borrows br
// //         JOIN users u ON br.user_id = u.id
// //         JOIN books b ON br.book_id = b.id
// //         WHERE br.status = 'active' AND br.due_date < NOW()
// //     `
// //     rows, err := h.db.GetConn().Query(query)
// //     if err != nil {
// //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка должников"})
// //         return
// //     }
// //     defer rows.Close()

// //     var overdueList []gin.H
// //     for rows.Next() {
// //         var id, days int
// //         var name, phone, title, bDate, dDate string
// //         rows.Scan(&id, &name, &phone, &title, &bDate, &dDate, &days)
// //         overdueList = append(overdueList, gin.H{
// //             "id":          id,
// //             "name":        name,
// //             "phone":       phone,
// //             "book_title":  title,
// //             "borrow_date": bDate,
// //             "dueDate":     dDate,
// //             "daysOverdue": days,
// //         })
// //     }
// //     c.JSON(http.StatusOK, overdueList)
// // }

// package handlers

// import (
// 	"net/http"
// 	"mobile-library/internal/models"
// 	"mobile-library/internal/repositories"
// 	"github.com/gin-gonic/gin"
// )

// type AdminHandler struct {
// 	db *repositories.DB
// }

// func NewAdminHandler(db *repositories.DB) *AdminHandler {
// 	return &AdminHandler{db: db}
// }

// // GetStats godoc
// // @Summary      Получить статистику
// // @Description  Возвращает количество книг, юзеров и аренды
// // @Tags         admin
// // @Security     ApiKeyAuth
// // @Success      200 {object} map[string]interface{}
// // @Failure      500 {object} map[string]string
// // @Router       /admin/stats [get]
// func (h *AdminHandler) GetStats(c *gin.Context) {
// 	var totalBooks, totalMembers, activeBorrows, overdueBooks int
// 	conn := h.db.GetConn()
	
// 	conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&totalBooks)
// 	conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalMembers)
// 	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'active'").Scan(&activeBorrows)
// 	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'overdue' OR (status = 'active' AND due_date < NOW())").Scan(&overdueBooks)

// 	c.JSON(http.StatusOK, gin.H{
// 		"total_books":    totalBooks,
// 		"total_members":  totalMembers,
// 		"active_borrows": activeBorrows,
// 		"overdue_books":  overdueBooks,
// 		"revenue":        0,
// 	})
// }

// // GetAllMembers godoc
// // @Summary      Список всех пользователей
// // @Description  Получить список всех зарегистрированных в системе пользователей
// // @Tags         admin
// // @Security     ApiKeyAuth
// // @Success      200 {object} map[string]interface{}
// // @Router       /admin/members [get]
// func (h *AdminHandler) GetAllMembers(c *gin.Context) {
// 	rows, err := h.db.GetConn().Query("SELECT id, name, email, phone, role, job_position, date_of_birth, created_at FROM users")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
// 		return
// 	}
// 	defer rows.Close()

// 	var response []models.UserResponse
// 	for rows.Next() {
// 		var u models.User
// 		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.JobPosition, &u.DateOfBirth, &u.CreatedAt); err != nil {
// 			continue
// 		}
// 		response = append(response, u.ToResponse())
// 	}
// 	c.JSON(http.StatusOK, gin.H{"members": response})
// }

// // AcceptAdminRequest godoc
// // @Summary      Подтвердить права админа
// // @Router       /admin/members/{id}/accept [post]
// func (h *AdminHandler) AcceptAdminRequest(c *gin.Context) {
// 	id := c.Param("id")
// 	_, err := h.db.GetConn().Exec("UPDATE users SET role = 'admin', is_pending_admin = false WHERE id = $1", id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить роль"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Пользователь подтвержден как админ"})
// }

// // DeleteMember godoc
// // @Summary      Удалить пользователя
// // @Router       /admin/members/{id} [delete]
// func (h *AdminHandler) DeleteMember(c *gin.Context) {
// 	id := c.Param("id")
// 	_, err := h.db.GetConn().Exec("DELETE FROM users WHERE id = $1", id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить пользователя"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
// }

// // GetChartData godoc
// // @Summary      Статистика для графика
// // @Router       /admin/chart [get]
// func (h *AdminHandler) GetChartData(c *gin.Context) {
// 	query := `
// 		SELECT 
// 			TO_CHAR(borrow_date, 'MM') as month,
// 			TO_CHAR(borrow_date, 'DD-MM-YYYY') as date,
// 			COUNT(*) FILTER (WHERE status = 'active') as borrowed,
// 			COUNT(*) FILTER (WHERE status = 'overdue' OR (status = 'active' AND due_date < NOW())) as overdue
// 		FROM borrows
// 		GROUP BY 1, 2
// 		ORDER BY 2 ASC
// 	`
// 	rows, err := h.db.GetConn().Query(query)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных"})
// 		return
// 	}
// 	defer rows.Close()

// 	var stats []gin.H
// 	for rows.Next() {
// 		var month, date string
// 		var borrowed, overdue int
// 		rows.Scan(&month, &date, &borrowed, &overdue)
// 		stats = append(stats, gin.H{
// 			"month": month, "date": date, "borrowed": borrowed, "overdue": overdue,
// 		})
// 	}
// 	c.JSON(http.StatusOK, stats)
// }

// // GetOverdueMembers godoc
// // @Summary      Список должников
// // @Router       /admin/overdue [get]
// func (h *AdminHandler) GetOverdueMembers(c *gin.Context) {
// 	query := `
// 		SELECT u.id, u.name, u.phone, b.title, br.borrow_date, br.due_date,
// 		(NOW()::date - br.due_date::date) as days_overdue
// 		FROM borrows br
// 		JOIN users u ON br.user_id = u.id
// 		JOIN books b ON br.book_id = b.id
// 		WHERE br.status = 'active' AND br.due_date < NOW()
// 	`
// 	rows, err := h.db.GetConn().Query(query)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
// 		return
// 	}
// 	defer rows.Close()

// 	var overdueList []gin.H
// 	for rows.Next() {
// 		var id, days int
// 		var name, phone, title, bDate, dDate string
// 		rows.Scan(&id, &name, &phone, &title, &bDate, &dDate, &days)
// 		overdueList = append(overdueList, gin.H{
// 			"id": id, "name": name, "phone": phone, "book_title": title,
// 			"borrow_date": bDate, "dueDate": dDate, "daysOverdue": days,
// 		})
// 	}
// 	c.JSON(http.StatusOK, overdueList)
// }



// // GetUserBookshelf godoc
// // @Summary      Книги пользователя
// // @Description  Возвращает список книг пользователя (активные и история)
// // @Tags         admin
// // @Security     ApiKeyAuth
// // @Param        user_id path string true "User ID"
// // @Success      200 {array} models.UserBookshelfItem
// // @Failure      500 {object} map[string]string
// // @Router       /admin/members/{user_id}/books [get]
// func (h *AdminHandler) GetUserBookshelf(c *gin.Context) {
// 	userID := c.Param("user_id")

// 	query := `
//         SELECT b.id, b.image_url, b.title, b.author, 
//                COALESCE(br.borrow_date::text, ''), 
//                COALESCE(br.due_date::text, ''), 
//                br.status
//         FROM borrows br
//         JOIN books b ON br.book_id = b.id
//         WHERE br.user_id = $1
//         ORDER BY br.borrow_date DESC
//     `

// 	rows, err := h.db.GetConn().Query(query, userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении книг пользователя"})
// 		return
// 	}
// 	defer rows.Close()

// 	var bookshelf []models.UserBookshelfItem
// 	for rows.Next() {
// 		var item models.UserBookshelfItem
// 		if err := rows.Scan(&item.ID, &item.ImageURL, &item.Title, &item.Author, 
//                             &item.BorrowDate, &item.DueDate, &item.Status); err != nil {
// 			continue
// 		}
// 		bookshelf = append(bookshelf, item)
// 	}

// 	c.JSON(http.StatusOK, bookshelf)
// }


// // AcceptReceiveRequest подтверждает заявку на получение книги
// func (h *AdminHandler) AcceptReceiveRequest(c *gin.Context) {
// 	requestID := c.Param("id")

// 	// Начинаем транзакцию
// 	tx, err := h.db.GetConn().Begin()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка транзакции"})
// 		return
// 	}

// 	// 1. Получаем данные заявки, чтобы перенести их в аренду
// 	var userID, bookID string
// 	err = tx.QueryRow("DELETE FROM receive_requests WHERE id = $1 RETURNING user_id, book_id", requestID).Scan(&userID, &bookID)
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Заявка не найдена"})
// 		return
// 	}

// 	// 2. Создаем запись об аренде
// 	_, err = tx.Exec("INSERT INTO borrows (user_id, book_id, borrow_date, due_date, status) VALUES ($1, $2, NOW(), NOW() + INTERVAL '14 days', 'active')", userID, bookID)
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании аренды"})
// 		return
// 	}

// 	// 3. Уменьшаем количество книг
// 	_, err = tx.Exec("UPDATE books SET available_copies = available_copies - 1 WHERE id = $1", bookID)
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Книга недоступна"})
// 		return
// 	}

// 	tx.Commit()
// 	c.JSON(http.StatusOK, gin.H{"message": "Заявка одобрена, книга выдана"})
// }

// // GetReceiveRequests godoc
// // @Summary      Список заявок на получение
// // @Router       /admin/receive-requests [get]
// func (h *AdminHandler) GetReceiveRequests(c *gin.Context) {
//     query := `
//         SELECT rr.id, u.name, b.title, rr.request_date
//         FROM receive_requests rr
//         JOIN users u ON rr.user_id = u.id
//         JOIN books b ON rr.book_id = b.id
//     `
//     rows, err := h.db.GetConn().Query(query)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
//         return
//     }
//     defer rows.Close()

//     var requests []gin.H
//     for rows.Next() {
//         var id, name, title, date string
//         rows.Scan(&id, &name, &title, &date)
//         requests = append(requests, gin.H{
//             "id": id, "receiver_name": name, "book_title": title, "request_date": date,
//         })
//     }
//     c.JSON(http.StatusOK, requests)
// }

// // DeleteReceiveRequest godoc
// // @Summary      Отклонить/удалить заявку
// // @Router       /admin/receive-requests/{id} [delete]
// func (h *AdminHandler) DeleteReceiveRequest(c *gin.Context) {
//     id := c.Param("id")
//     _, err := h.db.GetConn().Exec("DELETE FROM receive_requests WHERE id = $1", id)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить заявку"})
//         return
//     }
//     c.JSON(http.StatusOK, gin.H{"message": "Заявка удалена"})
// }


// // AcceptReturnRequest подтверждает возврат книги
// // @Summary      Подтвердить возврат книги
// // @Router       /admin/return-requests/{id}/accept [post]
// func (h *AdminHandler) AcceptReturnRequest(c *gin.Context) {
//     requestID := c.Param("id")

//     tx, err := h.db.GetConn().Begin()
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка транзакции"})
//         return
//     }

//     // 1. Получаем данные заявки на возврат и обновляем статус аренды
//     var bookID string
//     // Предполагаем, что таблица return_requests связывает аренду (borrow_id)
//     err = tx.QueryRow(`
//         DELETE FROM return_requests WHERE id = $1 
//         RETURNING book_id
//     `, requestID).Scan(&bookID)
    
//     if err != nil {
//         tx.Rollback()
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Заявка не найдена"})
//         return
//     }

//     // 2. Обновляем статус аренды на 'returned'
//     _, err = tx.Exec("UPDATE borrows SET status = 'returned' WHERE book_id = $1 AND status = 'active'", bookID)
    
//     // 3. Возвращаем книгу в фонд (увеличиваем количество)
//     _, err = tx.Exec("UPDATE books SET available_copies = available_copies + 1 WHERE id = $1", bookID)
//     if err != nil {
//         tx.Rollback()
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении книги"})
//         return
//     }

//     tx.Commit()
//     c.JSON(http.StatusOK, gin.H{"message": "Книга принята обратно, заявка закрыта"})
// }


// // GetReturnRequests godoc
// // @Summary      Список заявок на возврат
// // @Router       /admin/return-requests [get]
// func (h *AdminHandler) GetReturnRequests(c *gin.Context) {
//     query := `
//         SELECT rr.id, u.name, b.title, rr.request_date
//         FROM return_requests rr
//         JOIN users u ON rr.user_id = u.id
//         JOIN books b ON rr.book_id = b.id
//     `
//     rows, err := h.db.GetConn().Query(query)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
//         return
//     }
//     defer rows.Close()

//     var requests []gin.H
//     for rows.Next() {
//         var id, name, title, date string
//         rows.Scan(&id, &name, &title, &date)
//         requests = append(requests, gin.H{
//             "id": id, "returner_name": name, "book_title": title, "request_date": date,
//         })
//     }
//     c.JSON(http.StatusOK, requests)
// }

// // DeleteReturnRequest godoc
// // @Summary      Отклонить/удалить заявку на возврат
// // @Router       /admin/return-requests/{id} [delete]
// func (h *AdminHandler) DeleteReturnRequest(c *gin.Context) {
//     id := c.Param("id")
//     _, err := h.db.GetConn().Exec("DELETE FROM return_requests WHERE id = $1", id)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить заявку"})
//         return
//     }
//     c.JSON(http.StatusOK, gin.H{"message": "Заявка на возврат удалена"})
// }
// // AddNotification создает уведомление
// // func (h *AdminHandler) AddNotification(c *gin.Context) {
// //     var req struct {
// //         Member               string `json:"member"` // "user_id" или "all_users"
// //         NotificationType     string `json:"notification_type"`
// //         NotificationImageURL string `json:"notification_image_url"`
// //         Title                string `json:"title"`
// //         Description          string `json:"description"`
// //     }
// //     if err := c.ShouldBindJSON(&req); err != nil {
// //         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// //         return
// //     }

// //     // Логика обработки поля Member
// //     var userID *uint
// //     if req.Member != "all_users" {
// //         // Конвертируем строку-ID в uint (нужна будет утилита или Atoi)
// //         // ... (логика парсинга userID)
// //     }

// //     newNotif := models.Notification{
// //         UserID:               userID,
// //         NotificationType:     req.NotificationType,
// //         NotificationImageURL: req.NotificationImageURL,
// //         Title:                req.Title,
// //         Description:          req.Description,
// //         CreatedAt:            time.Now(),
// //     }

// //     h.db.Create(&newNotif)
// //     c.JSON(http.StatusCreated, newNotif)
// // }

// // // GetNotifications возвращает список
// // func (h *AdminHandler) GetNotifications(c *gin.Context) {
// //     var notifs []models.Notification
// //     h.db.Find(&notifs)
// //     c.JSON(http.StatusOK, notifs)
// // }


// func (h *AdminHandler) AddNotification(c *gin.Context) {
//     var req struct {
//         Member               string `json:"member"`
//         NotificationType     string `json:"notification_type"`
//         NotificationImageURL string `json:"notification_image_url"`
//         Title                string `json:"title"`
//         Description          string `json:"description"`
//     }
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     // Создаем модель (используем Message для совместимости)
//     newNotif := models.Notification{
//         Message:              req.Description, 
//         NotificationType:     req.NotificationType,
//         NotificationImageURL: req.NotificationImageURL,
//         Title:                req.Title,
//         Description:          req.Description,
//         CreatedAt:            time.Now(),
//     }

//     // ВЫЗЫВАЕМ РЕПОЗИТОРИЙ, а не h.db.Create
//     if err := h.notifRepo.Create(&newNotif); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в БД"})
//         return
//     }
    
//     c.JSON(http.StatusCreated, newNotif)
// }

// func (h *AdminHandler) GetNotifications(c *gin.Context) {
//     // Чистый SQL запрос
//     rows, err := h.db.GetConn().Query("SELECT id, message, is_read, created_at FROM notifications")
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
//         return
//     }
//     defer rows.Close()

//     var notifs []models.Notification
//     for rows.Next() {
//         var n models.Notification
//         // Сканируем поля. Убедись, что список сканирования совпадает с SELECT
//         if err := rows.Scan(&n.ID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
//             continue
//         }
//         notifs = append(notifs, n)
//     }
    
//     c.JSON(http.StatusOK, notifs)
// }











































































// package handlers

// import (
// 	"net/http"
// 	"mobile-library/internal/models"
// 	"mobile-library/internal/repositories"
// 	"github.com/gin-gonic/gin"
// 	"time"
// )

// type AdminHandler struct {
// 	db        *repositories.DB
// 	notifRepo *repositories.NotificationRepository
// }

// func NewAdminHandler(db *repositories.DB, nr *repositories.NotificationRepository) *AdminHandler {
// 	return &AdminHandler{db: db, notifRepo: nr}
// }

// // --- СТАТИСТИКА И ДАШБОРД ---

// func (h *AdminHandler) GetStats(c *gin.Context) {
// 	var totalBooks, totalMembers, activeBorrows, overdueBooks int
// 	conn := h.db.GetConn()
// 	conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&totalBooks)
// 	conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalMembers)
// 	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'active'").Scan(&activeBorrows)
// 	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'overdue' OR (status = 'active' AND due_date < NOW())").Scan(&overdueBooks)

// 	c.JSON(http.StatusOK, gin.H{
// 		"total_books": totalBooks, "total_members": totalMembers,
// 		"active_borrows": activeBorrows, "overdue_books": overdueBooks, "revenue": 0,
// 	})
// }

// // --- ПОЛЬЗОВАТЕЛИ ---

// func (h *AdminHandler) GetAllMembers(c *gin.Context) {
// 	rows, err := h.db.GetConn().Query("SELECT id, name, email, phone, role, job_position, date_of_birth, created_at FROM users")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
// 		return
// 	}
// 	defer rows.Close()

// 	var response []models.UserResponse
// 	for rows.Next() {
// 		var u models.User
// 		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.JobPosition, &u.DateOfBirth, &u.CreatedAt); err == nil {
// 			response = append(response, u.ToResponse())
// 		}
// 	}
// 	c.JSON(http.StatusOK, gin.H{"members": response})
// }

// func (h *AdminHandler) AcceptAdminRequest(c *gin.Context) {
// 	id := c.Param("id")
// 	_, err := h.db.GetConn().Exec("UPDATE users SET role = 'admin', is_pending_admin = false WHERE id = $1", id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Статус обновлен"})
// }

// func (h *AdminHandler) DeleteMember(c *gin.Context) {
// 	id := c.Param("id")
// 	h.db.GetConn().Exec("DELETE FROM users WHERE id = $1", id)
// 	c.JSON(http.StatusOK, gin.H{"message": "Удалено"})
// }

// // --- ЗАЯВКИ (RECEIVE & RETURN) ---

// func (h *AdminHandler) AcceptReceiveRequest(c *gin.Context) {
// 	requestID := c.Param("id")
// 	tx, _ := h.db.GetConn().Begin()
// 	var userID, bookID string
// 	err := tx.QueryRow("DELETE FROM receive_requests WHERE id = $1 RETURNING user_id, book_id", requestID).Scan(&userID, &bookID)
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Заявка не найдена"})
// 		return
// 	}
// 	tx.Exec("INSERT INTO borrows (user_id, book_id, borrow_date, due_date, status) VALUES ($1, $2, NOW(), NOW() + INTERVAL '14 days', 'active')", userID, bookID)
// 	tx.Exec("UPDATE books SET available_copies = available_copies - 1 WHERE id = $1", bookID)
// 	tx.Commit()
// 	c.JSON(http.StatusOK, gin.H{"message": "Книга выдана"})
// }

// // --- УВЕДОМЛЕНИЯ ---

// // func (h *AdminHandler) AddNotification(c *gin.Context) {
// // 	var req struct {
// // 		Member               string `json:"member"`
// // 		NotificationType     string `json:"notification_type"`
// // 		NotificationImageURL string `json:"notification_image_url"`
// // 		Title                string `json:"title"`
// // 		Description          string `json:"description"`
// // 	}
// // 	if err := c.ShouldBindJSON(&req); err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
// // 		return
// // 	}

// // 	newNotif := models.Notification{
// //         UserID:    req.UserID, // Передавай как есть, если в req это *uint
// //         Message:   req.Message,
// //         IsRead:    false,
// //         CreatedAt: time.Now(),
// //     }

// //     // ИСПОЛЬЗУЙ РЕПОЗИТОРИЙ, а не h.db.Create!
// //     if err := h.notifRepo.Create(&newNotif); err != nil {
// //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
// //         return
// //     }
// //     c.JSON(http.StatusCreated, newNotif)
// // }

// func (h *AdminHandler) AddNotification(c *gin.Context) {
//     var req struct {
//         UserID  *uint  `json:"user_id"`
//         Message string `json:"message"`
//     }
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     newNotif := models.Notification{
//         UserID:    req.UserID,
//         Message:   req.Message,
//         IsRead:    false,
//         CreatedAt: time.Now(),
//     }

//     // ИСПОЛЬЗУЙ РЕПОЗИТОРИЙ, а не h.db.Create
//     if err := h.notifRepo.Create(&newNotif); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
//         return
//     }
//     c.JSON(http.StatusCreated, newNotif)
// }


// func (h *AdminHandler) GetUserBookshelf(c *gin.Context) {
// 	userID := c.Param("user_id")

// 	// SQL запрос для получения списка книг пользователя
// 	query := `
// 		SELECT b.id, b.title, b.author, br.borrow_date, br.status
// 		FROM borrows br
// 		JOIN books b ON br.book_id = b.id
// 		WHERE br.user_id = $1
// 		ORDER BY br.borrow_date DESC`

// 	rows, err := h.db.GetConn().Query(query, userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении книг из БД"})
// 		return
// 	}
// 	defer rows.Close()

// 	// Структура для ответа
// 	type BookInfo struct {
// 		ID         int       `json:"id"`
// 		Title      string    `json:"title"`
// 		Author     string    `json:"author"`
// 		BorrowDate time.Time `json:"borrow_date"`
// 		Status     string    `json:"status"`
// 	}

// 	var books []BookInfo
// 	for rows.Next() {
// 		var b BookInfo
// 		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.BorrowDate, &b.Status); err == nil {
// 			books = append(books, b)
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"user_id": userID,
// 		"books":   books,
// 	})
// }

// func (h *AdminHandler) GetAdminManagement(c *gin.Context) {
//     // 1. Получаем список действующих админов
//     adminRows, err := h.db.GetConn().Query("SELECT id, name, email, role, job_position FROM users WHERE role = 'admin'")
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка админов"})
//         return
//     }
//     defer adminRows.Close()

//     type AdminInfo struct {
//         ID          int    `json:"id"`
//         Name        string `json:"name"`
//         Email       string `json:"email"`
//         Role        string `json:"role"`
//         JobPosition string `json:"job_position"`
//     }

//     var admins []AdminInfo
//     for adminRows.Next() {
//         var a AdminInfo
//         if err := adminRows.Scan(&a.ID, &a.Name, &a.Email, &a.Role, &a.JobPosition); err == nil {
//             admins = append(admins, a)
//         }
//     }

//     // 2. Получаем список заявок на админство (те, кто ожидает подтверждения: is_pending_admin = true)
//     pendingRows, err := h.db.GetConn().Query("SELECT id, name, email, job_position FROM users WHERE is_pending_admin = true")
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении заявок"})
//         return
//     }
//     defer pendingRows.Close()

//     type PendingAdmin struct {
//         ID          int    `json:"id"`
//         Name        string `json:"name"`
//         Email       string `json:"email"`
//         JobPosition string `json:"job_position"`
//     }

//     var pendingUsers []PendingAdmin
//     for pendingRows.Next() {
//         var p PendingAdmin
//         if err := pendingRows.Scan(&p.ID, &p.Name, &p.Email, &p.JobPosition); err == nil {
//             pendingUsers = append(pendingUsers, p)
//         }
//     }

//     // Отдаем фронтенду сгруппированный чистый JSON
//     c.JSON(http.StatusOK, gin.H{
//         "current_admins": admins,
//         "pending_admins": pendingUsers,
//     })
// }

// func (h *AdminHandler) GetChartData(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"data": []int{}}) }
// func (h *AdminHandler) GetOverdueMembers(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"members": []string{}}) }
// func (h *AdminHandler) GetReceiveRequests(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"requests": []string{}}) }
// func (h *AdminHandler) DeleteReceiveRequest(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "deleted"}) }
// func (h *AdminHandler) GetReturnRequests(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"requests": []string{}}) }
// func (h *AdminHandler) AcceptReturnRequest(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "accepted"}) }
// func (h *AdminHandler) DeleteReturnRequest(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "deleted"}) }
// func (h *AdminHandler) GetNotifications(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"notifications": []string{}}) }













































package handlers

import (
	"github.com/gin-gonic/gin"
	"mobile-library/internal/models"
	"mobile-library/internal/repositories"
	"net/http"
	"time"
)

type AdminHandler struct {
	db        *repositories.DB
	notifRepo *repositories.NotificationRepository
}

func NewAdminHandler(db *repositories.DB, nr *repositories.NotificationRepository) *AdminHandler {
	return &AdminHandler{db: db, notifRepo: nr}
}

// --- СТАТИСТИКА И ДАШБОРД ---

func (h *AdminHandler) GetStats(c *gin.Context) {
	var totalBooks, totalMembers, activeBorrows, overdueBooks int
	conn := h.db.GetConn()
	conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&totalBooks)
	conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalMembers)
	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'active'").Scan(&activeBorrows)
	conn.QueryRow("SELECT COUNT(*) FROM borrows WHERE status = 'overdue' OR (status = 'active' AND due_date < NOW())").Scan(&overdueBooks)

	c.JSON(http.StatusOK, gin.H{
		"total_books": totalBooks, "total_members": totalMembers,
		"active_borrows": activeBorrows, "overdue_books": overdueBooks, "revenue": 0,
	})
}

// --- ПОЛЬЗОВАТЕЛИ ---

func (h *AdminHandler) GetAllMembers(c *gin.Context) {
	rows, err := h.db.GetConn().Query("SELECT id, name, email, phone, role, job_position, date_of_birth, created_at FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}
	defer rows.Close()

	var response []models.UserResponse
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.JobPosition, &u.DateOfBirth, &u.CreatedAt); err == nil {
			response = append(response, u.ToResponse())
		}
	}
	c.JSON(http.StatusOK, gin.H{"members": response})
}

func (h *AdminHandler) AcceptAdminRequest(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.GetConn().Exec("UPDATE users SET role = 'admin', is_pending_admin = false WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Статус обновлен"})
}

func (h *AdminHandler) DeleteMember(c *gin.Context) {
	id := c.Param("id")
	h.db.GetConn().Exec("DELETE FROM users WHERE id = $1", id)
	c.JSON(http.StatusOK, gin.H{"message": "Удалено"})
}

// --- ЗАЯВКИ (RECEIVE & RETURN) ---

func (h *AdminHandler) AcceptReceiveRequest(c *gin.Context) {
	requestID := c.Param("id")
	tx, _ := h.db.GetConn().Begin()
	var userID, bookID string
	err := tx.QueryRow("DELETE FROM receive_requests WHERE id = $1 RETURNING user_id, book_id", requestID).Scan(&userID, &bookID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Заявка не найдена"})
		return
	}
	tx.Exec("INSERT INTO borrows (user_id, book_id, borrow_date, due_date, status) VALUES ($1, $2, NOW(), NOW() + INTERVAL '14 days', 'active')", userID, bookID)
	tx.Exec("UPDATE books SET available_copies = available_copies - 1 WHERE id = $1", bookID)
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Книга выдана"})
}

// --- УВЕДОМЛЕНИЯ ---

func (h *AdminHandler) AddNotification(c *gin.Context) {
	var req struct {
		UserID  *uint  `json:"user_id"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newNotif := models.Notification{
		UserID:    req.UserID,
		Message:   req.Message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := h.notifRepo.Create(&newNotif); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
		return
	}
	c.JSON(http.StatusCreated, newNotif)
}

func (h *AdminHandler) GetUserBookshelf(c *gin.Context) {
	userID := c.Param("user_id")

	query := `
        SELECT b.id, b.title, b.author, br.borrow_date, br.status
        FROM borrows br
        JOIN books b ON br.book_id = b.id
        WHERE br.user_id = $1
        ORDER BY br.borrow_date DESC`

	rows, err := h.db.GetConn().Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении книг из БД"})
		return
	}
	defer rows.Close()

	type BookInfo struct {
		ID         int       `json:"id"`
		Title      string    `json:"title"`
		Author     string    `json:"author"`
		BorrowDate time.Time `json:"borrow_date"`
		Status     string    `json:"status"`
	}

	var books []BookInfo
	for rows.Next() {
		var b BookInfo
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.BorrowDate, &b.Status); err == nil {
			books = append(books, b)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"books":   books,
	})
}

func (h *AdminHandler) GetAdminManagement(c *gin.Context) {
	adminRows, err := h.db.GetConn().Query("SELECT id, name, email, role, job_position FROM users WHERE role = 'admin'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка админов"})
		return
	}
	defer adminRows.Close()

	type AdminInfo struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		Role        string `json:"role"`
		JobPosition string `json:"job_position"`
	}

	var admins []AdminInfo
	for adminRows.Next() {
		var a AdminInfo
		if err := adminRows.Scan(&a.ID, &a.Name, &a.Email, &a.Role, &a.JobPosition); err == nil {
			admins = append(admins, a)
		}
	}

	pendingRows, err := h.db.GetConn().Query("SELECT id, name, email, job_position FROM users WHERE is_pending_admin = true")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении заявок"})
		return
	}
	defer pendingRows.Close()

	type PendingAdmin struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		JobPosition string `json:"job_position"`
	}

	var pendingUsers []PendingAdmin
	for pendingRows.Next() {
		var p PendingAdmin
		if err := pendingRows.Scan(&p.ID, &p.Name, &p.Email, &p.JobPosition); err == nil {
			pendingUsers = append(pendingUsers, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"current_admins": admins,
		"pending_admins": pendingUsers,
	})
}

// Заглушки для роутов
func (h *AdminHandler) GetChartData(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"data": []int{}}) }
func (h *AdminHandler) GetOverdueMembers(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"members": []string{}}) }
func (h *AdminHandler) GetReceiveRequests(c *gin.Context)   { c.JSON(http.StatusOK, gin.H{"requests": []string{}}) }
func (h *AdminHandler) DeleteReceiveRequest(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "deleted"}) }
func (h *AdminHandler) GetReturnRequests(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"requests": []string{}}) }
func (h *AdminHandler) DeleteReturnRequest(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"message": "deleted"}) }
func (h *AdminHandler) GetNotifications(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"notifications": []string{}}) }
func (h *AdminHandler) AcceptReturnRequest(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"message": "accepted"}) }

// // Новые методы для управления фильтрами
// func (h *AdminHandler) GetFilters(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"filters": []string{}}) }
// func (h *AdminHandler) AddFilter(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"message": "added"}) }
// func (h *AdminHandler) UpdateFilter(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"message": "updated"}) }
// func (h *AdminHandler) DeleteFilter(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"message": "deleted"}) }







// Структура для парсинга входящего JSON-запроса от фронтенда
type FilterRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// GetFilters — Получение всех категорий/фильтров из базы данных
func (h *AdminHandler) GetFilters(c *gin.Context) {
	query := `SELECT id, name FROM categories ORDER BY name ASC`
	rows, err := h.db.GetConn().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении фильтров из базы данных"})
		return
	}
	defer rows.Close()

	type FilterResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var list []FilterResponse
	for rows.Next() {
		var f FilterResponse
		if err := rows.Scan(&f.ID, &f.Name); err == nil {
			list = append(list, f)
		}
	}

	if list == nil {
		list = []FilterResponse{} // Возвращаем пустой массив [] вместо null
	}

	c.JSON(http.StatusOK, gin.H{"filters": list})
}

// AddFilter — Реальное сохранение нового фильтра/категории в БД
func (h *AdminHandler) AddFilter(c *gin.Context) {
	var req FilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя фильтра обязательно (минимум 2 символа)"})
		return
	}

	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	var newID int
	err := h.db.GetConn().QueryRow(query, req.Name).Scan(&newID)
	if err != nil {
		// Проверяем на дубликат (если уникальное имя уже существует)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Такой фильтр или категория уже существует"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Фильтр успешно добавлен",
		"filter": gin.H{
			"id":   newID,
			"name": req.Name,
		},
	})
}

// UpdateFilter — Редактирование существующего фильтра по ID
func (h *AdminHandler) UpdateFilter(c *gin.Context) {
	id := c.Param("id")
	var req FilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Новое имя фильтра некорректно"})
		return
	}

	query := `UPDATE categories SET name = $1 WHERE id = $2`
	res, err := h.db.GetConn().Exec(query, req.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении фильтра"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильтр с таким ID не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Фильтр успешно обновлен"})
}

// DeleteFilter — Удаление фильтра из БД
func (h *AdminHandler) DeleteFilter(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM categories WHERE id = $1`
	res, err := h.db.GetConn().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нельзя удалить фильтр, так как он используется в книгах"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильтр с таким ID не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Фильтр успешно удален"})
}