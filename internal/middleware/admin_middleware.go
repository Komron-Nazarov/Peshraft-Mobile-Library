package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        // userID мы уже получаем в AuthMiddleware, 
        // но нам нужно проверить роль.
        // В идеале роль нужно положить в context внутри AuthMiddleware.
        role, exists := c.Get("userRole") 
        if !exists || role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен: только для администраторов"})
            c.Abort()
            return
        }
        c.Next()
    }
}