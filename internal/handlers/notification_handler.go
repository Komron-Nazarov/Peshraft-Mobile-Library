package handlers

import (
	"net/http"
	"strconv"

	"mobile-library/internal/models"
	"mobile-library/internal/services"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(s *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

// List godoc
// @Summary      Get user notifications
// @Description  Get all notifications and unread count for the authenticated user
// @Tags         notifications
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	userID := val.(uint)

	notifs, unread, err := h.service.GetUserNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	if notifs == nil {
		notifs = []models.NotificationResponse{}
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifs,
		"unread_count":  unread,
	})
}

// MarkRead godoc
// @Summary      Mark notification as read
// @Description  Update the status of a specific notification to read
// @Tags         notifications
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Notification ID"
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /notifications/{id}/read [post]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	userID := val.(uint)

	notifID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	if err := h.service.MarkAsRead(uint(notifID), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied or not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}
