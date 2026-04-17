package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//"/admin", middleware.CheckAdmin()  poner estas rutas en el main

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, ok := c.Get("role")
		if !ok || result.(string) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: se requiere rol de administrador"})
			return
		}

		c.Next()
	}
}
