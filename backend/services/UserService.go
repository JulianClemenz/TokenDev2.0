package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	PostUser(user *dto.UserRegisterDTO) (*dto.UserResponseDTO, error)
	GetUsers() ([]*dto.UserResponseDTO, error)
	GetUserByID(id string) (*dto.UserResponseDTO, error)
	PutUser(user *dto.UserModifyDTO) (*dto.UserModifyResponseDTO, error)
	PasswordModify(dto dto.PasswordChange, id string) (bool, error)
	//DELETE?
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserService(UserRepository repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: UserRepository,
	}
}

func (service *UserService) PostUser(userDto *dto.UserRegisterDTO) (*dto.UserResponseDTO, error) {
	//VALIDACIONES
	if userDto.Weight < 0 { //comprobamos que el peso ingresado no sea negativo
		return nil, fmt.Errorf("tu peso no puede ser menor a 0")
	}

	if userDto.BirthDate.After(time.Now()) { //comprobamos que la feha de nacimiento no sea mayor a hoy
		return nil, fmt.Errorf("error en fecha de nacimiento")
	}

	userDto.Email = strings.ToLower(strings.TrimSpace(userDto.Email))
	userDto.UserName = strings.TrimSpace(userDto.UserName)

	verificacionEmail, err := service.UserRepository.ExistByEmail(userDto.Email) //verifica si existe el email en la base de datos
	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar email: %w", err)
	}

	if verificacionEmail {
		return nil, fmt.Errorf("dicho email ya existe")
	}

	verificacionUserName, err := service.UserRepository.ExistByUserName(userDto.UserName) //verifica si existe el username en la base de datos

	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar nombre de usuario: %w", err)
	}

	if verificacionUserName {
		return nil, fmt.Errorf("dicho nombre de usuario ya existe")
	}

	//LOGICA
	userDB := userDto.GetModelUserRegister() //convertimos el dto para registrar en model

	hashed, err := utils.HashPassword(userDB.Password)
	if err != nil { //comprobamos que no suceda ningun error en el hasheo de la contraseña
		return nil, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	userDB.Password = hashed //hasheamos la contraseña

	result, err := service.UserRepository.PostUser(userDB)
	if err != nil {
		return nil, fmt.Errorf("error al insertar usuario: %w", err)
	}

	userDB.ID = result.InsertedID.(primitive.ObjectID) //asignamos el ID generado por Mongo al userDB
	userResponse := dto.NewUserResponseDTO(userDB)     //convertimos el model a dto para devolverlo
	return userResponse, nil
}

func (services *UserService) GetUsers() ([]*dto.UserResponseDTO, error) {
	usersDB, err := services.UserRepository.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}

	var users []*dto.UserResponseDTO
	for _, userDB := range usersDB {
		user := dto.NewUserResponseDTO(userDB)
		users = append(users, user)
	}

	return users, nil
}

func (services *UserService) GetUserByID(id string) (*dto.UserResponseDTO, error) {
	userDB, err := services.UserRepository.GetUsersByID(id)

	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}
	return dto.NewUserResponseDTO(userDB), nil
}

func (s *UserService) PutUser(newData *dto.UserModifyDTO) (*dto.UserModifyResponseDTO, error) {

	if newData == nil {
		return nil, fmt.Errorf("datos vacíos")
	}
	user, err := s.UserRepository.GetUsersByID(newData.ID)

	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario")
	}
	if user.ID.IsZero() {
		return nil, fmt.Errorf("no existe ningun usuario con ese ID") //verificamos que existe el usuario antes de modificar
	}

	newData.UserName = strings.TrimSpace(newData.UserName)

	if newData.UserName != strings.TrimSpace(user.UserName) { //solo verificamos si el user name es diferente al q ya estaba
		exist, err := s.UserRepository.ExistByUserNameExceptID(newData.ID, newData.UserName)
		if err != nil {
			return nil, fmt.Errorf("no se pudo verificar nombre de usuario: %w", err)
		}

		if exist {
			return nil, fmt.Errorf("dicho nombre de usuario ya existe")
		}
	}

	userDB, err := dto.GetModelUserModify(newData)
	if err != nil {
		return nil, err
	}
	userResp := dto.NewUserModifyResponseDTO(userDB)

	if _, err := s.UserRepository.PutUser(userDB); err != nil {
		return nil, fmt.Errorf("error al modificar usuario: %w", err)
	}

	return userResp, nil
}

func (s *UserService) PasswordModify(dto dto.PasswordChange, id string) (bool, error) {
	current := strings.TrimSpace(dto.CurrentPassword)
	newPassword := strings.TrimSpace(dto.NewPassword)
	confirm := strings.TrimSpace(dto.ConfirmPassword)

	userDB, err := s.UserRepository.GetUsersByID(id)
	if err != nil {
		return false, err
	}

	if userDB.ID.IsZero() {
		return false, fmt.Errorf("no se encontro ningun usuario con ese ID")
	}
	// verificamos contraseña actual
	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(current))
	if err != nil {
		return false, fmt.Errorf("contraseña no coincide")
	}
	// Evitar reutilizar la misma
	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(newPassword))
	if err == nil {
		return false, fmt.Errorf("la nueva contraseña no pued ser igual a la anterior")
	}
	// Confirmación
	if newPassword != confirm {
		return false, fmt.Errorf("la nueva contraseña y su confirmacion no se iguales")
	}

	// Hashear nueva
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return false, fmt.Errorf("error al hashear contraseña")
	}

	dto.NewPassword = hashed

	modified, err := s.UserRepository.UpdateNewPassword(dto, id)
	if err != nil {
		return false, err
	}

	if modified == 0 {
		return false, fmt.Errorf("no se realizaron cambios")
	}

	return true, nil
}
