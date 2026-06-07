package handlers

import (
	"net/http"
	"strconv"
	"mobile-library/internal/models"
	"mobile-library/internal/repositories"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	repo *repositories.ReviewRepository
}

func NewReviewHandler(repo *repositories.ReviewRepository) *ReviewHandler {
	return &ReviewHandler{repo: repo}
}

// CreateReview — оставить отзыв к книге (POST /books/:id/reviews)
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вытаскиваем ID авторизованного юзера из контекста (наша вспомогательная логика)
	var currentUserID uint
	if id, exists := c.Get("userID"); exists {
		if uIntId, ok := id.(uint); ok { currentUserID = uIntId }
	} else if id, exists := c.Get("user_id"); exists {
		if uIntId, ok := id.(uint); ok { currentUserID = uIntId }
	}

	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	review := models.Review{
		BookID:         uint(bookID),
		UserID:         currentUserID,
		Rating:         req.Rating,
		Review:         req.Review,
		ReviewCategory: req.ReviewCategory,
	}

	if err := h.repo.Create(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add review: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetBookReviews — получить отзывы книги (GET /books/:id/reviews)
func (h *ReviewHandler) GetBookReviews(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	reviews, err := h.repo.GetByBookID(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	if reviews == nil {
		reviews = []models.Review{} // Возвращаем пустой массив вместо null для фронта
	}

	c.JSON(http.StatusOK, reviews)
}