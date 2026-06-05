// package handlers

// import (
// 	"net/http"

// 	"mobile-library/internal/models"
// 	"mobile-library/internal/services"

// 	"github.com/gin-gonic/gin"
// )

// type AuthHandler struct {
// 	service *services.AuthService
// }

// func NewAuthHandler(s *services.AuthService) *AuthHandler {
// 	return &AuthHandler{service: s}
// }

// // Register godoc
// // @Summary      User Registration
// // @Description  Create a new account in Shogun's Library with confirmation and birth date
// // @Tags         auth
// // @Accept       json
// // @Produce      json
// // @Param        input body models.RegisterRequest true "Registration Info"
// // @Success      201 {object} map[string]interface{} "Successfully registered"
// // @Failure      400 {object} map[string]interface{} "Invalid input data (e.g. passwords don't match)"
// // @Failure      500 {object} map[string]interface{} "Internal server error"
// // @Router       /auth/register [post]

// func (h *AuthHandler) Register(c *gin.Context) {
// 	var req models.RegisterRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
// 		return
// 	}

// 	user, token, err := h.service.Register(req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "Registration successful",
// 		"user":    user,
// 		"token":   token,
// 	})
// }

// // @Summary Login user
// // @Tags auth
// // @Accept json
// // @Produce json
// // @Param body body models.LoginRequest true "Login"
// // @Success 200 {object} map[string]interface{}
// // @Failure 401 {object} map[string]interface{}
// // @Router /auth/login [post]
// func (h *AuthHandler) Login(c *gin.Context) {
// 	var req models.LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
// 		return
// 	}

// 	user, token, err := h.service.Login(req)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Login successful",
// 		"user":    user,
// 		"token":   token,
// 	})
// }



package handlers

import (
    "net/http"
    "mobile-library/internal/models"
    "mobile-library/internal/services"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
    return &AuthHandler{service: s}
}

// @Summary      User Registration
// @Description  Create a new account with trim and lowercase email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.RegisterRequest true "Registration Info"
// @Success      201 {object} map[string]interface{}
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    user, token, err := h.service.Register(req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Registration successful",
        "user":    user,
        "token":   token,
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    user, token, err := h.service.Login(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "user":    user,
        "token":   token,
    })
}