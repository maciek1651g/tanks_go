package main

import (
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

var corsHandler = handlers.CORS(
	handlers.AllowedOrigins([]string{"*"}),
	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	handlers.AllowedHeaders([]string{"Content-Type"}),
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
	http.Handle("/tanks/objects:exchange", corsHandler(http.HandlerFunc(handleTanksConnection)))
}

func initializeHttpEndpoints() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/api/users:migrate", corsHandler(http.HandlerFunc(migrateUser)))
	http.Handle("/api/users:migrate-add", corsHandler(http.HandlerFunc(migrateAddUser)))
}

func initializeServer() {
	args := os.Args
	var port string

	if len(args) > 1 {
		port = args[1]
	} else {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	println("You server run " + port)
	http.ListenAndServe(":"+port, nil)
}
