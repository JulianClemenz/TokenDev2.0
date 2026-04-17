package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutInterface interface {
	PostWorkout(*dto.WorkoutRegisterDTO) (*dto.WorkoutResponseDTO, error)
	GetWorkouts(idUser string) ([]*dto.WorkoutResponseDTO, error)
	GetWorkoutByID(workoutID string, userID string) (*dto.WorkoutResponseDTO, error)
	DeleteWorkout(dto.WorkoutDeleteDTO) error
	GetWorkoutStats(userID string) (*dto.WorkoutStatsDTO, error)
}

type WorkoutService struct {
	WorkoutRepository repositories.WorkoutRepositoryInterface
	RoutineRepository repositories.RoutineRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewWorkoutService(workoutRepository repositories.WorkoutRepositoryInterface, routineRepository repositories.RoutineRepositoryInterface, userRepository repositories.UserRepositoryInterface) *WorkoutService {
	return &WorkoutService{
		WorkoutRepository: workoutRepository,
		RoutineRepository: routineRepository,
		UserRepository:    userRepository,
	}
}

func (ws WorkoutService) PostWorkout(workoutDTO *dto.WorkoutRegisterDTO /*UserID se setea en handler*/) (*dto.WorkoutResponseDTO, error) {
	result, err := ws.RoutineRepository.GetRoutineByID(workoutDTO.RoutineID)
	if err != nil {
		return nil, fmt.Errorf("rutina no encontrada: %w", err)
	}
	if result.ID.IsZero() {
		return nil, fmt.Errorf("rutina no encontrada")
	}

	workoutModel, err := dto.GetModelWorkoutRegisterDTO(workoutDTO)
	if err != nil {
		return nil, err
	}
	workoutModel.Date = time.Now()
	workoutModel.RoutineName = result.Name

	insertResult, err := ws.WorkoutRepository.PostWorkout(workoutModel)
	if err != nil {
		return nil, fmt.Errorf("error al crear el workout: %w", err)
	}
	if insertResult.InsertedID == nil {
		return nil, fmt.Errorf("no se pudo crear el workout")
	}

	//recuperar workout creado
	//createdWorkout, err := ws.WorkoutRepository.GetWorkoutByID(insertResult.InsertedID.(string))  este era el original, gemini me dijo que no era correcto y se cambio por las dos lineas de abajo
	insertedID := insertResult.InsertedID.(primitive.ObjectID)
	createdWorkout, err := ws.WorkoutRepository.GetWorkoutByID(insertedID.Hex())

	if err != nil {
		return nil, fmt.Errorf("error al obtener el workout creado: %w", err)
	}
	if createdWorkout.ID.IsZero() {
		return nil, fmt.Errorf("workout creado no encontrado")
	}

	//convertir a dto y devolver
	workoutResponse := dto.NewWorkoutResponseDTO(createdWorkout)

	return workoutResponse, nil
}

// GetWorkouts obtiene todos los workouts de un usuario específico
func (ws WorkoutService) GetWorkouts(idUser string) ([]*dto.WorkoutResponseDTO, error) {

	//validar existencia de user
	user, err := ws.UserRepository.GetUsersByID(idUser)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}
	if user.ID.IsZero() {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	//obtener workouts del user
	workoutsModel, err := ws.WorkoutRepository.GetWorkoutsByUserID(idUser)
	if err != nil {
		return nil, fmt.Errorf("error al obtener workouts: %w", err)
	}
	if len(workoutsModel) == 0 {
		return nil, fmt.Errorf("no se encontraron workouts para el usuario")
	}

	var workoutsDTO []*dto.WorkoutResponseDTO
	for _, workout := range workoutsModel {
		workoutDTO := dto.NewWorkoutResponseDTO(workout)
		workoutsDTO = append(workoutsDTO, workoutDTO)
	}

	return workoutsDTO, nil
}

// GetWorkoutByID obtiene un workout por su ID
func (ws WorkoutService) GetWorkoutByID(workoutID string, userID string) (*dto.WorkoutResponseDTO, error) {
	//obtener workout por id
	workoutModel, err := ws.WorkoutRepository.GetWorkoutByID(workoutID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener workout: %w", err)
	}
	if workoutModel.ID.IsZero() {
		return nil, fmt.Errorf("workout no encontrado")
	}
	// (Opcional: aquí podrías añadir una validación para ver si el usuario logueado es el dueño)
	userIDconvert := utils.GetStringIDFromObjectID(workoutModel.UserID)
	if userIDconvert != userID {
		return nil, fmt.Errorf("el usuario no tiene permiso para acceder a workouts")
	}

	workoutDTO := dto.NewWorkoutResponseDTO(workoutModel)
	return workoutDTO, nil
}

func (ws WorkoutService) DeleteWorkout(delete dto.WorkoutDeleteDTO) error {
	//validacion de existencia de workout
	workout, err := ws.WorkoutRepository.GetWorkoutByID(delete.RoutineID)
	if err != nil {
		return fmt.Errorf("error al obtener workout: %w", err)
	}
	if workout.ID.IsZero() {
		return fmt.Errorf("workout no encontrado")
	}
	//validamos si dicho editor es dueño de la rutina a eliminar
	if delete.UserID != workout.UserID.Hex() {
		return fmt.Errorf("al no ser el creador de dicho workout no tienes permisos para esta accion")
	}

	//eliminar workout
	result, err := ws.WorkoutRepository.DeleteWorkout(delete.RoutineID)
	if err != nil {
		return fmt.Errorf("error al eliminar workout: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no se pudo eliminar el workout")
	}
	return nil
}

func (ws WorkoutService) GetWorkoutStats(userID string) (*dto.WorkoutStatsDTO, error) {

	// validacion de existencia de user
	userModel, err := ws.UserRepository.GetUsersByID(userID)
	if err != nil {
		return nil, fmt.Errorf("Error al recuperar usuario")
	}
	if userModel.ID.IsZero() {
		return nil, fmt.Errorf("No se encontro user")
	}

	// logica
	workoutsUser, err := ws.WorkoutRepository.GetWorkoutsByUserID(userID) // lista de workouts del user
	if err != nil {
		return nil, fmt.Errorf("Error al obtener workouts")
	}

	// Inicializamos el DTO con valores por defecto para evitar nulos en el JSON
	status := &dto.WorkoutStatsDTO{
		TotalWorkouts:    len(workoutsUser),
		MostUsedRoutines: []dto.RoutineUsageDTO{},
		ProgressOverTime: []dto.ProgressPointDTO{},
		WeeklyFrequency:  0.0,
	}

	// Si hay 0 o 1 workout, devolvemos lo básico porque no se puede calcular frecuencia entre fechas
	if len(workoutsUser) <= 1 {
		return status, nil
	}

	// (logica: (cantidad de dias entre el primer y el ult workots - 1) / cantidad de entrenamientos)
	sort.Slice(workoutsUser, func(i, j int) bool { // ordena ascendentemente por fecha
		return workoutsUser[i].Date.Before(workoutsUser[j].Date)
	})

	primerWorkout := workoutsUser[0]
	ultWorkout := workoutsUser[len(workoutsUser)-1]

	dayDifference := ultWorkout.Date.Sub(primerWorkout.Date).Hours() / 24

	if dayDifference < 1 {
		// Si todos los entrenamientos fueron el mismo día (diferencia < 1 día)
		status.WeeklyFrequency = float64(status.TotalWorkouts)
	} else {
		// Fórmula ajustada: Workouts totales dividido por la cantidad de semanas que han pasado
		status.WeeklyFrequency = float64(status.TotalWorkouts) / (dayDifference / 7.0)
	}

	// --- MostUsedRoutines (ranking de rutinas mas usadas) ---
	counts := make(map[string]int)
	for _, w := range workoutsUser {
		counts[w.RoutineName]++
	}

	for name, c := range counts {
		status.MostUsedRoutines = append(status.MostUsedRoutines, dto.RoutineUsageDTO{
			RoutineName: name,
			Count:       c,
		})
	}
	// ordenamos primero por mas usados, segundo alfabeticamente
	sort.Slice(status.MostUsedRoutines, func(i, j int) bool {
		if status.MostUsedRoutines[i].Count == status.MostUsedRoutines[j].Count {
			return status.MostUsedRoutines[i].RoutineName < status.MostUsedRoutines[j].RoutineName
		}
		return status.MostUsedRoutines[i].Count > status.MostUsedRoutines[j].Count
	})

	// ---grafica ---
	buckets := make(map[string]int)
	for _, w := range workoutsUser {
		key := w.Date.Format("2006-01") // yyyy-mm
		buckets[key]++
	}

	// Pasar a slice ordenado cronológicamente
	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		status.ProgressOverTime = append(status.ProgressOverTime, dto.ProgressPointDTO{
			Date:  k,
			Count: buckets[k],
		})
	}

	return status, nil
}
