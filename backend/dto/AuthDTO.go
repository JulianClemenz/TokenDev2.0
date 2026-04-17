package dto

// LoginRequestDTO es lo que el usuario envía para iniciar sesión
type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponseDTO es lo que el servidor devuelve al loguearse
type LoginResponseDTO struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	User         *UserResponseDTO `json:"user"` // Para que el frontend sepa quién se logueó
}

// RefreshRequestDTO es lo que el usuario envía para refrescar su token
type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponseDTO es la respuesta con el nuevo access token
type RefreshResponseDTO struct {
	AccessToken string `json:"access_token"`
}
