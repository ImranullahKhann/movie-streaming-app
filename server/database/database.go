package database

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os"
)

var dbClient *mongo.Client

func ConnectDB() {
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatal("environment variables not found")
	}
	fmt.Println(mongodbURI)

	client, err := mongo.Connect(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal(err)
	}

	// Checking connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	dbClient = client
}

func OpenCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	collection := dbClient.Database(dbName).Collection(collectionName)

	return collection
}
