package handlers

import (
	"AppFitness/dto"
	"AppFitness/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthInterface
}

func NewAuthHandler(authService services.AuthInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// PostLogin maneja la solicitud de inicio de sesión
func (h *AuthHandler) PostLogin(c *gin.Context) {
	var loginDTO dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	response, err := h.authService.Login(&loginDTO)
	if err != nil {
		// credenciales inválidas 401
		if strings.Contains(err.Error(), "credenciales inválidas") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		//  internos
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PostLogout maneja la solicitud de cierre de sesión
func (h *AuthHandler) PostLogout(c *gin.Context) {
	var refreshDTO dto.RefreshRequestDTO
	if err := c.ShouldBindJSON(&refreshDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token de refresco requerido: " + err.Error()})
		return
	}

	err := h.authService.Logout(&refreshDTO)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}

// PostRefresh maneja la solicitud de refresco de token
func (h *AuthHandler) PostRefresh(c *gin.Context) {
	var refreshDTO dto.RefreshRequestDTO
	if err := c.ShouldBindJSON(&refreshDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token de refresco requerido: " + err.Error()})
		return
	}

	response, err := h.authService.Refresh(&refreshDTO)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
