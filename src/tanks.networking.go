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

func handleTanksConnection(w http.ResponseWriter, r *http.Request) {

	conn, err := connector.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("Error occurred : %s", err)
	}

	fmt.Println("Initialized connection for : " + conn.RemoteAddr().String())

	clients = append(clients, *conn)

	for {
		msgType, message, err := conn.ReadMessage()

		if err != nil {
			return
		}

		fmt.Printf("%s : Received payload = '%s'\n", conn.RemoteAddr(), string(message))
		var payload, payloadError = createTanksPayload(message)

		if payloadError != nil {
			return
		}

		for _, client := range clients {
			if client.RemoteAddr() != conn.RemoteAddr() {
				fmt.Printf("%s : Sending payload : '%s'\n", client.RemoteAddr(), string(payload))
				if err = client.WriteMessage(msgType, payload); err != nil {
					return
				}
			}
		}
	}
}

func createTanksPayload(message []byte) ([]byte, error) {
	var requestPayload TanksPayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)

	if unmarshallErr != nil {
		return nil, unmarshallErr
	}

	return json.Marshal(requestPayload)
}
