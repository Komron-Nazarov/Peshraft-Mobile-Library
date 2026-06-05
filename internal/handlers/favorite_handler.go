package handlers

import (
	"net/http"
	"strconv"

	"mobile-library/internal/repositories"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	repo *repositories.FavoriteRepository
}

func NewFavoriteHandler(r *repositories.FavoriteRepository) *FavoriteHandler {
	return &FavoriteHandler{repo: r}
}

// Add godoc
// @Summary      Add book to favorites
// @Description  Add a specific book to the authenticated user's favorites list
// @Tags         favorites
// @Security     ApiKeyAuth
// @Param        book_id  path      int  true  "Book ID"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /favorites/{book_id} [post]
func (h *FavoriteHandler) Add(c *gin.Context) {
	userID := c.GetUint("userID")
	bookID, err := strconv.ParseUint(c.Param("book_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	if err := h.repo.Add(userID, uint(bookID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Added to favorites"})
}

// Remove godoc
// @Summary      Remove book from favorites
// @Description  Remove a specific book from the authenticated user's favorites list
// @Tags         favorites
// @Security     ApiKeyAuth
// @Param        book_id  path      int  true  "Book ID"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /favorites/{book_id} [delete]
func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID := c.GetUint("userID")
	bookID, err := strconv.ParseUint(c.Param("book_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	if err := h.repo.Remove(userID, uint(bookID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}

// List godoc
// @Summary      List favorite books
// @Description  Get a list of all favorite books for the authenticated user
// @Tags         favorites
// @Security     ApiKeyAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /favorites [get]
func (h *FavoriteHandler) List(c *gin.Context) {
	userID := c.GetUint("userID")
	items, err := h.repo.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"favorites": items})
}
