// package handlers

// import (
// 	"net/http"
// 	"strconv"

// 	"mobile-library/internal/models"
// 	"mobile-library/internal/services"

// 	"github.com/gin-gonic/gin"
// )

// type BookHandler struct {
// 	service *services.BookService
// }

// func NewBookHandler(s *services.BookService) *BookHandler {
// 	return &BookHandler{service: s}
// }

// // GetAll godoc
// // @Summary      Get all books
// // @Description  Get a paginated list of all books
// // @Tags         books
// // @Produce      json
// // @Param        page       query     int  false  "Page number"
// // @Param        page_size  query     int  false  "Items per page"
// // @Success      200        {object}  map[string]interface{}
// // @Failure      500        {object}  map[string]string
// // @Router       /books [get]
// func (h *BookHandler) GetAll(c *gin.Context) {
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

// 	books, total, err := h.service.GetAll(page, pageSize)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":  books,
// 		"total": total,
// 		"page":  page,
// 	})
// }

// // GetByID godoc
// // @Summary      Get book by ID
// // @Description  Get detailed information about a specific book
// // @Tags         books
// // @Produce      json
// // @Param        id   path      int  true  "Book ID"
// // @Success      200  {object}  map[string]interface{}
// // @Failure      400  {object}  map[string]string
// // @Failure      404  {object}  map[string]string
// // @Router       /books/{id} [get]
// func (h *BookHandler) GetByID(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
// 		return
// 	}
// 	book, err := h.service.GetByID(uint(id))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"book": book})
// }

// // Search godoc
// // @Summary      Search books
// // @Description  Search books by title, author or category
// // @Tags         books
// // @Produce      json
// // @Param        title     query     string  false  "Book title"
// // @Param        author    query     string  false  "Book author"
// // @Param        category  query     string  false  "Book category"
// // @Param        page      query     int     false  "Page number"
// // @Param        page_size query     int     false  "Items per page"
// // @Success      200       {object}  map[string]interface{}
// // @Failure      400       {object}  map[string]string
// // @Failure      500       {object}  map[string]string
// // @Router       /books/search [get]
// func (h *BookHandler) Search(c *gin.Context) {
// 	var params models.BookSearchParams
// 	if err := c.ShouldBindQuery(&params); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
// 	params.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))

// 	books, total, err := h.service.Search(params)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"data":  books,
// 		"total": total,
// 	})
// }

// // Create godoc
// // @Summary      Create a new book
// // @Description  Add a new book to the library
// // @Tags         books
// // @Accept       json
// // @Produce      json
// // @Param        book  body      models.Book  true  "Book object"
// // @Success      201   {object}  models.Book
// // @Failure      400   {object}  map[string]string
// // @Failure      401   {object}  map[string]string
// // @Failure      500   {object}  map[string]string
// // @Security     ApiKeyAuth
// // @Router       /books [post]
// // func (h *BookHandler) Create(c *gin.Context) {
// // 	var book models.Book
// // 	if err := c.ShouldBindJSON(&book); err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// // 		return
// // 	}

// // 	if err := h.service.Create(&book); err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// // 		return
// // 	}

// // 	c.JSON(http.StatusCreated, book)
// // }

// func (h *BookHandler) Create(c *gin.Context) {
//     var req models.CreateBookRequest // Используем DTO
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     // Преобразуем DTO в модель для сервиса
//     book := models.Book{
//     Title:           req.Title,
//     Author:          req.Author,
//     Description:     req.Description,
//     Category:        req.Category,
//     Year:            req.Year,
//     AvailableCopies: req.AvailableCopies,
//     ImageURL:        req.ImageURL,
//     Language:        req.Language,
//     PageCount:       req.PageCount,
//     Status:          "available",
//     }

//     if err := h.service.Create(&book); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusCreated, book)
// }




// // Update обновляет данные о книге (PUT /books/:id)
// func (h *BookHandler) Update(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
// 		return
// 	}

// 	var input models.Book
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
// 		return
// 	}

// 	// Вызываем сервис, а не h.db
// 	if err := h.service.Update(uint(id), &input); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Книга обновлена"})
// }

// // Delete удаляет книгу (DELETE /books/:id)
// func (h *BookHandler) Delete(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
// 		return
// 	}

// 	// Вызываем сервис, а не h.db
// 	if err := h.service.Delete(uint(id)); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить книгу"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Книга удалена"})
// }


// // GetFilters godoc
// // @Summary      Список категорий
// // @Description  Возвращает список всех уникальных категорий книг
// // @Tags         books
// // @Security     ApiKeyAuth
// // @Success      200 {array} map[string]string
// // @Router       /admin/filters [get]
// func (h *BookHandler) GetFilters(c *gin.Context) {
//     // ВАЖНО: Предполагаем, что у тебя в BookService есть метод GetCategories
//     // Если его нет, создай его (см. ниже)
//     categories, err := h.service.GetCategories() 
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить фильтры"})
//         return
//     }

//     // Преобразуем список строк в нужный формат: [{"id": "1", "filterName": "Finance"}, ...]
//     var response []gin.H
//     for i, cat := range categories {
//         response = append(response, gin.H{
//             "id":         strconv.Itoa(i + 1),
//             "filterName": cat,
//         })
//     }

//     c.JSON(http.StatusOK, response)
// }

// func (h *BookHandler) AddFilter(c *gin.Context) {
//     var req struct { FilterName string `json:"filterName"` }
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
//         return
//     }
//     if err := h.service.CreateCategory(req.FilterName); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать фильтр"})
//         return
//     }
//     c.JSON(http.StatusCreated, gin.H{"message": "Фильтр создан"})
// }

