package handlers

import (
	"AppFitness/dto"
	"AppFitness/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type WorkoutHandler struct {
	WorkoutService services.WorkoutInterface
}

func NewWorkoutHadler(workoutService services.WorkoutInterface) *WorkoutHandler {
	return &WorkoutHandler{
		WorkoutService: workoutService,
	}
}

func (h *WorkoutHandler) PostWorkout(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	idRoutine := c.Param("id_routine")
	newWorkout := &dto.WorkoutRegisterDTO{}

	newWorkout.RoutineID = idRoutine
	newWorkout.UserID = idEditor.(string)
	result, err := h.WorkoutService.PostWorkout(newWorkout)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "rutina no encontrada"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		case strings.Contains(msg, "error al crear el workout"),
			strings.Contains(msg, "no se pudo crear el workout"),
			strings.Contains(msg, "error al obtener el workout creado"),
			strings.Contains(msg, "workout creado no encontrado"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al crear el workout"}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusCreated, result)
}

func (h *WorkoutHandler) GetWorkouts(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	result, err := h.WorkoutService.GetWorkouts(idEditor.(string))
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "usuario no encontrado"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		case strings.Contains(msg, "no se encontraron workouts para el usuario"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return
			// Alternativa: 204 No Content si preferís no tratarlo como error

		case strings.Contains(msg, "error al obtener usuario"),
			strings.Contains(msg, "error al obtener workouts"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener workouts"}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *WorkoutHandler) GetWorkoutByID(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	idUser, ok := idEditor.(string)
	if !ok || strings.TrimSpace(idUser) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	idWork := c.Param("id")
	if strings.TrimSpace(idWork) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere id de workout"})
		return
	}

	result, err := h.WorkoutService.GetWorkoutByID(idWork, idUser)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "usuario no encontrado"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		case strings.Contains(msg, "workout no encontrado"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		case strings.Contains(msg, "error al obtener usuario"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) // 500
			return

		case strings.Contains(msg, "error al obtener workout"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) // 500
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *WorkoutHandler) DeleteWorkout(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	idWorkout := c.Param("id")

	var delete dto.WorkoutDeleteDTO
	delete.RoutineID = idWorkout
	delete.UserID = idEditor.(string)

	err := h.WorkoutService.DeleteWorkout(delete)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "workout no encontrado"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		case strings.Contains(msg, "al no ser el creador"):
			c.JSON(http.StatusForbidden, gin.H{"error": msg}) // 403
			return

		case strings.Contains(msg, "error al obtener workout"),
			strings.Contains(msg, "error al eliminar workout"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) // 500
			return

		case strings.Contains(msg, "no se pudo eliminar el workout"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) // 409
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) // 500
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout eliminado correctamente"})
}

func (h *WorkoutHandler) GetWorkoutStats(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	result, err := h.WorkoutService.GetWorkoutStats(idEditor.(string))
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "No se encontro user"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		case strings.Contains(msg, "Error al recuperar usuario"),
			strings.Contains(msg, "Error al obtener workouts"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener estadísticas"}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	// Nota: si el usuario tiene 0 o 1 workout, el service devuelve DTO vacío y err == nil → respondemos 200 OK igual.
	c.JSON(http.StatusOK, result)

}
