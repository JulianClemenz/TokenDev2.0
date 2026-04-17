package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
	UserID primitive.ObjectID `bson:"user_id,omitempty" json:"user_id" binding:"required"`
	//TokenID   primitive.ObjectID `bson:"token, omitempty" json:"token" binding:"required"`
	ExpiresAt time.Time `bson:"expires" json:"expires"`
	CreatedAt time.Time `bson:"created" json:"created"`
	IsActive  bool      `bson:"estatus" json:"status"`
}
