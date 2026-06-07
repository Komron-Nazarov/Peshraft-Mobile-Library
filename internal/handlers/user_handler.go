// package handlers

// import (
// 	"net/http"

// 	"mobile-library/internal/models"
// 	"mobile-library/internal/services"

// 	"github.com/gin-gonic/gin"
// )

// type UserHandler struct {
// 	service *services.UserService
// }

// func NewUserHandler(s *services.UserService) *UserHandler {
// 	return &UserHandler{service: s}
// }

// // GetProfile godoc
// // @Summary      Get user profile
// // @Description  Get the profile details of the currently authenticated user
// // @Tags         user
// // @Security     ApiKeyAuth
// // @Produce      json
// // @Success      200  {object}  map[string]interface{}
// // @Failure      404  {object}  map[string]string
// // @Router       /user/profile [get]
// func (h *UserHandler) GetProfile(c *gin.Context) {
// 	userID := c.GetUint("userID")
// 	profile, err := h.service.GetProfile(userID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"profile": profile})
// }

// // UpdateProfile godoc
// // @Summary      Update user profile
// // @Description  Update profile information for the authenticated user
// // @Tags         user
// // @Security     ApiKeyAuth
// // @Accept       json
// // @Produce      json
// // @Param        body  body      models.UpdateProfileRequest  true  "Update profile details"
// // @Success      200   {object}  map[string]interface{}
// // @Failure      400   {object}  map[string]string
// // @Router       /user/profile [put]
// func (h *UserHandler) UpdateProfile(c *gin.Context) {
// 	userID := c.GetUint("userID")
// 	var req models.UpdateProfileRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	profile, err := h.service.UpdateProfile(userID, req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Profile updated", "profile": profile})
// }

// // ChangePassword godoc
// // @Summary      Change password
// // @Description  Change the password for the authenticated user
// // @Tags         user
// // @Security     ApiKeyAuth
// // @Accept       json
// // @Produce      json
// // @Param        body  body      models.ChangePasswordRequest  true  "Password change details"
// // @Success      200   {object}  map[string]string
// // @Failure      400   {object}  map[string]string
// // @Router       /user/change-password [post]
// func (h *UserHandler) ChangePassword(c *gin.Context) {
// 	userID := c.GetUint("userID")
// 	var req models.ChangePasswordRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := h.service.ChangePassword(userID, req); err != nil {
// 		status := http.StatusInternalServerError
// 		if err.Error() == "incorrect old password" {
// 			status = http.StatusBadRequest
// 		}
// 		c.JSON(status, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
// }

// // GetHistory godoc
// // @Summary      Get borrow history
// // @Description  Get the book borrowing history for the authenticated user
// // @Tags         user
// // @Security     ApiKeyAuth
// // @Produce      json
// // @Success      200  {object}  map[string]interface{}
// // @Failure      500  {object}  map[string]string
// // @Router       /user/history [get]
// func (h *UserHandler) GetHistory(c *gin.Context) {
// 	userID := c.GetUint("userID")
// 	borrowService := c.MustGet("borrowService").(*services.BorrowService)
// 	history, err := borrowService.GetHistory(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	history_data := history
// 	c.JSON(http.StatusOK, gin.H{"history": history_data})
// }



// package handlers

// import (
//     "net/http"
//     "mobile-library/internal/models"
//     "mobile-library/internal/services"
//     "github.com/gin-gonic/gin"
// )

// type UserHandler struct {
//     service *services.UserService
// }

// func NewUserHandler(s *services.UserService) *UserHandler {
//     return &UserHandler{service: s}
// }

// func (h *UserHandler) GetProfile(c *gin.Context) {
//     userID := c.GetUint("userID")
//     profile, err := h.service.GetProfile(userID)
//     if err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//         return
//     }
//     c.JSON(http.StatusOK, gin.H{"profile": profile})
// }

// // UpdateProfile godoc
// func (h *UserHandler) UpdateProfile(c *gin.Context) {
//     userID := c.GetUint("userID")
//     var req models.UpdateProfileRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     // Передаем обновленный req, который теперь включает DateOfBirth
//     profile, err := h.service.UpdateProfile(userID, req)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }
//     c.JSON(http.StatusOK, gin.H{"message": "Profile updated", "profile": profile})
// }

// func (h *UserHandler) ChangePassword(c *gin.Context) {
//     userID := c.GetUint("userID")
//     var req models.ChangePasswordRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
//     if err := h.service.ChangePassword(userID, req); err != nil {
//         status := http.StatusInternalServerError
//         if err.Error() == "incorrect old password" {
//             status = http.StatusBadRequest
//         }
//         c.JSON(status, gin.H{"error": err.Error()})
//         return
//     }
//     c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
// }

// func (h *UserHandler) GetHistory(c *gin.Context) {
//     userID := c.GetUint("userID")
//     borrowService := c.MustGet("borrowService").(*services.BorrowService)
//     history, err := borrowService.GetHistory(userID)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }
//     c.JSON(http.StatusOK, gin.H{"history": history})
// }
package handlers

import (
	"net/http"
	"mobile-library/internal/models"
	"mobile-library/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Вспомогательная функция безопасного извлечения userID из контекста
func getAuthUserID(c *gin.Context) uint {
	if id, exists := c.Get("userID"); exists {
		if uIntId, ok := id.(uint); ok { return uIntId }
		if intId, ok := id.(int); ok { return uint(intId) }
		if floatId, ok := id.(float64); ok { return uint(floatId) }
	}
	if id, exists := c.Get("user_id"); exists {
		if uIntId, ok := id.(uint); ok { return uIntId }
		if intId, ok := id.(int); ok { return uint(intId) }
		if floatId, ok := id.(float64); ok { return uint(floatId) }
	}
	return c.GetUint("userID") // Фоллбэк на стандартный метод
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := getAuthUserID(c)
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// UpdateProfile godoc
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := getAuthUserID(c)
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	profile, err := h.service.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated", "profile": profile})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := getAuthUserID(c)
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ChangePassword(userID, req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "incorrect old password" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *UserHandler) GetHistory(c *gin.Context) {
	userID := getAuthUserID(c)
	borrowService := c.MustGet("borrowService").(*services.BorrowService)
	history, err := borrowService.GetHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"history": history})
}