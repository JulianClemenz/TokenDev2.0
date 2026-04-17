package utils

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret = []byte("firma_secretisima_del_token")

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwtv5.RegisteredClaims
}

func GenerateToken(userID primitive.ObjectID, email string, role string) (string, error) {
	claims := Claims{
		UserID: userID.Hex(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// También deberás actualizar tu función ValidateToken
func ValidateToken(tokenString string) (*Claims, error) {
	// Usar el alias `jwtv5` en todas partes
	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
