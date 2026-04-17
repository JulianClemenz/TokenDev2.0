package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
}

func NewMongoDB() *MongoDB {
	instancia := &MongoDB{}
	instancia.Connect()

	return instancia
}

func (mongoDB *MongoDB) GetClient() *mongo.Client {
	return mongoDB.Client
}

// La dejamos privada, se ejecuta cuando se crea el objeto
func (mongoDB *MongoDB) Connect() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	mongoDB.Client = client

	return nil
}

func (mongoDB *MongoDB) Disconnect() error {
	return mongoDB.Client.Disconnect(context.Background())
}
