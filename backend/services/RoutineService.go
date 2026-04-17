package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutineInterface interface {
	PostRoutine(routineDTO *dto.RoutineRegisterDTO) (*dto.RoutineResponseDTO, error)
	GetRoutines() ([]*dto.RoutineResponseDTO, error)
	GetRoutineByID(id string) (*dto.RoutineResponseDTO, error)
	PutRoutine(modify dto.RoutineModifyDTO) (*dto.RoutineResponseDTO, error)
	AddExcerciseToRoutine(routineID string, exercise *dto.ExcerciseInRoutineDTO, idEditor string) (*dto.RoutineResponseDTO, error)
	RemoveExcerciseFromRoutine(idEditor string, remove dto.RoutineRemoveDTO) (*dto.RoutineResponseDTO, error)
	UpdateExerciseInRoutine(idEditor string, exerciseMod *dto.ExcerciseInRoutineModifyDTO) (*dto.RoutineResponseDTO, error)
	DeleteRoutine(id string, idEditor string) (bool, error)
}

type RoutineService struct {
	RoutineRepository   repositories.RoutineRepositoryInterface
	ExcerciseRepository repositories.ExcerciseRepositoryInterface
}

func NewRoutineService(routineRepository repositories.RoutineRepositoryInterface, excerciseRspository repositories.ExcerciseRepositoryInterface) *RoutineService {
	return &RoutineService{
		RoutineRepository:   routineRepository,
		ExcerciseRepository: excerciseRspository,
	}
}

func (service *RoutineService) PostRoutine(routineDTO *dto.RoutineRegisterDTO) (*dto.RoutineResponseDTO, error) {
	routineDTO.Name = strings.ToLower(strings.TrimSpace(routineDTO.Name))

	//validacion de campos obligatorios
	if routineDTO.Name == "" {
		return nil, fmt.Errorf("el nombre de la rutina no puede estar vacío")
	}
	if routineDTO.CreatorUserID == "" {
		return nil, fmt.Errorf("el ID del usuario creador no puede estar vacío")
	}

	verificacion, err := service.RoutineRepository.ExistByRutineName(routineDTO.Name)
	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar si existe una rutina con el mismo nombre: %w", err)
	}

	if verificacion {
		return nil, fmt.Errorf("dicho nombre de rutina ya existe")
	}

	//LOGICA
	model, err := dto.GetModelRoutineRegisterDTO(routineDTO)
	if err != nil {
		return nil, err
	}
	result, err := service.RoutineRepository.PostRoutine(*model) //insertamos la rutina en la base de datos
	if err != nil {
		return nil, fmt.Errorf("error al crear la rutina en RoutineService.PostRoutine(): %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("no se pudo obtener el ObjectID insertado")
	}
	idStr := utils.GetStringIDFromObjectID(oid)

	routineModel, err := service.RoutineRepository.GetRoutineByID(idStr) //obtenemos la rutina de tipo response creada para devolver

	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina creada en RoutineService.PostRoutine(): %v", err)
	}

	return dto.NewRoutineResponseDTO(*routineModel), nil
}

func (service *RoutineService) GetRoutines() ([]*dto.RoutineResponseDTO, error) {
	routinesDB, err := service.RoutineRepository.GetRoutines()
	if err != nil {
		return nil, fmt.Errorf("error al obtener rutinas %v:", err)
	}

	if len(routinesDB) == 0 {
		return nil, fmt.Errorf("no existen rutinas registradas")
	}

	var routines []*dto.RoutineResponseDTO
	for _, routineDB := range routinesDB {
		routine := dto.NewRoutineResponseDTO(*routineDB)
		routines = append(routines, routine)
	}

	return routines, nil
}

func (service *RoutineService) GetRoutineByID(id string) (*dto.RoutineResponseDTO, error) {
	routineDB, err := service.RoutineRepository.GetRoutineByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina por ID en RoutineService.GetRoutineByID(): %v", err)
	}
	if routineDB == nil || routineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con el ID proporcionado")
	}
	return dto.NewRoutineResponseDTO(*routineDB), nil
}

