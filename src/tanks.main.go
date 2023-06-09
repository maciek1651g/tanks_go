package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	var session, err = initializeMongoSession()
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	save(UserCreatePayload{Id: "", Health: 100.00, Coordinates: Coordinates{X: 200, Y: 600}}, session, "users")
	initializeEngine()
	configureConnector()
	initializeWebSockets()
	initializeHttpEndpoints()
	initializeServer()
}

func configureConnector() {
	connector.CheckOrigin = func(r *http.Request) bool { return true }
}

func initializeWebSockets() {
	http.HandleFunc("/tanks/objects:exchange", handleTanksConnection)
}

func initializeHttpEndpoints() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
}

func initializeServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	println("You server run " + port)
	http.ListenAndServe(":"+port, nil)
}
