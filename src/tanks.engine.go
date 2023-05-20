package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var chests []Chest

func initializeEngine() {
	createChests()
}

func createChests() {
	chests = []Chest{
		{
			Id:          uuid.New().String(),
			Coordinates: Coordinates{X: 200, Y: 600},
			Destroyed:   false,
		},
		{
			Id:          uuid.New().String(),
			Coordinates: Coordinates{X: 300, Y: 1000},
			Destroyed:   false,
		},
		{
			Id:          uuid.New().String(),
			Coordinates: Coordinates{X: 500, Y: 750},
			Destroyed:   false,
		},
		{
			Id:          uuid.New().String(),
			Coordinates: Coordinates{X: 1100, Y: 700},
			Destroyed:   false,
		},
		{
			Id:          uuid.New().String(),
			Coordinates: Coordinates{X: 950, Y: 730},
			Destroyed:   false,
		},
	}
}

func handleActionPayload(conn *websocket.Conn, bytes []byte) {
	var payload Payload
	var error = json.Unmarshal(bytes, &payload)

	if error != nil {
		fmt.Printf("Could not deserialize payload %s", error)
		return
	}

	switch payload.MessageType {
	case "status":
		handleStatusPayload(conn, payload)
	case "user_attack":
		handleUserAttackPayload(conn, payload)
	case "chest_grab":
		handleChestGrabPayload(conn, payload)
	}
}

func handleStatusPayload(conn *websocket.Conn, payload Payload) {
	var userStatusPayload, err = createUserStatusPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'UserStatusPayload' %s", err)
		return
	}

	updateCoordinates(conn, userStatusPayload.Coordinates)
	broadcastPayload(conn, userStatusPayload)
}

func handleUserAttackPayload(conn *websocket.Conn, payload Payload) {
	var userAttackPayload, err = createUserAttackPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'UserAttackPayload' %s", err)
		return
	}

	broadcastPayload(conn, userAttackPayload)
}

func updateCoordinates(client *websocket.Conn, coordinates Coordinates) {
	var id, _ = metadata.Load(client.RemoteAddr())
	var status, _ = objects.Load(id)
	var currentStatus = status.(UserStatusPayload)
	currentStatus.Coordinates = coordinates
	objects.Store(id, currentStatus)
}

func handleChestGrabPayload(conn *websocket.Conn, payload Payload) {
	var chestGrabPayload, err = createChestGrabPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'ChestGrabPayload' %s", err)
		return
	}

	if deleteChest(chestGrabPayload.Id) == true {
		broadcastPayload(conn, createChestDestroyedPayload(chestGrabPayload.Id))
	}
}

func deleteChest(id string) bool {
	for index, chest := range chests {
		if chest.Id == id {
			chests[index].Destroyed = true
			return true
		}
	}

	return false
}
