package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminLevel string

const (
	Admin  AdminLevel = "admin"
	Client AdminLevel = "client"
)

type ExperienceLevel string

const (
	Beginner     ExperienceLevel = "beginner"
	Intermediate ExperienceLevel = "intermediate"
	Advanced     ExperienceLevel = "advanced"
)

type ObjetiveLevel string

const (
	LoseWeight ObjetiveLevel = "lose_weight"
	GainWeight ObjetiveLevel = "gain_weight"
	Maintain   ObjetiveLevel = "maintain"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name" binding:"required"`
	LastName        string             `bson:"last_name" json:"last_name" binding:"required"`
	UserName        string             `bson:"user_name" json:"user_name" binding:"required"`
	Email           string             `bson:"email" json:"email" binding:"required,email"`
	Password        string             `bson:"password" json:"-"`
	BirthDate       time.Time          `bson:"birth_date" json:"birth_date"`
	Role            AdminLevel         `bson:"role" json:"role" binding:"required, oneof=admin client"`
	Weight          float32            `bson:"weight" json:"weight"`
	Height          float32            `bson:"height" json:"height"`
	Experience      ExperienceLevel    `bson:"experience" json:"experience" binding:"required, oneof=beginner intermediate advanced"`
	Objetive        ObjetiveLevel      `bson:"objetive" json:"objetive" binding:"required, oneof=lose_weight gain_weight maintain"`
	EditionDate     time.Time          `bson:"edition_date" json:"edition_date"`
	EliminationDate time.Time          `bson:"elimination_date" json:"elimination_date"`
	CreationDate    time.Time          `bson:"creation_date" json:"creation_date"`
}

// PARA FRONTEND
/*type LoginRequest struct {
	Email    string `bson:"email" binding:"required,email"`
	Password string `bson:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}*/
