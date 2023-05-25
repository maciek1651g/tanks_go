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
	Health      float64     `json:"health"`
	Coordinates Coordinates `json:"coordinates"`
}

type UserCreatePayload struct {
	Id          string      `json:"id"`
	MessageType string      `json:"messageType"`
	Health      float64     `json:"health"`
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

type UserDamagePayload struct {
	Id          string  `json:"id"`
	TargetId    string  `json:"targetId"`
	Damage      float64 `json:"damage"`
	MessageType string  `json:"messageType"`
}

type CreateChestPayload struct {
	Id          string      `json:"id"`
	MessageType string      `json:"messageType"`
	Coordinates Coordinates `json:"coordinates"`
}

type MobDestroyedPayload struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
}

type ChestGrabPayload struct {
	Id string `json:"id"`
}

type ChestDestroyedPayload struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
}

type MobCreatedPayload struct {
	Id          string      `json:"id"`
	MessageType string      `json:"messageType"`
	Health      float64     `json:"health"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	X          int `json:"x"`
	Y          int `json:"y"`
	DirectionX int `json:"directionX"`
}

type GameMasterPayload struct {
	Id          string `json:"id"`
	MessageType string `json:"messageType"`
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

func createUserDamagePayload(message []byte) (UserDamagePayload, error) {
	var requestPayload UserDamagePayload
	unmarshallErr := json.Unmarshal(message, &requestPayload)
	return requestPayload, unmarshallErr
}

func createChestCreatePayload(chest Chest) CreateChestPayload {
	return CreateChestPayload{Id: chest.Id, MessageType: "create_chest", Coordinates: chest.Coordinates}
}

func createChestDestroyedPayload(id string) ChestDestroyedPayload {
	return ChestDestroyedPayload{Id: id, MessageType: "chest_destroy"}
}

func createMobCreatedPayload(mob Mob) MobCreatedPayload {
	return MobCreatedPayload{Id: mob.Id, MessageType: "mob_create", Coordinates: mob.Coordinates, Health: mob.Health}
}

func createMobDestroyedPayload(id string) MobDestroyedPayload {
	return MobDestroyedPayload{Id: id, MessageType: "mob_destroy"}
}

func createStandardMob(id string, coordinates Coordinates) Mob {
	return Mob{Id: id, Coordinates: coordinates, Destroyed: false, Health: 100.000000}
}

func createGameMasterPayload(id string) GameMasterPayload {
	return GameMasterPayload{Id: id, MessageType: "game_master"}
}

func createUserStatusPayloadFromPlayer(player Player) UserStatusPayload {
	return UserStatusPayload{Id: player.Id, MessageType: "status", Coordinates: player.Coordinates, Health: player.Health}
}

func createPlayerCreatePayload(player Player) UserCreatePayload {
	return UserCreatePayload{Id: player.Id, MessageType: "create_player", Coordinates: player.Coordinates, Health: player.Health}
}

type Chest struct {
	Id          string
	Coordinates Coordinates
	Destroyed   bool
}

type Mob struct {
	Id          string
	Coordinates Coordinates
	Destroyed   bool
	Health      float64
}

type Player struct {
	Id          string
	Coordinates Coordinates
	Destroyed   bool
	Health      float64
	Master      bool
}

func (mob *Mob) dealDamage(damage float64) {
	mob.Health = mob.Health - damage
	mob.Destroyed = mob.Health <= 0
}
