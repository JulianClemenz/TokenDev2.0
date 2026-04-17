package handlers

import (
	"AppFitness/dto"
	"AppFitness/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	ExerciseService services.ExcerciseInterface
}

func NewExerciseHandler(exerciseService services.ExcerciseInterface) *ExerciseHandler {
	return &ExerciseHandler{
		ExerciseService: exerciseService,
	}
}

func (h *ExerciseHandler) GetByFilters(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	var filter dto.ExerciseFilterDTO
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exercises, err := h.ExerciseService.GetByFilters(filter)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "al menos un filtro"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) // 400
			return
		case strings.Contains(msg, "obtener ejercicios"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) // 500
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}
	c.JSON(http.StatusOK, exercises)

}

func (h *ExerciseHandler) GetExcerciseByID(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "se requiere un ID de ejercicio"})
		return
	}
	exercise, err := h.ExerciseService.GetExcerciseByID(id)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no se encontró"):
			c.JSON(http.StatusNotFound, gin.H{"error": "No existe un ejercicio con ese ID"}) // 404
			return
		case strings.Contains(msg, "obtener ejercicio"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al obtener el ejercicio"}) // 500
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) // 500 genérico
			return
		}
	}
	c.JSON(http.StatusOK, exercise)
}

func (h *ExerciseHandler) GetExcercises(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	collection, err := h.ExerciseService.GetExcercises()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error al obtener ejercicios"}) //404
		return
	}

	c.JSON(http.StatusOK, collection)
}

func (h *ExerciseHandler) PostExcercise(c *gin.Context) {
	idUser, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var exercise dto.ExcerciseRegisterDTO
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exercise.CreatorUserID = idUser.(string)

	resultado, err := h.ExerciseService.PostExcercise(&exercise)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no puede estar vacío"),
			strings.Contains(msg, "no puede estar vacía"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) //400
			return

		case strings.Contains(msg, "ya existe un ejercicio con ese nombre"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) //409
			return

		case strings.Contains(msg, "no se pudo verificar el nombre del ejercicio"),
			strings.Contains(msg, "error al insertar"),
			strings.Contains(msg, "insertar ejercicio"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al registrar ejercicio"}) //500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusCreated, resultado)

}

func (h *ExerciseHandler) PutExcercise(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el id del ejercico es requerido en la url"})
	}

	var exercise dto.ExcerciseModifyDTO

	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //400
		return
	}
	exercise.ID = id

	res, err := h.ExerciseService.PutExcercise(&exercise)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "datos vac"),
			strings.Contains(msg, "inválid"),
			strings.Contains(msg, "id del ejercicio no puede estar vacío"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) //400
			return

		case strings.Contains(msg, "no se encontró el ejercicio"),
			strings.Contains(msg, "no existe"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) //404
			return

		case strings.Contains(msg, "no se modificó ningún ejercicio"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) //409
			return

		case strings.Contains(msg, "obtener el ejercicio a modificar"),
			strings.Contains(msg, "obtener el ejercicio modificado"),
			strings.Contains(msg, "actualizar el ejercicio"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al modificar ejercicio"}) //500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *ExerciseHandler) DeleteExcercise(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	idExcercise := c.Param("id")
	if strings.TrimSpace(idExcercise) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere id de excercise"})
		return
	}

	deleted, err := h.ExerciseService.DeleteExcercise(idExcercise)
	if err != nil {
		// ... (Maneja los errores como en tus otros handlers, ej. 404 si no existe, 500 si falla)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se pudo eliminar el ejercicio o no fue encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ejercicio eliminado correctamente"})
}