// PUT SOLO DE NAME
func (service *RoutineService) PutRoutine(modify dto.RoutineModifyDTO) (*dto.RoutineResponseDTO, error) {

	routineDB, err := service.RoutineRepository.GetRoutineByID(modify.IDRoutine)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina a modificar en RoutineService.PutRoutine(): %v", err)
	}
	if routineDB == nil || routineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID")
	}
	newName := strings.ToLower(strings.TrimSpace(modify.Name))
	if newName == "" {
		return nil, fmt.Errorf("el nombre de la rutina no puede estar vacío")
	}
	if strings.ToLower(strings.TrimSpace(routineDB.Name)) == newName {
		return nil, fmt.Errorf("el nuevo nombre de la rutina no puede ser igual al anterior")
	}

	routineDB.Name = strings.ToLower(strings.TrimSpace(modify.Name))
	routineDB.EditionDate = time.Now()

	result, err := service.RoutineRepository.PutRoutine(*routineDB)
	if err != nil {
		return nil, fmt.Errorf("error al modificar la rutina en RoutineService.PutRoutine(): %v", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se modificó ninguna rutina")
	}
	updatedRoutineDB, err := service.RoutineRepository.GetRoutineByID(modify.IDRoutine)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina modificada en RoutineService.PutRoutine(): %v", err)
	}
	return dto.NewRoutineResponseDTO(*updatedRoutineDB), nil
}

func (service *RoutineService) AddExcerciseToRoutine(routineID string, exercise *dto.ExcerciseInRoutineDTO, idEditor string) (*dto.RoutineResponseDTO, error) {

	routineDB, err := service.RoutineRepository.GetRoutineByID(routineID)
	if err != nil {
		return nil, err
	}
	if routineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID")
	}
	idCreator := utils.GetStringIDFromObjectID(routineDB.CreatorUserID)

	if idCreator != idEditor {
		return nil, fmt.Errorf("Al no ser el creador de esta rutina no se brinda permisos para dicha accion")
	}

	exerciseDB, err := service.ExcerciseRepository.GetExcerciseByID(exercise.ExcerciseID)
	if err != nil {
		return nil, err
	}
	if exerciseDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ningún ejercicio con ese ID")
	}

	exerciseModel, err := dto.GetModelExerciseInRoutineDTO(exercise)
	if err != nil {
		return nil, err
	}
	exerciseModel.CreationDate = time.Now()

	result, err := service.RoutineRepository.AddExerciseRutine(exerciseModel, routineDB.ID)
	if err != nil {
		return nil, fmt.Errorf("error al agregar el ejercicio a la rutina: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se agregó ningún ejercicio a la rutina")
	}

	updatedRoutineDB, err := service.RoutineRepository.GetRoutineByID(routineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina modificada en RoutineService.AddExcerciseToRoutine(): %v", err)
	}
	if updatedRoutineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID, error al agregar ejercicio")
	}
	return dto.NewRoutineResponseDTO(*updatedRoutineDB), nil
}

func (service *RoutineService) RemoveExcerciseFromRoutine(idEditor string, remove dto.RoutineRemoveDTO) (*dto.RoutineResponseDTO, error) {
	//validaciones
	routineDB, err := service.RoutineRepository.GetRoutineByID(remove.IDRoutine)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina a modificar: %w", err)
	}
	if routineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID")
	}
	idCreator := utils.GetStringIDFromObjectID(routineDB.CreatorUserID)

	if idCreator != idEditor {
		return nil, fmt.Errorf("Al no ser el creador de esta rutina no se brinda permisos para dicha accion")
	}
	routineObjectID, err := utils.GetObjectIDFromStringID(remove.IDRoutine)
	if err != nil {
		// Este es un error de formato de ID (400)
		return nil, fmt.Errorf("ID de rutina con formato inválido: %w", err)
	}

	exerciseDB, err := service.ExcerciseRepository.GetExcerciseByID(remove.IDExercise)
	if err != nil {
		return nil, err
	}
	if exerciseDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ningún ejercicio con ese ID")
	}
	exerciseObjectID, err := utils.GetObjectIDFromStringID(remove.IDExercise)
	if err != nil {
		// Este es un error de formato de ID (400)
		return nil, fmt.Errorf("ID de ejercicio con formato inválido: %w", err)
	}

	//lógica de eliminación
	result, err := service.RoutineRepository.DeleteExerciseToRutine(routineObjectID, exerciseObjectID)
	if err != nil {
		return nil, fmt.Errorf("error al eliminar el ejercicio de la rutina: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se eliminó ningún ejercicio de la rutina")
	}

	//modificar la fecha de edición de la rutina
	routineDB.EditionDate = time.Now()
	_, err = service.RoutineRepository.PutRoutine(*routineDB)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar la fecha de edición de la rutina en RoutineService.RemoveExcerciseFromRoutine(): %v", err)
	}

	//buscamos rutina para devolver
	updatedRoutineDB, err := service.RoutineRepository.GetRoutineByID(remove.IDRoutine)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina modificada en RoutineService.RemoveExcerciseFromRoutine(): %v", err)
	}
	if updatedRoutineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID, error al eliminar ejercicio")
	}
	return dto.NewRoutineResponseDTO(*updatedRoutineDB), nil
}

