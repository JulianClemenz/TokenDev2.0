package handlers

import (
	"AppFitness/dto"
	"AppFitness/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserInterface
}

func NewUserHandler(userService services.UserInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	collection, err := h.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"mensaje": "error al obtener usuarios"}) //Error 404
		return
	}

	c.JSON(http.StatusOK, collection)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no se encontró") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) //Error 404
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener cliente"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) PostUser(c *gin.Context) {
	var user dto.UserRegisterDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.userService.PostUser(&user) // dependiendo el error lanzamos un status
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "ya existe"): // email o username duplicados / eror 409
			c.JSON(http.StatusConflict, gin.H{"error": msg})
			return
		case strings.Contains(msg, "no se pudo verificar"),
			strings.Contains(msg, "error al insertar usuario"),
			strings.Contains(msg, "hashear contraseña"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al registrar usuario"}) //error 500
			return
		default:
			//  validaciones de negocio
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) // error 400
			return
		}
	}

	c.JSON(http.StatusCreated, resultado)
}

func (h *UserHandler) PutUser(c *gin.Context) {
	var user dto.UserModifyDTO

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //400
		return
	}

	user.ID = c.Param("id")
	result, err := h.userService.PutUser(&user)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "datos vacíos"),
			strings.Contains(msg, "inválid"),
			strings.Contains(msg, "no puede ser menor a 0"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) // error 400
			return

		case strings.Contains(msg, "ya existe"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) // eror 409
			return

		case strings.Contains(msg, "no existe ningun usuario"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) //Error 404
			return

		case strings.Contains(msg, "no se pudo verificar"),
			strings.Contains(msg, "error al buscar usuario"),
			strings.Contains(msg, "error al modificar usuario"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al modificar usuario"}) //error 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) PasswordModify(c *gin.Context) {
	var change dto.PasswordChange
	if err := c.ShouldBindJSON(&change); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //400
		return
	}
	id := c.Param("id")

	result, err := h.userService.PasswordModify(change, id)
	if err != nil {
		msg := err.Error()
		switch {
		// Autenticación
		case strings.Contains(msg, "contraseña no coincide"):
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg}) // 401
			return

		// Reglas de negocio
		case strings.Contains(msg, "la nueva contraseña no pued ser igual"),
			strings.Contains(msg, "la nueva contraseña y su confirmacion no se iguales"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) // 400
			return

		// No se encuentra
		case strings.Contains(msg, "no se encontro ningun usuario"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		// Errores internos
		case strings.Contains(msg, "error al obtener usuario por id"),
			strings.Contains(msg, "error al hashear contraseña"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al modificar contraseña"}) // 500
			return

		// Sin cambios
		case strings.Contains(msg, "no se realizaron cambios"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) // 409
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) // 500
			return
		}
	}

	if !result {
		c.JSON(http.StatusConflict, gin.H{"error": "Problema en lógica de negocio al intentar modificar contraseña"}) // 409
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "contraseña actualizada con exito"})

}
