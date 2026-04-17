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

type ExcerciseInterface interface { //POST, PUT y DELETE son accesibles solo por admins (reciben actor)
	PostExcercise(excercise *dto.ExcerciseRegisterDTO) (*dto.ExcerciseResponseDTO, error)
	PutExcercise(newData *dto.ExcerciseModifyDTO) (*dto.ExcerciseModifyResponseDTO, error)
	DeleteExcercise(id string) (bool, error)
	GetExcercises() ([]*dto.ExcerciseResponseDTO, error)
	GetExcerciseByID(id string) (*dto.ExcerciseResponseDTO, error)
	GetByFilters(filterDTO dto.ExerciseFilterDTO) ([]*dto.ExcerciseResponseDTO, error)
}

type ExcerciseService struct {
	ExcerciseRepository repositories.ExcerciseRepositoryInterface
}

func NewExcerciseService(ExcerciseRepository repositories.ExcerciseRepositoryInterface) *ExcerciseService {
	return &ExcerciseService{
		ExcerciseRepository: ExcerciseRepository,
	}
}

func (service *ExcerciseService) PostExcercise(excerciseDto *dto.ExcerciseRegisterDTO) (*dto.ExcerciseResponseDTO, error) {

	// Validaciones de campos obligatorios
	if strings.TrimSpace(excerciseDto.Name) == "" {
		return nil, fmt.Errorf("el nombre del ejercicio no puede estar vacío")
	}
	if excerciseDto.DifficultLevel == "" {
		return nil, fmt.Errorf("el nivel de dificultad no puede estar vacío")
	}
	if excerciseDto.MainMuscleGroup == "" {
		return nil, fmt.Errorf("el grupo muscular no puede estar vacío")
	}
	if excerciseDto.Description == "" {
		return nil, fmt.Errorf("la descripción del ejercicio no puede estar vacía")
	}
	if excerciseDto.Category == "" {
		return nil, fmt.Errorf("la categoría del ejercicio no puede estar vacía")
	}

	// Validacion de existencia por nombre
	nameExist, err := service.ExcerciseRepository.ExistByName(excerciseDto.Name)
	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar el nombre del ejercicio en la base de datos: %w", err)
	}
	if nameExist == true {
		return nil, fmt.Errorf("ya existe un ejercicio con ese nombre")
	}

	idObjCreator, err := utils.GetObjectIDFromStringID(excerciseDto.CreatorUserID)
	if err != nil {
		return nil, fmt.Errorf("ID del creador con formato inválido: %w", err)
	}

	//LOGICA
	excerciseModel := dto.GetModelExcerciseRegister(excerciseDto)             //convertimos el dto a modelo para el repository
	excerciseModel.CreatorUserID = idObjCreator                               //asignamos el ObjectID del usuario que crea el ejercicio
	result, err := service.ExcerciseRepository.PostExcercise(*excerciseModel) //ejecutamos post en repository
	if err != nil {
		return nil, err
	}

	excerciseModel.ID = result.InsertedID.(primitive.ObjectID)        //asignamos el ID generado por MongoDB al model
	excerciseResponse := dto.NewExcerciseResponseDTO(*excerciseModel) //convertimos el modelo a dto para la respuesta
	return excerciseResponse, nil
}

func (service *ExcerciseService) PutExcercise(newData *dto.ExcerciseModifyDTO) (*dto.ExcerciseModifyResponseDTO, error) {
	// Validaciones de campos obligatorios
	if strings.TrimSpace(newData.Name) == "" {
		return nil, fmt.Errorf("datos vacios")
	}
	ObjetiveID, err := utils.GetObjectIDFromStringID(newData.ID)
	if err != nil {
		return nil, fmt.Errorf("formato de ID de ejercicio invalido: %w", err)
	}

	if ObjetiveID.IsZero() {
		return nil, fmt.Errorf("el id del ejercicio no puede estar vacío")
	}

	//LOGICA
	_, err = service.ExcerciseRepository.GetExcerciseByID(newData.ID) //comprobamos que el ejercicio a modificar existe
	if err != nil {
		return nil, err
	}

	excerciseModel := dto.GetModelExcerciseModify(newData) //convertimos el dto a modelo para el repository
	excerciseModel.ID = ObjetiveID                         //asignamos el ObjectID del ejercicio a modificar
	excerciseModel.EditionDate = time.Now()                //actualizamos la fecha de edición

	result, err := service.ExcerciseRepository.PutExcercise(*excerciseModel) //ejecutamos put en repository
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se modificó ningún ejercicio")
	}

	excerciseModify, err := service.ExcerciseRepository.GetExcerciseByID(newData.ID) //obtenemos el ejercicio modificado para devolverlo
	if err != nil {
		return nil, fmt.Errorf("error al obtener el ejercicio modificado: %w", err)
	}

	return dto.NewExcerciseModifyResponseDTO(excerciseModify), nil
}

func (service *ExcerciseService) GetExcercises() ([]*dto.ExcerciseResponseDTO, error) {
	excercisesDB, err := service.ExcerciseRepository.GetExcercises()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ejercicios: %w", err)
	}

	var excercises []*dto.ExcerciseResponseDTO
	for _, excerciseDB := range excercisesDB {
		excercise := dto.NewExcerciseResponseDTO(excerciseDB)
		excercises = append(excercises, excercise)
	}
	return excercises, nil
}

func (service *ExcerciseService) GetExcerciseByID(id string) (*dto.ExcerciseResponseDTO, error) {
	userDB, err := service.ExcerciseRepository.GetExcerciseByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener ejercicio: %w", err)
	}
	return dto.NewExcerciseResponseDTO(userDB), nil
}

func (service *ExcerciseService) GetByFilters(filterDTO dto.ExerciseFilterDTO) ([]*dto.ExcerciseResponseDTO, error) {
	if filterDTO.Name == "" && filterDTO.Category == "" && filterDTO.MuscleGroup == "" {
		return nil, fmt.Errorf("debe ingresar al menos un filtro de búsqueda (nombre, categoría o grupo muscular)")
	}

	excercisesDB, err := service.ExcerciseRepository.GetByFilters(filterDTO)
	if err != nil {
		return nil, fmt.Errorf("error al botener ejercicios aplicando filtros")
	}

	var excercises []*dto.ExcerciseResponseDTO
	for _, excerciseDB := range excercisesDB {
		excercise := dto.NewExcerciseResponseDTO(*excerciseDB)
		excercises = append(excercises, excercise)
	}
	return excercises, nil
}

func (service *ExcerciseService) DeleteExcercise(id string) (bool, error) {
	//VALIDACIONES
	result, err := service.ExcerciseRepository.GetExcerciseByID(id) //comprobamos que el ejercicio a eliminar existe
	if err != nil {
		return false, fmt.Errorf("error al obtener el ejercicio a eliminar: %w", err)
	}
	if result.ID.IsZero() {
		return false, fmt.Errorf("no existe un ejercicio con el id proporcionado")
	}

	//LOGICA
	deleteResult, err := service.ExcerciseRepository.DeleteExcercise(id) //ejecutamos delete en repository
	if err != nil {
		return false, err
	}
	if deleteResult.DeletedCount == 0 {
		return false, fmt.Errorf("no se eliminó ningún ejercicio")
	}
	return true, nil
}
