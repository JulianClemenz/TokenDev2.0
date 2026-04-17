package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Leemos el rol que seteó el AuthMiddleware
		role, exists := c.Get("role")

		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: se requieren permisos de administrador"})
			c.Abort()
			return
		}

		c.Next()
	}
}
