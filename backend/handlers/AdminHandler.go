package handlers

import (
	"AppFitness/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminService services.AdminInterface
}

func NewAdminHandler(adminService services.AdminInterface) *AdminHandler {
	return &AdminHandler{
		AdminService: adminService,
	}
}

func (h *AdminHandler) GetGlobalStats(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	result, err := h.AdminService.GetGlobalStats()
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "vacio"):
			// El service indica que no hay datos (len == 0)
			c.Status(http.StatusNoContent) // 204
			return

		case strings.Contains(msg, "error al recuperar rutinas"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener estadísticas globales"}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) GetLogs(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	result, count, err := h.AdminService.GetLogs()
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "vacio"):
			// No hay usuarios registrados
			c.Status(http.StatusNoContent) // 204
			return

		case strings.Contains(msg, "error al recuperar users"):
			// Error interno (DB o repositorio)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener logs de usuarios"}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	// Éxito: devolvemos la lista de usuarios + cantidad
	c.JSON(http.StatusOK, gin.H{
		"total": count,
		"users": result,
	})
}
