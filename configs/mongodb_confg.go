package configs

import (
	"context"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fiberlog.Info("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("testing").Collection(collectionName)
	return collection
}
