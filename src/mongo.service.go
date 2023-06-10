package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
)

var mongoClients = sync.Map{}

func saveInMongo(object interface{}, client *mongo.Client, collectionName string) {
	var collection = client.Database("game_stats").Collection(collectionName)
	_, err := collection.InsertOne(context.TODO(), object)
	if err != nil {
		fmt.Printf("Could not saveInMongo data in collectvvvion %s: %s", collection, err)
	}
}

func initializeMongoSession() (*mongo.Client, error) {
	return mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
}

func storeSession(client *mongo.Client, conn *websocket.Conn) {
	mongoClients.Store(conn.RemoteAddr(), client)
}
