package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var connector = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients []websocket.Conn
var metadata = sync.Map{}
var objects = sync.Map{}

func handleTanksConnection(w http.ResponseWriter, r *http.Request) {

	conn, err := connector.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("Error occurred : %s", err)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		handleUserDisconnection(conn)
		return nil
	})

	fmt.Println("Initialized connection for : " + conn.RemoteAddr().String())

	clients = append(clients, *conn)
	performUserInitialization(conn)
	sendHistoricalPayloadsForObjects(conn)

	for true {
		_, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			return
		}

		fmt.Printf("%s : Received payload = '%s'\n", conn.RemoteAddr(), string(message))
		var payload, payloadError = createCoordinatesPayload(message)

		if payloadError != nil {
			fmt.Printf(payloadError.Error() + "\n")
			return
		}

		saveCoordinates(payload.Id, payload)
		broadcastPayload(conn, payload)
	}
}

func saveCoordinates(id string, payload CoordinatesChangedPayload) {
	objects.Store(id, payload)
}

func performUserInitialization(client *websocket.Conn) {

	_, message, err := client.ReadMessage()

	if err != nil {
		fmt.Printf("Error when reading initialization payload : " + err.Error() + "\n")
		return
	}

	fmt.Printf("%s : Received initialization informations = '%s'\n", client.RemoteAddr(), string(message))
	var payload, payloadError = createUserConnectedPayload(message)

	if payloadError != nil {
		fmt.Printf("Error when convertin initialization payload to Object : " + payloadError.Error() + "\n")
		return
	}

	metadata.Store(client.RemoteAddr(), payload.Id)
	var initializationPayload = InitializePayload{Id: payload.Id, MessageType: "create_player", Health: 100, Coordinates: Coordinates{X: 200, Y: 600}}
	if err = client.WriteJSON(initializationPayload); err != nil {
		fmt.Printf("Error when sending 'create_player' payload : " + err.Error() + "\n")
	}
	broadcastPayload(client, initializationPayload)
}

func sendHistoricalPayloadsForObjects(client *websocket.Conn) {
	objects.Range(func(key, value interface{}) bool {
		fmt.Printf("Sending historical data to %s : %s\n", client.RemoteAddr(), value)
		var err = client.WriteJSON(value)

		if err != nil {
			fmt.Printf("There was an error when sending payload to %s : %s\n", client.RemoteAddr(), err.Error())
		}

		return true
	})
}

func handleUserDisconnection(client *websocket.Conn) {
	metadata.Range(func(address, id interface{}) bool {
		if address == client.RemoteAddr() {
			sendUserDisconnection(client, id.(string))
			metadata.Delete(address)
		}
		return true
	})
}

func sendUserDisconnection(client *websocket.Conn, id string) {
	var payload = UserDisconnectedPayload{Id: id, MessageType: "user_disconnected"}
	fmt.Printf("Broadcasting 'user_disconnected' payload for id = %s\n", id)
	broadcastPayload(client, payload)
}

func broadcastPayload(currentConnection *websocket.Conn, payload any) {
	for _, client := range clients {
		if client.RemoteAddr() != currentConnection.RemoteAddr() {
			fmt.Printf("%s : Sending payload : '%s'\n", client.RemoteAddr(), payload)
			if err := client.WriteJSON(payload); err != nil {
				fmt.Printf(err.Error() + "\n")
				return
			}
		}
	}
}
