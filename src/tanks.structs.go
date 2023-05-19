package main

import (
	"encoding/json"
)

type UserConnectedPayload struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
}

type CoordinatesChangedPayload struct {
	Id          string      `json:"id"`
	MessageType string      `json:"messageType"`
	Coordinates Coordinates `json:"coordinates"`
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

type Coordinates struct {
	X          int `json:"x"`
	Y          int `json:"y"`
	DirectionX int `json:"directionX"`
}

func createCoordinatesPayload(message []byte) (CoordinatesChangedPayload, error) {
	var requestPayload CoordinatesChangedPayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)
	return requestPayload, unmarshallErr
}

func createUserConnectedPayload(message []byte) (UserConnectedPayload, error) {
	var requestPayload UserConnectedPayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)
	return requestPayload, unmarshallErr
}

func createStatusPayload(id string, health int, coordinates Coordinates) UserStatusPayload {
	return UserStatusPayload{Id: id, MessageType: "status", Health: health, Coordinates: coordinates}
}
