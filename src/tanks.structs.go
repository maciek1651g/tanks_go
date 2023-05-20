package main

import (
	"encoding/json"
)

type Payload struct {
	MessageType string `json:"messageType"`
	Data        string `json:"data"`
}

type UserConnectedPayload struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
}

type UserStatusPayload struct {
	Id          string      `json:"id"`
	MessageType string      `json:"messageType"`
	Health      int         `json:"health"`
	Coordinates Coordinates `json:"coordinates"`
}

type UserDisconnectedPayload struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
}

type UserAttackPayload struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
}

type CreateChestPayload struct {
	Id          string      `json:"id"`
	MessageType string      `json:"messageType"`
	Coordinates Coordinates `json:"coordinates"`
}

type ChestGrabPayload struct {
	Id string `json:"id"`
}

type ChestDestroyedPayload struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
}

type Coordinates struct {
	X          int `json:"x"`
	Y          int `json:"y"`
	DirectionX int `json:"directionX"`
}

func createUserConnectedPayload(message []byte) (UserConnectedPayload, error) {
	var requestPayload UserConnectedPayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)
	return requestPayload, unmarshallErr
}

func createUserStatusPayload(message []byte) (UserStatusPayload, error) {
	var requestPayload UserStatusPayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)
	return requestPayload, unmarshallErr
}

func createUserAttackPayload(message []byte) (UserAttackPayload, error) {
	var requestPayload UserAttackPayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)
	return requestPayload, unmarshallErr
}

func createChestGrabPayload(message []byte) (ChestGrabPayload, error) {
	var requestPayload ChestGrabPayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)
	return requestPayload, unmarshallErr
}

func createChestCreatePayload(chest Chest) CreateChestPayload {
	return CreateChestPayload{Id: chest.Id, MessageType: "create_chest", Coordinates: chest.Coordinates}
}

func createChestDestroyedPayload(id string) ChestDestroyedPayload {
	return ChestDestroyedPayload{Id: id, MessageType: "chest_destroy"}
}

type Chest struct {
	Id          string
	Coordinates Coordinates
	Destroyed   bool
}
