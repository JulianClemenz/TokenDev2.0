package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminInterface interface {
	GetGlobalStats() ([]*dto.TopUsedExcerciseDTO, error)
	GetLogs() ([]*dto.UserResponseDTO, int, error) //lista de users, cantidad
}

type AdminService struct {
	UserRepository      repositories.UserRepositoryInterface
	ExcerciseRepository repositories.ExcerciseRepositoryInterface
	RoutineRepository   repositories.RoutineRepositoryInterface
	SessionRepository   repositories.SessionRepositoryInterface
}

func NewAdminService(
	userRepo repositories.UserRepositoryInterface, exerciseRepo repositories.ExcerciseRepositoryInterface, routineRepo repositories.RoutineRepositoryInterface, sessionRepo repositories.SessionRepositoryInterface) AdminInterface {
	return &AdminService{
		UserRepository:      userRepo,
		ExcerciseRepository: exerciseRepo,
		RoutineRepository:   routineRepo,
		SessionRepository:   sessionRepo,
	}
}

func (a *AdminService) GetGlobalStats() ([]*dto.TopUsedExcerciseDTO, error) {

	routinesDB, err := a.RoutineRepository.GetRoutines()
	if err != nil {
		return nil, fmt.Errorf("error al recuperar rutinas")
	}
	if len(routinesDB) == 0 {
		return []*dto.TopUsedExcerciseDTO{}, nil // o nil, nil
	}
	//recuperamos todos los exc
	var excerciseList []primitive.ObjectID
	for _, r := range routinesDB {
		for _, e := range r.ExcerciseList {
			excerciseList = append(excerciseList, e.ExcerciseID)
		}
	}
	//agrupamos por id
	counts := make(map[primitive.ObjectID]int, len(excerciseList))
	for _, e := range excerciseList {
		counts[e]++
	}
	//mapear a dto TopUsedExcerciseDTO id, name, count

	topList := make([]*dto.TopUsedExcerciseDTO, 0, len(counts))
	for te, i := range counts {
		excerciseName, _ := a.ExcerciseRepository.GetExcerciseByID(utils.GetStringIDFromObjectID(te))
		topList = append(topList, &dto.TopUsedExcerciseDTO{
			ExcerciseID:   utils.GetStringIDFromObjectID(te),
			ExcerciseName: excerciseName.Name,
			Count:         i,
		})

	}
	//ordenamos primero por mas usados, segundo alfabeticamente
	sort.Slice(topList, func(i, j int) bool {
		if topList[i].Count == topList[j].Count {
			return topList[i].ExcerciseName < topList[j].ExcerciseName
		}
		return topList[i].Count > topList[j].Count
	})
	return topList, nil
}

func (a *AdminService) GetLogs() ([]*dto.UserResponseDTO, int, error) {

	//buscar users
	usersDB, err := a.UserRepository.GetUsers()
	if err != nil {
		return nil, 0, fmt.Errorf("error al recuperar users")
	}
	if len(usersDB) == 0 {
		// Devuelve vacío en lugar de error
		return []*dto.UserResponseDTO{}, 0, nil
	}

	var usersResponse []*dto.UserResponseDTO
	for _, u := range usersDB {
		userDTO := dto.NewUserResponseDTO(u)
		isActive, err := a.SessionRepository.IsUserActive(u.ID.Hex())
		if err != nil {
			// Si hay un error, inactivo por seguridad
			userDTO.IsActive = false
		} else {
			userDTO.IsActive = isActive // si esta activo le ponemos true sino false
		}

		usersResponse = append(usersResponse, userDTO)
	}

	return usersResponse, len(usersResponse), nil

}
