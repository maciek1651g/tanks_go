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
	sendCurrentUserStatuses(conn)
	sendCurrentChests(conn)

	for true {
		_, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			return
		}

		fmt.Printf("%s : Received payload = '%s'\n", conn.RemoteAddr(), string(message))
		handleActionPayload(conn, message)
	}
}

func saveUserStatus(id string, payload UserStatusPayload) {
	objects.Store(id, payload)
}

func performUserInitialization(client *websocket.Conn) {

	_, message, err := client.ReadMessage()

	fmt.Printf("message: %s \n", string(message))

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

	fmt.Printf("payload: %s \n", payload)
	metadata.Store(client.RemoteAddr(), payload.Id)
	var initializationPayload = UserStatusPayload{Id: payload.Id, MessageType: "create_player", Health: 100, Coordinates: Coordinates{X: 200, Y: 600}}
	saveUserStatus(payload.Id, initializationPayload)
	broadcastPayload(client, initializationPayload)
}

func sendCurrentUserStatuses(client *websocket.Conn) {
	var id, _ = metadata.Load(client.RemoteAddr())
	objects.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %s \n value: %s \n", key, value)
		if id != key {
			fmt.Printf("Sending historical data to %s : %s\n", client.RemoteAddr(), value)
			var err = client.WriteJSON(value)

			if err != nil {
				fmt.Printf("There was an error when sending payload to %s : %s\n", client.RemoteAddr(), err.Error())
			}
		}

		return true
	})
}

func handleUserDisconnection(client *websocket.Conn) {
	metadata.Range(func(address, id interface{}) bool {
		if address == client.RemoteAddr() {
			sendUserDisconnection(client, id.(string))
			var id, _ = metadata.Load(address)
			metadata.Delete(address)
			objects.Delete(id)
		}
		return true
	})
}

func sendCurrentChests(client *websocket.Conn) {
	for _, chest := range chests {
		if chest.Destroyed == false {
			var payload = createChestCreatePayload(chest)
			if err := client.WriteJSON(payload); err != nil {
				fmt.Printf("Error occurred when sending 'CreateChestPayload' %s\n", payload)
			}
		}
	}
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