// func (h *BookHandler) UpdateFilter(c *gin.Context) {
//     id := c.Param("id")
//     var req struct { FilterName string `json:"filterName"` }
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
//         return
//     }
//     if err := h.service.UpdateCategory(id, req.FilterName); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
//         return
//     }
//     c.JSON(http.StatusOK, gin.H{"message": "Фильтр обновлен"})
// }

// func (h *BookHandler) DeleteFilter(c *gin.Context) {
//     id := c.Param("id")
//     if err := h.service.DeleteCategory(id); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении"})
//         return
//     }
//     c.JSON(http.StatusOK, gin.H{"message": "Фильтр удален"})
// }



package handlers

import (
	"net/http"
	"strconv"

	"mobile-library/internal/models"
	"mobile-library/internal/services"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(s *services.BookService) *BookHandler {
	return &BookHandler{service: s}
}

// Вспомогательная функция для безопасного получения userID из JWT контекста
func getUserID(c *gin.Context) uint {
	// Проверяем как "userID", так и "user_id" (в зависимости от того, как у тебя написано в middleware)
	if id, exists := c.Get("userID"); exists {
		if uIntId, ok := id.(uint); ok {
			return uIntId
		}
		if floatId, ok := id.(float64); ok { // JWT библиотеки иногда парсят числа как float64
			return uint(floatId)
		}
	}
	if id, exists := c.Get("user_id"); exists {
		if uIntId, ok := id.(uint); ok {
			return uIntId
		}
		if floatId, ok := id.(float64); ok {
			return uint(floatId)
		}
	}
	return 0 // Если токена нет или мы в публичной зоне
}

// GetAll godoc
func (h *BookHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	currentUserID := getUserID(c)

	books, total, err := h.service.GetAll(page, pageSize, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"total": total,
		"page":  page,
	})
}

// GetByID godoc
func (h *BookHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	
	currentUserID := getUserID(c)
	book, err := h.service.GetByID(uint(id), currentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

// Search godoc
func (h *BookHandler) Search(c *gin.Context) {
	var params models.BookSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	params.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))

	currentUserID := getUserID(c)
	books, total, err := h.service.Search(params, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"total": total,
	})
}



// Локальное DTO прямо в хендлере, чтобы обойти любые проблемы с импортами и кэшем
type LocalCreateBookRequest struct {
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author" binding:"required"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	Year            int    `json:"year"`
	AvailableCopies int    `json:"available_copies"`
	ImageURL        string `json:"image_url"`
	BgImageURL      string `json:"bg_image_url"` // Точно на месте
	Language        string `json:"language"`
	PageCount       int    `json:"page_count"`
}

func (h *BookHandler) Create(c *gin.Context) {
	var req LocalCreateBookRequest // Используем локальную структуру
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем в модель для сервиса
	book := models.Book{
		Title:           req.Title,
		Author:          req.Author,
		Description:     req.Description,
		Category:        req.Category,
		Year:            req.Year,
		AvailableCopies: req.AvailableCopies,
		ImageURL:        req.ImageURL,
		BgImageURL:      req.BgImageURL, // Теперь компилятор 100% увидит это поле!
		Language:        req.Language,
		PageCount:       req.PageCount,
		Status:          "Available",
	}

	if err := h.service.Create(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}


// // Create — теперь учитывает bg_image_url, запрошенный фронтендом
// func (h *BookHandler) Create(c *gin.Context) {
// 	var req models.CreateBookRequest // Используем твой DTO
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Преобразуем DTO в модель для сервиса с учетом фона книги
// 	book := models.Book{
// 		Title:           req.Title,
// 		Author:          req.Author,
// 		Description:     req.Description,
// 		Category:        req.Category,
// 		Year:            req.Year,
// 		AvailableCopies: req.AvailableCopies,
// 		ImageURL:        req.ImageURL,
// 		BgImageURL:      req.BgImageURL, // Передаем фон
// 		Language:        req.Language,
// 		PageCount:       req.PageCount,
// 		Status:          "Available",
// 	}

// 	if err := h.service.Create(&book); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, book)
// }

// Update обновляет данные о книге (PUT /books/:id)
func (h *BookHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	if err := h.service.Update(uint(id), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга обновлена"})
}

// Delete удаляет книгу (DELETE /books/:id)
func (h *BookHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить книгу"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга удалена"})
}

// GetFilters godoc
func (h *BookHandler) GetFilters(c *gin.Context) {
	categories, err := h.service.GetCategories() 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить фильтры"})
		return
	}

	// Форматируем под мобилку фронта: [{"id": "1", "filterName": "Finance"}, ...]
	var response []gin.H
	for i, cat := range categories {
		response = append(response, gin.H{
			"id":         strconv.Itoa(i + 1),
			"filterName": cat,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *BookHandler) AddFilter(c *gin.Context) {
	var req struct { FilterName string `json:"filterName"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}
	if err := h.service.CreateCategory(req.FilterName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать фильтр"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Фильтр создан"})
}

func (h *BookHandler) UpdateFilter(c *gin.Context) {
	id := c.Param("id")
	var req struct { FilterName string `json:"filterName"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}
	if err := h.service.UpdateCategory(id, req.FilterName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Фильтр обновлен"})
}

func (h *BookHandler) DeleteFilter(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Фильтр удален"})
}