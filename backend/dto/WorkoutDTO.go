package dto

import (
	"AppFitness/models"
	"AppFitness/utils"
	"fmt"
	"time"
)

type WorkoutRegisterDTO struct {
	RoutineID   string `json:"routine_id" binding:"required"`
	RoutineName string
	UserID      string
}
type WorkoutResponseDTO struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	RoutineID   string    `json:"routine_id"`
	RoutineName string    `json:"RoutineName"`
	DoneAt      time.Time `json:"DoneAt"`
}

func GetModelWorkoutRegisterDTO(dto *WorkoutRegisterDTO) (models.Workout, error) {
	routineOID, err := utils.GetObjectIDFromStringID(dto.RoutineID)
	if err != nil {
		return models.Workout{}, fmt.Errorf("ID de rutina con formato inválido: %w", err)
	}

	userOID, err := utils.GetObjectIDFromStringID(dto.UserID)
	if err != nil {
		// Esto no debería pasar si el token es válido, pero lo chequeamos
		return models.Workout{}, fmt.Errorf("ID de usuario con formato inválido: %w", err)
	}

	return models.Workout{
		RoutineID:   routineOID,
		UserID:      userOID,
		RoutineName: dto.RoutineName,
	}, nil // <--- 4. Devuelve nil como error
}

func NewWorkoutResponseDTO(workout models.Workout) *WorkoutResponseDTO {
	return &WorkoutResponseDTO{
		UserID:      utils.GetStringIDFromObjectID(workout.UserID),
		RoutineID:   utils.GetStringIDFromObjectID(workout.RoutineID),
		RoutineName: workout.RoutineName,
		DoneAt:      workout.Date,
	}
}

type WorkoutStatsDTO struct {
	TotalWorkouts    int                //cantidad total de workouts del user
	WeeklyFrequency  float64            // promedio de entrenamientos desde que se realizop el primero (ir contando la cantidad de dias que hay entre entrenamientos (desde el primero hasta el ult) y dividir por la cantidad de entrenamientos)
	MostUsedRoutines []RoutineUsageDTO  //ranking de rutinas mas usadas
	ProgressOverTime []ProgressPointDTO //para grafica entrenamientos-dias
}

type RoutineUsageDTO struct {
	RoutineName string
	Count       int
}

type ProgressPointDTO struct {
	Date  string
	Count int
}

type WorkoutDeleteDTO struct {
	RoutineID string
	UserID    string
}
