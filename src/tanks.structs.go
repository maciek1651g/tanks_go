package main

type TanksPayload struct {
	Id          string      `json:"id"`
	MessageType string      `json:"messageType"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}
