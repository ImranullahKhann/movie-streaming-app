package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
)

func ConnectDB() (*mongo.Client, error) {
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		return nil, fmt.Errorf("environment variables not found: MONGODB_URI")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Checking connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	collection := client.Database(dbName).Collection(collectionName)

	return collection
}
