package repositories

import (
	"context"
	"os" // ✅ CORRECCIÓN: Para leer MONGO_URI

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB() *MongoDB {
	instancia := &MongoDB{}
	instancia.Connect()
	return instancia
}

func (mongoDB *MongoDB) GetClient() *mongo.Client {
	return mongoDB.Client
}

func (mongoDB *MongoDB) Connect() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // Fallback por seguridad
	}

	clientOptions := options.Client().ApplyURI(uri)
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
