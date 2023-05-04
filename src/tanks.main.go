package main

import (
	"net/http"
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
	println("You server run 8080")
	http.ListenAndServe(":8080", nil)
}
