package main

import (
	"net/http"
	"os"
)

func main() {
	configureConnector()
	initializeWebSockets()
	initializeServer()
}

func configureConnector() {
	connector.CheckOrigin = func(r *http.Request) bool { return true }
}

func initializeWebSockets() {
	http.HandleFunc("/tanks/objects:exchange", handleTanksConnection)
}

func initializeServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	println("You server run " + port)
	http.ListenAndServe(":"+port, nil)
}
