package repositories

import (
	"AppFitness/models"
	"AppFitness/utils"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoutineRepositoryInterface interface {
	PostRoutine(models.Routine) (*mongo.InsertOneResult, error)
	GetRoutines() ([]*models.Routine, error)
	GetRoutineByID(id string) (*models.Routine, error)
	PutRoutine(routine models.Routine) (*mongo.UpdateResult, error)
	DeleteRoutine(id string) (*mongo.DeleteResult, error)
	AddExerciseRutine(exercise models.ExcerciseInRoutine, idRutine primitive.ObjectID) (*mongo.UpdateResult, error)
	UpdateExerciseInRoutine(idRutine primitive.ObjectID, idExercise primitive.ObjectID, exerciseMod models.ExcerciseInRoutine) (*mongo.UpdateResult, error)
	DeleteExerciseToRutine(rutineID primitive.ObjectID, exerciseID primitive.ObjectID) (*mongo.UpdateResult, error)
	ExistByRutineName(rutineName string) (bool, error)
}

type RoutineRepository struct {
	db DB
}

func NewRoutineRepository(db DB) *RoutineRepository {
	return &RoutineRepository{
		db: db,
	}
}

func (repository RoutineRepository) PostRoutine(routine models.Routine) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	result, err := collection.InsertOne(nil, routine)
	if err != nil {
		return result, fmt.Errorf("error al insertar la rutina en RoutineRepository.PostRoutine(): %v", err)
	}
	return result, nil
}

func (repository RoutineRepository) GetRoutines() ([]*models.Routine, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error en Find() RoutineRepository.GetRoutines(): %v", err)
	}
	defer cursor.Close(context.TODO())

	var routines []*models.Routine
	for cursor.Next(context.Background()) {
		var routine *models.Routine
		err := cursor.Decode(&routine)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar la rutina en RoutineRepository.GetRoutines(): %v", err)
		}
		routines = append(routines, routine)
	}
	return routines, err
}

func (repository RoutineRepository) GetRoutineByID(id string) (*models.Routine, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	objID, err := utils.GetObjectIDFromStringID(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	result := collection.FindOne(context.TODO(), filter)

	var routine *models.Routine
	err = result.Decode(&routine)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina en RoutineRepository.GetRoutineByID(): %v", err)
	}
	return routine, nil
}

func (repository RoutineRepository) PutRoutine(routine models.Routine) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	filter := bson.M{"_id": routine.ID}
	entity := bson.M{"$set": routine}
	result, err := collection.UpdateOne(context.TODO(), filter, entity)
	if err != nil {
		return result, fmt.Errorf("error al actualizar la rutina en RoutineRepository.PutRoutine(): %v", err)
	}
	return result, nil
}

func (repository RoutineRepository) DeleteRoutine(id string) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	objectID, err := utils.GetObjectIDFromStringID(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return result, fmt.Errorf("error al eliminar la rutina en RoutineRepository.DeleteRoutine(): %v", err)
	}
	return result, nil
}

func (repository RoutineRepository) AddExerciseRutine(exercise models.ExcerciseInRoutine, idRutine primitive.ObjectID) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")

	update := bson.M{
		"$push": bson.M{
			"exercise_list": exercise,
		},
		"$set": bson.M{
			"edition_date": time.Now(),
		},
	}

	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": idRutine},
		update,
	)

	if err != nil {
		return nil, fmt.Errorf("error al agregar ejercicio a la rutina: %v", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no se encontró ninguna rutina con el ID especificado")
	}

	return result, nil
}

func (repository RoutineRepository) UpdateExerciseInRoutine(idRutine primitive.ObjectID, idExercise primitive.ObjectID, exerciseMod models.ExcerciseInRoutine) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")

	set := bson.M{}
	if exerciseMod.Series > 0 {
		set["exercise_list.$[e].series"] = exerciseMod.Series
	}
	if exerciseMod.Repetitions > 0 {
		set["exercise_list.$[e].repetitions"] = exerciseMod.Repetitions
	}
	if exerciseMod.Weight >= 0 {
		set["exercise_list.$[e].weight"] = exerciseMod.Weight
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("no se enviaron campos para actualizar")
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{ //Esa línea crea las opciones que le dicen a MongoDB qué elemento del array debe modificar, en lugar de tocar todos.
		Filters: []interface{}{bson.M{"e.excercise_id": idExercise}},
	})

	res, err := collection.UpdateOne( //y aca buscamos la rutina con su id y seteamos la linea de opciones anteriormente obtenidas
		context.TODO(),
		bson.M{"_id": idRutine},
		bson.M{"$set": set},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar ejercicio en rutina: %v", err)
	}
	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("no se encontró la rutina")
	}
	if res.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se encontró el ejercicio dentro de la rutina")
	}
	return res, nil
}

func (repository RoutineRepository) DeleteExerciseToRutine(rutineID primitive.ObjectID, exerciseID primitive.ObjectID) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")

	update := bson.M{
		"$pull": bson.M{
			"exercise_list": bson.M{
				"excercise_id": exerciseID,
			},
		},
	}
	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": rutineID},
		update,
	)
	if err != nil {
		return nil, fmt.Errorf("error al eliminar ejercicio de la rutina: %v", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no se encontró la rutina especificada")
	}

	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se encontró el ejercicio dentro de la rutina")
	}

	return result, nil
}

func (r RoutineRepository) ExistByRutineName(rutineName string) (bool, error) { //verificamos si existe algun usuario con el mismo nombre de ussuario y lo llamamos de UserService
	collection := r.db.GetClient().Database("AppFitness").Collection("routines")
	filter := bson.M{"name": rutineName}

	count, err := collection.CountDocuments(context.TODO(), filter)

	if err != nil {
		return false, fmt.Errorf("error al verificar nombre de usuario: %w", err)
	}

	return count > 0, err
}
