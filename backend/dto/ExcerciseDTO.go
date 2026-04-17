package dto

import (
	"AppFitness/models"
	"AppFitness/utils"
	"time"
)

// ExcerciseRegisterDTO
// ExcerciseResponseDTO
// ExcerciseModifyDTO

type ExcerciseRegisterDTO struct {
	CreatorUserID   string
	Name            string `json:"name" bson:"name" binding:"required"`
	Description     string `json:"description" bson:"description" binding:"required"`
	Category        string `json:"category" bson:"category" binding:"required"`
	MainMuscleGroup string `json:"main_muscle_group" bson:"main_muscle_group" binding:"required"`
	DifficultLevel  string `json:"difficult_level" bson:"difficult_level" binding:"required"`
	Example         string `json:"example" bson:"example" binding:"required"`
	Instructions    string `json:"instructions" bson:"instructions" binding:"required"`
}

func GetModelExcerciseRegister(excercise *ExcerciseRegisterDTO) *models.Excercise {
	return &models.Excercise{
		Name:            excercise.Name,
		Description:     excercise.Description,
		Category:        models.CategoryLevel(excercise.Category),
		MainMuscleGroup: excercise.MainMuscleGroup,
		DifficultLevel:  excercise.DifficultLevel,
		Example:         excercise.Example,
		Instructions:    excercise.Instructions,
	}
}

type ExcerciseResponseDTO struct {
	ID              string `json:"id"`
	Name            string
	Description     string
	CreatorUserID   string
	Category        string
	MainMuscleGroup string
	DifficultLevel  string
	Example         string
	Instructions    string
	EditionDate     time.Time
	EliminationDate time.Time
	CreationDate    time.Time
}

func NewExcerciseResponseDTO(excercise models.Excercise) *ExcerciseResponseDTO {
	return &ExcerciseResponseDTO{
		ID:              utils.GetStringIDFromObjectID(excercise.ID),
		Name:            excercise.Name,
		Description:     excercise.Description,
		CreatorUserID:   utils.GetStringIDFromObjectID(excercise.CreatorUserID),
		Category:        string(excercise.Category),
		MainMuscleGroup: excercise.MainMuscleGroup,
		DifficultLevel:  excercise.DifficultLevel,
		Example:         excercise.Example,
		Instructions:    excercise.Instructions,
		EditionDate:     excercise.EditionDate,
		EliminationDate: excercise.EliminationDate,
		CreationDate:    excercise.CreationDate,
	}
}

type ExcerciseModifyDTO struct {
	ID              string
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`
	Category        string `json:"category" binding:"required"`
	MainMuscleGroup string `json:"main_muscle_group" binding:"required"`
	DifficultLevel  string `json:"difficult_level" binding:"required"`
	Example         string `json:"example" binding:"required"`
	Instructions    string `json:"instructions" binding:"required"`
}

func GetModelExcerciseModify(excercise *ExcerciseModifyDTO) *models.Excercise {
	return &models.Excercise{
		Name:            excercise.Name,
		Description:     excercise.Description,
		Category:        models.CategoryLevel(excercise.Category),
		MainMuscleGroup: excercise.MainMuscleGroup,
		DifficultLevel:  excercise.DifficultLevel,
		Example:         excercise.Example,
		Instructions:    excercise.Instructions,
	}
}

type ExcerciseModifyResponseDTO struct {
	Name            string
	Description     string
	CreatorUserID   string
	Category        string
	MainMuscleGroup string
	DifficultLevel  string
	Example         string
	Instructions    string
	EditionDate     time.Time
}

func NewExcerciseModifyResponseDTO(excercise models.Excercise) *ExcerciseModifyResponseDTO {
	return &ExcerciseModifyResponseDTO{
		Name:            excercise.Name,
		Description:     excercise.Description,
		CreatorUserID:   utils.GetStringIDFromObjectID(excercise.CreatorUserID),
		Category:        string(excercise.Category),
		MainMuscleGroup: excercise.MainMuscleGroup,
		DifficultLevel:  excercise.DifficultLevel,
		Example:         excercise.Example,
		Instructions:    excercise.Instructions,
		EditionDate:     excercise.EditionDate,
	}
}

type ExerciseFilterDTO struct {
	Name        string `json:"name,omitempty"`
	Category    string `json:"category,omitempty"`
	MuscleGroup string `json:"muscle_group,omitempty"`
}
