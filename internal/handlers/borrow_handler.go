package handlers

import (
	"net/http"
	"strconv"

	"mobile-library/internal/services"

	"github.com/gin-gonic/gin"
)

type BorrowHandler struct {
	service *services.BorrowService
}

func NewBorrowHandler(s *services.BorrowService) *BorrowHandler {
	return &BorrowHandler{service: s}
}

// BorrowBook godoc
// @Summary      Borrow a book
// @Description  Borrow a book by its ID for the authenticated user
// @Tags         borrow
// @Security     ApiKeyAuth
// @Param         id   path      int  true  "Book ID"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Router       /borrow/{id} [post]
func (h *BorrowHandler) BorrowBook(c *gin.Context) {
	userID := c.GetUint("userID")
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	borrow, err := h.service.BorrowBook(userID, uint(bookID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Book borrowed", "borrow": borrow})
}

// ReturnBook godoc
// @Summary      Return a book
// @Description  Return a borrowed book by its ID
// @Tags         borrow
// @Security     ApiKeyAuth
// @Param         id   path      int  true  "Book ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Router       /borrow/{id}/return [post]
func (h *BorrowHandler) ReturnBook(c *gin.Context) {
	userID := c.GetUint("userID")
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	borrow, err := h.service.ReturnBook(userID, uint(bookID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book returned", "borrow": borrow})
}
