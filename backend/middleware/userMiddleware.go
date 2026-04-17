package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//"/user", middleware.CheckUser() esto va en main.go

func CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, ok := c.Get("role")
		if !ok || result.(string) != "client" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: se requiere rol de usuario"})
			return
		}
		c.Next()
	}
}
