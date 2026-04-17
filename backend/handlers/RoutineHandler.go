package handlers

import (
	"AppFitness/dto"
	"AppFitness/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoutineHandler struct {
	RoutineService services.RoutineInterface
}

func NewRoutineHandler(routineService services.RoutineInterface) *RoutineHandler {
	return &RoutineHandler{
		RoutineService: routineService,
	}
}

func (h *RoutineHandler) PostRoutine(c *gin.Context) {
	idUser, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var routine dto.RoutineRegisterDTO
	if err := c.ShouldBindJSON(&routine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routine.CreatorUserID = idUser.(string)

	result, err := h.RoutineService.PostRoutine(&routine)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no puede estar vacío"),
			strings.Contains(msg, "no puede estar vacía"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) //400
			return

		case strings.Contains(msg, "dicho nombre de rutina ya existe"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) //409
			return

		case strings.Contains(msg, "no se pudo verificar si existe una rutina"),
			strings.Contains(msg, "error al crear la rutina"),
			strings.Contains(msg, "no se pudo obtener el ObjectID insertado"),
			strings.Contains(msg, "error al obtener la rutina creada"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al registrar la rutina"}) //500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) GetRoutines(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	result, err := h.RoutineService.GetRoutines()
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no existen rutinas registradas"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) //404
			return

		case strings.Contains(msg, "error al obtener rutinas"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener rutinas"}) //500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) GetRoutineByID(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"}) //401
		return
	}

	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "se requiere un ID de rutina para hacer la busqueda"})
		return
	}

	result, err := h.RoutineService.GetRoutineByID(id)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no existe ninguna rutina"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) //404
			return

		case strings.Contains(msg, "error al obtener la rutina por ID"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener la rutina"}) //500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) PutRoutine(c *gin.Context) {
	_, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	idRoutine := c.Param("id")

	var routineModify dto.RoutineModifyDTO
	if err := c.ShouldBindJSON(&routineModify); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routineModify.IDRoutine = idRoutine
	result, err := h.RoutineService.PutRoutine(routineModify)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no puede estar vacío"),
			strings.Contains(msg, "no puede ser igual al anterior"):
			c.JSON(http.StatusBadRequest, gin.H{"error": msg}) //400
			return

		case strings.Contains(msg, "no se modificó ninguna rutina"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) //409
			return

		case strings.Contains(msg, "error al obtener la rutina a modificar"),
			strings.Contains(msg, "error al modificar la rutina"),
			strings.Contains(msg, "error al obtener la rutina modificada"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al modificar la rutina"}) //500
			return

		case strings.Contains(msg, "no existe ninguna rutina con ese ID"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) //404
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) AddExcerciseToRoutine(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	idRoutine := c.Param("id")
	var exercise dto.ExcerciseInRoutineDTO
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.RoutineService.AddExcerciseToRoutine(idRoutine, &exercise, idEditor.(string))
	if err != nil {
		msg := err.Error()
		switch {
		// permisos
		case strings.Contains(msg, "Al no ser el creador de esta rutina"):
			c.JSON(http.StatusForbidden, gin.H{"error": msg}) // 403
			return

		// no encontrados
		case strings.Contains(msg, "no existe ninguna rutina con ese ID"),
			strings.Contains(msg, "no existe ningún ejercicio con ese ID"),
			strings.Contains(msg, "no se encontró ningún ejercicio"),
			strings.Contains(msg, "no se encontró ninguna rutina"),
			strings.Contains(msg, "no existe ninguna rutina con ese ID, error al agregar ejercicio"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		//no se agregó nada
		case strings.Contains(msg, "no se agregó ningún ejercicio a la rutina"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) // 409 (conflicto de negocio)
			return

		// errores DB
		case
			strings.Contains(msg, "error al obtener el ejercicio a agregar"),
			strings.Contains(msg, "error al agregar el ejercicio a la rutina"),
			strings.Contains(msg, "error al actualizar la fecha de edición de la rutina"),
			strings.Contains(msg, "error al obtener la rutina modificada"),
			strings.Contains(msg, "error al decodificar"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al agregar ejercicio a la rutina"}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) RemoveExerciseFromRoutine(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var exerciseRem dto.RoutineRemoveDTO
	if err := c.ShouldBindJSON(&exerciseRem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.RoutineService.RemoveExcerciseFromRoutine(idEditor.(string), exerciseRem)
	if err != nil {
		msg := err.Error()
		switch {
		// No se encontro
		case strings.Contains(msg, "no existe ninguna rutina con ese ID"),
			strings.Contains(msg, "no existe ningún ejercicio con ese ID"),
			strings.Contains(msg, "no existe ninguna rutina con ese ID, error al eliminar ejercicio"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		// no se eliminó nada
		case strings.Contains(msg, "no se eliminó ningún ejercicio de la rutina"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) // 409
			return

		// Errores en db
		case strings.Contains(msg, "error al obtener la rutina a modificar"),
			strings.Contains(msg, "error al obtener el ejercicio a eliminar"),
			strings.Contains(msg, "error al eliminar el ejercicio de la rutina"),
			strings.Contains(msg, "error al actualizar la fecha de edición de la rutina"),
			strings.Contains(msg, "error al obtener la rutina modificada"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al eliminar ejercicio de la rutina"}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) UpdateExerciseInRoutine(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	idExercise := c.Param("exercise_id")
	idRoutine := c.Param("id")

	var exerciseUpd dto.ExcerciseInRoutineModifyDTO
	if err := c.ShouldBindJSON(&exerciseUpd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exerciseUpd.ExcerciseID = idExercise
	exerciseUpd.RoutineID = idRoutine
	result, err := h.RoutineService.UpdateExerciseInRoutine(idEditor.(string), &exerciseUpd)
	if err != nil {
		msg := err.Error()
		switch {
		// No encontrados
		case strings.Contains(msg, "no existe ninguna rutina con ese ID"),
			strings.Contains(msg, "no existe ningún ejercicio con ese ID"),
			strings.Contains(msg, "no existe ninguna rutina con ese ID, error al modificar ejercicio"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		// Sin cambios
		case strings.Contains(msg, "no se modificó ningún ejercicio de la rutina"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) // 409
			return

		// Errores internos (repo/DB)
		case strings.Contains(msg, "error al obtener la rutina a modificar"),
			strings.Contains(msg, "error al obtener el ejercicio a modificar"),
			strings.Contains(msg, "error al modificar el ejercicio de la rutina"),
			strings.Contains(msg, "error al actualizar la fecha de edición de la rutina"),
			strings.Contains(msg, "error al obtener la rutina modificada"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al modificar ejercicio de la rutina"}) // 500
			return

		// Fallback
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *RoutineHandler) DeleteRoutine(c *gin.Context) {
	idEditor, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	idRoutine := c.Param("id")
	deleted, err := h.RoutineService.DeleteRoutine(idRoutine, idEditor.(string))
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no existe ninguna rutina con ese ID"):
			c.JSON(http.StatusNotFound, gin.H{"error": msg}) // 404
			return

		case strings.Contains(msg, "Al no ser el creador de esta rutina"):
			c.JSON(http.StatusForbidden, gin.H{"error": msg}) // 403
			return

		case strings.Contains(msg, "no se eliminó ninguna rutina"):
			c.JSON(http.StatusConflict, gin.H{"error": msg}) // 409
			return

		case strings.Contains(msg, "error al verificar la existencia de la rutina"),
			strings.Contains(msg, "error al eliminar la rutina"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al eliminar la rutina"}) // 500
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"deleted": deleted})
}
