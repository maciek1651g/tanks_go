package main

import (
	"net/http"
	"os"
)

func main() {
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
