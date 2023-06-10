package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func save(object interface{}, client *mongo.Client, collectionName string) {
	var collection = client.Database("game_stats").Collection(collectionName)
	_, err := collection.InsertOne(context.TODO(), object)
	if err != nil {
		fmt.Printf("Could not save data in collectvvvion %s: %s", collection, err)
	}
}

func initializeMongoSession() (*mongo.Client, error) {
	return mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
}
