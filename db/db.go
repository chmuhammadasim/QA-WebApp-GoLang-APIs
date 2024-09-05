package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Connect to MongoDB
func Connect() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB ping failed: ", err)
	}

	log.Println("Connected to MongoDB")
}

// GetCollection returns a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	if client == nil {
		log.Fatal("MongoDB client not initialized. Make sure Connect() is called before using GetCollection.")
	}
	return client.Database("qa-app").Collection(collectionName)
}
