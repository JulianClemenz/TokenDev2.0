package services

import (
	"AppFitness/dto"
	"AppFitness/models"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthInterface interface {
	Login(loginDTO *dto.LoginRequestDTO) (*dto.LoginResponseDTO, error)
	Logout(refreshDTO *dto.RefreshRequestDTO) error
	Refresh(refreshDTO *dto.RefreshRequestDTO) (*dto.RefreshResponseDTO, error)
}

type AuthService struct {
	UserRepo    repositories.UserRepositoryInterface
	SessionRepo repositories.SessionRepositoryInterface
}

func NewAuthService(userRepository repositories.UserRepositoryInterface, sessionRepo repositories.SessionRepositoryInterface) AuthInterface {
	return &AuthService{
		UserRepo:    userRepository,
		SessionRepo: sessionRepo,
	}
}

func (s *AuthService) Login(loginDTO *dto.LoginRequestDTO) (*dto.LoginResponseDTO, error) {
	user, err := s.UserRepo.GetUserByEmail(loginDTO.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("credenciales inválidas: usuario no encontrado")
		}
		return nil, fmt.Errorf("error al buscar usuario: %w", &err)
	}

	isValidPassword := utils.CheckPasswordHash(loginDTO.Password, user.Password)
	if !isValidPassword {
		return nil, fmt.Errorf("credenciales inválidas: contraseña incorrecta")
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("error al generar el access token: %w", err)
	}

	session := models.Session{
		ID:        primitive.NewObjectID(),            // ID único para sesión
		UserID:    user.ID,                            // Vinculamos la sesión al usuario
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7 días de duración
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	// Guardamos la sesión en la base de datos usando tu 'SessionRepository'
	_, err = s.SessionRepo.PostSession(session)
	if err != nil {
		return nil, fmt.Errorf("error al guardar la sesión: %w", err)
	}

	// El Refresh Token que devolvemos al usuario es el ID de la sesión en formato string
	refreshToken := session.ID.Hex()

	userResponse := dto.NewUserResponseDTO(user) // Creamos el DTO de usuario para el frontend

	return &dto.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResponse,
	}, nil
}

func (s *AuthService) Logout(refreshDTO *dto.RefreshRequestDTO) error {
	// El Refresh Token = ID de la sesión
	sessionID, err := primitive.ObjectIDFromHex(refreshDTO.RefreshToken)
	if err != nil {
		return fmt.Errorf("refresh token inválido")
	}

	// borramos la sesión de la base de datos
	deleteResult, err := s.SessionRepo.DeleteSession(sessionID.Hex())
	if err != nil {
		return fmt.Errorf("error al cerrar sesión: %w", err)
	}

	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("sesión no encontrada")
	}

	return nil
}

func (s *AuthService) Refresh(refreshDTO *dto.RefreshRequestDTO) (*dto.RefreshResponseDTO, error) {

	//Validar firmato Refresh
	sessionID, err := primitive.ObjectIDFromHex(refreshDTO.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token inválido: formato incorrecto")
	}

	// Buscamos la sesión en la base de datos
	session, err := s.SessionRepo.GetSessionByID(sessionID.Hex())
	if err != nil {
		return nil, fmt.Errorf("sesión inválida: no encontrada")
	}

	// Verificamos que la sesión no esté expirada
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("sesión expirada")
	}

	// Verificamos que esté activa
	if !session.IsActive {
		return nil, fmt.Errorf("sesión inactiva")
	}

	// Buscamos al usuario dueño de la sesión
	user, err := s.UserRepo.GetUsersByID(session.UserID.Hex())
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario de la sesión: %w", err)
	}

	// Usamos tu 'utils/jwt.go' para crear un nuevo token de corta duración
	newAccessToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("error al generar el nuevo access token: %w", err)
	}

	return &dto.RefreshResponseDTO{
		AccessToken: newAccessToken,
	}, nil
}
