package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var ProductsCollection *mongo.Collection
var USersCollection *mongo.Collection

func ConnectToMongo() (*mongo.Client, *mongo.Collection, error) {
	// MongoDB connection URI
	mongoURI := "mongodb://localhost:27017/"

	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, err
	}
	// Initialize the collections
	ProductsCollection = client.Database("prodctCollection").Collection("products")
	// usersCollection = client.Database("your-database-name").Collection("users")
	return client, ProductsCollection, nil

}
