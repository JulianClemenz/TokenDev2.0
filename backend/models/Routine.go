package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Routine struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name            string               `bson:"name" json:"name" binding:"required"`
	CreatorUserID   primitive.ObjectID   `bson:"creator_user_id" json:"creator_user_id"`
	ExcerciseList   []ExcerciseInRoutine `bson:"exercise_list,omitempty" json:"exercise_list,omitempty"` // no required
	EditionDate     time.Time            `bson:"edition_date" json:"edition_date"`
	EliminationDate time.Time            `bson:"elimination_date" json:"elimination_date"`
	CreationDate    time.Time            `bson:"creation_date" json:"creation_date"`
}

type ExcerciseInRoutine struct {
	ExcerciseID     primitive.ObjectID `bson:"excercise_id,omitempty" json:"excercise_id" binding:"required"`
	Repetitions     int                `bson:"repetitions"  json:"repetitions"  binding:"required,min=1"`
	Series          int                `bson:"series"       json:"series"       binding:"required,min=1"`
	Weight          float64            `bson:"weight"       json:"weight"       binding:"gte=0"`
	EliminationDate time.Time          `bson:"elimination_date" json:"elimination_date"`
	CreationDate    time.Time          `bson:"creation_date" json:"creation_date"`
}
