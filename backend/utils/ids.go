package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func GetObjectIDFromStringID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return objectID, nil
}

func GetStringIDFromObjectID(id primitive.ObjectID) string {
	return id.Hex()
}
