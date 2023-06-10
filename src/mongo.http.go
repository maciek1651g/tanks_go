package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func runMongoCommand(w http.ResponseWriter, r *http.Request) {
	var session, initError = initializeMongoSession()
	var query = r.URL.Query()

	//var users = session.Database("game_stats").Collection("users")

	if initError != nil {
		fmt.Printf("Error when initialization : %s", initError)
		w.WriteHeader(400)
	}

	cmdString := query.Get("query")
	var cmd bson.D
	var err = bson.UnmarshalExtJSON([]byte(cmdString), true, &cmd)

	if err != nil {
		fmt.Printf("Error when unmarshalling command : %s", err)
		w.WriteHeader(400)
		return
	}

	var cursor, error = session.Database("game_stats").RunCommandCursor(context.Background(), cmd)

	if error != nil {
		fmt.Printf("Error when returning response : %s", error)
		w.WriteHeader(400)
		return
	}

	var result bson.D
	var results []json.RawMessage
	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			fmt.Printf("Error when decoding : %s\n", err)
		}

		json, err := bson.MarshalExtJSON(result, false, true)

		if err != nil {
			fmt.Printf("Error when decoding : %s\n", err)
		}

		results = append(results, json)
	}

	var response, responseErr = json.Marshal(results)

	if responseErr != nil {
		fmt.Printf("Error when returning response : %s", responseErr)
		w.WriteHeader(400)
		return
	}

	w.Write(response)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	//response, err := users.Find(context.TODO(), nil)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer response.Close(context.TODO())
	//
	//// Iterate over the cursor and process each document
	//
	//var players []Player
	//for response.Next(context.TODO()) {
	//	var player Player
	//	err := response.Decode(&player)
	//	if err != nil {
	//		fmt.Printf("Problem occurred when decoding player : %s\n", err)
	//	}
	//
	//	players = append(players, player)
	//}
	//
	//var userResponse, unmarshallErr = json.Marshal(players)
	//
	//if unmarshallErr != nil {
	//	fmt.Printf("Error when retrieving users : %s\n", unmarshallErr)
	//}
	//
	//w.Write(userResponse)
	//w.WriteHeader(200)
}
