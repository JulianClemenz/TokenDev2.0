package repositories

import (
	"AppFitness/models"
	"AppFitness/utils"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutRepositoryInterface interface {
	PostWorkout(workout models.Workout) (*mongo.InsertOneResult, error)
	GetWorkouts() ([]models.Workout, error)
	GetWorkoutByID(id string) (models.Workout, error)
	GetWorkoutsByUserID(userID string) ([]models.Workout, error)
	PutWorkout(workout models.Workout) (*mongo.UpdateResult, error)
	DeleteWorkout(id string) (*mongo.DeleteResult, error)
}

type WorkoutRepository struct {
	db DB
}

func NewWorkoutRepository(db DB) *WorkoutRepository {
	return &WorkoutRepository{
		db: db,
	}
}

func (repository WorkoutRepository) PostWorkout(workout models.Workout) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("workouts")
	result, err := collection.InsertOne(nil, workout)
	if err != nil {
		return result, fmt.Errorf("error al insertar el workout en WorkoutRepository.PostWorkout(): %v", err)
	}
	return result, nil
}

func (repository WorkoutRepository) GetWorkouts() ([]models.Workout, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("workouts")
	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	var workouts []models.Workout
	for cursor.Next(context.Background()) {
		var workout models.Workout
		err := cursor.Decode(&workout)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar el workout en WorkoutRepository.GetWorkouts(): %v", err)
		}
		workouts = append(workouts, workout)
	}
	return workouts, err
}

func (repository WorkoutRepository) GetWorkoutByID(id string) (models.Workout, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("workouts")
	objectID, err := utils.GetObjectIDFromStringID(id)
	if err != nil {
		return models.Workout{}, err
	}
	filter := bson.M{"_id": objectID}

	result := collection.FindOne(context.TODO(), filter)

	var workout models.Workout
	err = result.Decode(&workout)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Workout{}, nil // Devuelve un workout vacío (ID.IsZero()) y SIN ERROR
		}
		return models.Workout{}, fmt.Errorf("error al obtener el workout en WorkoutRepository.GetWorkoutByID(): %v", err)
	}
	return workout, nil
}

func (repository WorkoutRepository) PutWorkout(workout models.Workout) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("workouts")
	filter := bson.M{"_id": workout.ID}
	entity := bson.M{"$set": workout}

	result, err := collection.UpdateOne(context.TODO(), filter, entity)
	if err != nil {
		return result, fmt.Errorf("error al actualizar el workout en WorkoutRepository.PutWorkout(): %v", err)
	}
	return result, nil
}

func (repository WorkoutRepository) DeleteWorkout(id string) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("workouts")
	objectID, err := utils.GetObjectIDFromStringID(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return result, fmt.Errorf("error al eliminar el workout en WorkoutRepository.DeleteWorkout(): %v", err)
	}
	return result, nil
}

func (repository WorkoutRepository) GetWorkoutsByUserID(userID string) ([]models.Workout, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("workouts")
	userObjectID, err := utils.GetObjectIDFromStringID(userID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"user_id": userObjectID}

	result, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta Find() en WorkoutRepository.GetWorkoutsByUserID(): %v", err)
	}
	defer result.Close(context.TODO())

	var workouts []models.Workout
	for result.Next(context.Background()) {
		var workout models.Workout
		err := result.Decode(&workout)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar el workout en WorkoutRepository.GetWorkoutsByUserID(): %v", err)
		}
		workouts = append(workouts, workout)
	}
	return workouts, nil
}
