package repositories

import (
	"AppFitness/dto"
	"AppFitness/models"
	"AppFitness/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface { //contrato que define metodos para manejar usuarios
	GetUsers() ([]models.User, error)
	GetUsersByID(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	PostUser(user models.User) (*mongo.InsertOneResult, error)
	PutUser(user models.User) (*mongo.UpdateResult, error)
	UpdateNewPassword(dto dto.PasswordChange, id string) (modified int64, err error)
	DeleteUser(id string) (*mongo.DeleteResult, error)
	ExistByEmail(email string) (bool, error)
	ExistByUserName(userName string) (bool, error)
	ExistByUserNameExceptID(id string, userName string) (bool, error)
}

type UserRepository struct { //campo para la conexion a la base de datos
	db DB
}

func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository UserRepository) PostUser(user models.User) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return result, fmt.Errorf("error al insertar el usuario en UserRepository.PostUser(): %v", err)
	}
	return result, err
}

func (repository UserRepository) GetUsers() ([]models.User, error) { //REVISAR Y AJUSTAR LA DEVOLUCION DE ERRORES
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	filter := bson.M{} //filtro vacio para traer todos los documentos

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta Find(): %w", err)
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar el usuario en UserRepository.GetUsers(): %v", err)
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar el cursor: %w", err)
	}

	return users, err
}

func (repository UserRepository) GetUsersByID(id string) (models.User, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	objectID, err := utils.GetObjectIDFromStringID(id)
	if err != nil {
		return models.User{}, fmt.Errorf("ID de formato inválido")
	}

	filter := bson.M{"_id": objectID}

	var user models.User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, fmt.Errorf("no se encontró ningún usuario con el ID proporcionado")
		}
		return models.User{}, fmt.Errorf("error al obtener usuario por ID: %w", err)
	}

	return user, err
}

func (repository UserRepository) GetUserByEmail(email string) (models.User, error) { //para recuperar usuario por email en el login y recuperar ids en service
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	filter := bson.M{"email": email}

	result := collection.FindOne(context.TODO(), filter)

	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("error al obtener el usuario en UserRepository.GetUserByEmail(): %v", err)
	}
	return user, nil
}

func (repository UserRepository) PutUser(user models.User) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	filter := bson.M{"_id": user.ID}

	entity := bson.M{"$set": bson.M{
		"user_name":  user.UserName,
		"email":      user.Email,
		"role":       user.Role,
		"weight":     user.Weight,
		"height":     user.Height,
		"experience": user.Experience,
		"objetive":   user.Objetive,
	}}
	result, err := collection.UpdateOne(context.TODO(), filter, entity)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repository UserRepository) DeleteUser(id string) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	objectID, err := utils.GetObjectIDFromStringID(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return result, fmt.Errorf("error al eliminar el usuario en UserRepository.DeleteUser(): %v", err)
	}

	return result, err
}

func (r UserRepository) ExistByEmail(email string) (bool, error) { //verificamos si existe algun usuario con el mismo email y lo llamamos de UserService
	collection := r.db.GetClient().Database("AppFitness").Collection("users")
	filter := bson.M{"email": email}

	count, err := collection.CountDocuments(context.TODO(), filter)

	if err != nil {
		return false, fmt.Errorf("error al verificar email: %w", err)
	}

	return count > 0, err
}

func (r UserRepository) ExistByUserName(userName string) (bool, error) { //verificamos si existe algun usuario con el mismo nombre de ussuario y lo llamamos de UserService
	collection := r.db.GetClient().Database("AppFitness").Collection("users")
	filter := bson.M{"user_name": userName}

	count, err := collection.CountDocuments(context.TODO(), filter)

	if err != nil {
		return false, fmt.Errorf("error al verificar nombre de usuario: %w", err)
	}

	return count > 0, err
}

func (r UserRepository) ExistByUserNameExceptID(id string, userName string) (bool, error) {
	collection := r.db.GetClient().Database("AppFitness").Collection("users")
	objectID, err := utils.GetObjectIDFromStringID(id)
	if err != nil {
		return false, err
	}
	filter := bson.M{
		"user_name": userName,
		"_id":       bson.M{"$ne": objectID},
	} //con este filtro buscamos si existe algun email pero q no se el de nuestro email

	count, err := collection.CountDocuments(context.TODO(), filter)

	if err != nil {
		return false, fmt.Errorf("error al verificar el nombre de usuario excepto el de nuestro id: %w", err)
	}

	return count > 0, err
}

func (r *UserRepository) UpdateNewPassword(dto dto.PasswordChange, id string) (modified int64, err error) {
	collection := r.db.GetClient().Database("AppFitness").Collection("users")

	ID, err := primitive.ObjectIDFromHex(strings.TrimSpace(id))
	if err != nil {
		return 0, fmt.Errorf("ID inválido")
	}

	filter := bson.M{"_id": ID}

	update := bson.M{"$set": bson.M{"password": dto.NewPassword}}

	res, err := collection.UpdateOne(
		context.TODO(),
		filter, // Usamos el filtro simplificado
		update,
	)
	if err != nil {
		return 0, err
	}
	//devolvemos la cantidad  de documentos q fueron modificados
	return res.ModifiedCount, nil
}
