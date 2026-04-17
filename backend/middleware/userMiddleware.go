package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		//Permitimos el paso si es user o admin
		if !exists || (role != "user" && role != "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: se requieren permisos de usuario"})
			c.Abort()
			return
		}

		c.Next()
	}
}
