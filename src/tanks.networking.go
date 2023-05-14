package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var connector = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients []websocket.Conn

var objects = make(map[string]TanksPayload)

func handleTanksConnection(w http.ResponseWriter, r *http.Request) {

	conn, err := connector.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("Error occurred : %s", err)
	}

	fmt.Println("Initialized connection for : " + conn.RemoteAddr().String())

	clients = append(clients, *conn)

	sendHistoricalPayloadsForObjects(conn)

	for true {
		_, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			return
		}

		fmt.Printf("%s : Received payload = '%s'\n", conn.RemoteAddr(), string(message))
		var payload, payloadError = createTanksPayload(message)

		if payloadError != nil {
			fmt.Printf(payloadError.Error() + "\n")
			return
		}

		savePayload(payload.Id, payload)

		for _, client := range clients {
			if client.RemoteAddr() != conn.RemoteAddr() {
				fmt.Printf("%s : Sending payload : '%s'\n", client.RemoteAddr(), payload)
				if err = client.WriteJSON(payload); err != nil {
					fmt.Printf(err.Error() + "\n")
					return
				}
			}
		}
	}
}

func createTanksPayload(message []byte) (TanksPayload, error) {
	var requestPayload TanksPayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)
	return requestPayload, unmarshallErr
}

func savePayload(id string, payload TanksPayload) {
	objects[id] = payload
}

func sendHistoricalPayloadsForObjects(client *websocket.Conn) {
	for _, value := range objects {
		fmt.Printf("Sending historical data to %s : %s\n", client.RemoteAddr(), value)
		var err = client.WriteJSON(value)

		if err != nil {
			fmt.Printf("There was an error when sending payload to %s : %s\n", client.RemoteAddr(), err.Error())
		}
	}
}