func (service *RoutineService) UpdateExerciseInRoutine(idEditor string, exerciseMod *dto.ExcerciseInRoutineModifyDTO) (*dto.RoutineResponseDTO, error) {

	//validaciones
	routineDB, err := service.RoutineRepository.GetRoutineByID(exerciseMod.RoutineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina a modificar: %w", err)
	}
	if routineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID")
	}
	idCreator := utils.GetStringIDFromObjectID(routineDB.CreatorUserID)

	if idCreator != idEditor {
		return nil, fmt.Errorf("Al no ser el creador de esta rutina no se brinda permisos para dicha accion")
	}
	routineObjectID, err := utils.GetObjectIDFromStringID(exerciseMod.RoutineID)
	if err != nil {
		return nil, fmt.Errorf("ID de rutina con formato inválido: %w", err)

	}

	exerciseDB, err := service.ExcerciseRepository.GetExcerciseByID(exerciseMod.ExcerciseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el ejercicio a modificar: %w", err)
	}
	if exerciseDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ningún ejercicio con ese ID")
	}
	exerciseObjectID, err := utils.GetObjectIDFromStringID(exerciseMod.ExcerciseID)
	if err != nil {
		return nil, fmt.Errorf("ID de ejercicio con formato inválido: %w", err)
	}

	//lógica de modificación
	exerciseModel := dto.GetModelFromExerciseInRoutineModifyDTO(exerciseMod)
	result, err := service.RoutineRepository.UpdateExerciseInRoutine(routineObjectID, exerciseObjectID, exerciseModel)
	if err != nil {
		return nil, fmt.Errorf("error al modificar el ejercicio de la rutina: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se modificó ningún ejercicio de la rutina")
	}

	//modificar la fecha de edición de la rutina
	routineDB.EditionDate = time.Now()
	_, err = service.RoutineRepository.PutRoutine(*routineDB)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar la fecha de edición de la rutina en RoutineService.UpdateExerciseInRoutine(): %v", err)
	}

	//buscamos rutina para devolver
	updatedRoutineDB, err := service.RoutineRepository.GetRoutineByID(exerciseMod.RoutineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina modificada en RoutineService.UpdateExerciseInRoutine(): %v", err)
	}
	if updatedRoutineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID, error al modificar ejercicio")
	}
	return dto.NewRoutineResponseDTO(*updatedRoutineDB), nil
}

func (service *RoutineService) DeleteRoutine(id string, idEditor string) (bool, error) {
	//validacion de existencia
	exist, err := service.RoutineRepository.GetRoutineByID(id)
	if err != nil {
		return false, fmt.Errorf("error al verificar la existencia de la rutina en RoutineService.DeleteRoutine(): %v", err)
	}
	if exist.ID.IsZero() {
		return false, fmt.Errorf("no existe ninguna rutina con ese ID")
	}
	idCreator := utils.GetStringIDFromObjectID(exist.CreatorUserID)

	if idCreator != idEditor {
		return false, fmt.Errorf("Al no ser el creador de esta rutina no se brinda permisos para dicha accion")
	}

	result, err := service.RoutineRepository.DeleteRoutine(id)
	if err != nil {
		return false, fmt.Errorf("error al eliminar la rutina en RoutineService.DeleteRoutine(): %v", err)
	}
	if result.DeletedCount == 0 {
		return false, fmt.Errorf("no se eliminó ninguna rutina")
	}
	return true, nil
}
