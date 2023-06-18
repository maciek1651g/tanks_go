package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

var chests = sync.Map{}
var mobs = sync.Map{}
var players = sync.Map{}

func initializeEngine() {
	createMobs()
	createChests()
}

func createMobs() {

	var coordinates = []Coordinates{
		{X: 200, Y: 200},
		{X: 300, Y: 300},
		{X: 400, Y: 400},
		{X: 500, Y: 500},
		{X: 600, Y: 600},
		{X: 700, Y: 700},
		{X: 800, Y: 800},
		{X: 900, Y: 900},
		{X: 1000, Y: 1000},
		{X: 1100, Y: 1100},
		{X: 550, Y: 600},
	}

	for _, coordinate := range coordinates {
		var id = uuid.New().String()
		chests.Store(id, createStandardMob(id, coordinate))
	}
}

func createChests() {
	var coordinates = []Coordinates{
		{X: 200, Y: 600},
		{X: 300, Y: 1000},
		{X: 500, Y: 750},
		{X: 1100, Y: 700},
		{X: 950, Y: 730},
	}

	for _, coordinate := range coordinates {
		var id = uuid.New().String()
		chests.Store(id, Chest{Id: id, Coordinates: coordinate, Destroyed: false})
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
	case "user_damage":
		handleUserDamagePayload(conn, payload)
	case "mob_status":
		handleMobStatusPayload(conn, payload)
	case "mob_damage":
		handleMobDamagePayload(conn, payload)
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

func handleMobStatusPayload(conn *websocket.Conn, payload Payload) {
	var mobStatusPayload, err = createMobStatusPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'MobStatusPayload' %s", err)
		return
	}

	if updateMobCoordinates(mobStatusPayload.Id, mobStatusPayload.Coordinates) {
		broadcastPayload(conn, mobStatusPayload)
	}
}

func updateMobCoordinates(id string, coordinates Coordinates) bool {
	var updated = false

	mobs.Range(func(key, value any) bool {
		var mob = value.(Mob)
		if mob.Id == id {
			mob.Coordinates = coordinates
			mobs.Store(id, mob)
			updated = true
			return false
		} else {
			return true
		}
	})

	return updated
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
	var status, _ = players.Load(id)
	var currentStatus = status.(Player)
	currentStatus.Coordinates = coordinates
	players.Store(id, currentStatus)
}

func handleChestGrabPayload(conn *websocket.Conn, payload Payload) {
	var chestGrabPayload, err = createChestGrabPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'ChestGrabPayload' %s", err)
		return
	}

	if deleteChest(chestGrabPayload.Id) == true {
		resolveScore(conn, chestGrabPayload.PlayerId, 10)
		broadcastPayload(conn, createChestDestroyedPayload(chestGrabPayload.Id))
	}
}

func resolveScore(client *websocket.Conn, id string, score float64) {
	players.Range(func(key, value any) bool {
		var player = value.(Player)
		if player.Id == id {
			var enrichedPlayer = assignScore(player, score)
			var payload = createUserScorePayload(enrichedPlayer)
			if err := client.WriteJSON(payload); err != nil {
				fmt.Printf("Error occurred when sending 'UserScorePayload' : %s\n", payload)
			}
			return false
		}
		return true
	})
}

func assignScore(player Player, amount float64) Player {
	player.Score = player.Score + amount
	players.Store(player.Id, player)
	return player
}

func handleUserDamagePayload(conn *websocket.Conn, payload Payload) {
	var userDamagePayload, err = createUserDamagePayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'UserDamagePayload' %s", err)
		return
	}

	var damaged, destroyed = dealDamageToUser(userDamagePayload)

	if damaged {
		var player, _ = players.Load(userDamagePayload.TargetId)
		broadcastPayloadToAll(createUserHealthPayload(player.(Player)))
	}

	if destroyed {
		broadcastPayloadToAll(createUserDestroyedPayload(userDamagePayload.TargetId))
	}
}

func handleMobDamagePayload(conn *websocket.Conn, payload Payload) {
	var mobDamagePayload, err = createMobDamagePayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'MobDamagePayload' %s", err)
		return
	}

	var damaged, destroyed = dealDamageToMob(mobDamagePayload)

	if damaged {
		broadcastPayload(conn, mobDamagePayload)
	}

	if destroyed {
		resolveScore(conn, mobDamagePayload.Id, 20)
		broadcastPayloadToAll(createMobDestroyedPayload(mobDamagePayload.TargetId))
	}
}

func dealDamageToUser(payload UserDamagePayload) (bool, bool) {
	var damaged, destroyed = false, false
	players.Range(func(key, value any) bool {
		var player = value.(Player)
		if player.Id == payload.TargetId {
			player.dealDamage(payload.Damage)
			players.Store(key, player)
			if player.Destroyed == true {
				damaged = true
				destroyed = true
			} else {
				damaged = true
				destroyed = false
			}
		}

		return true
	})

	return damaged, destroyed
}

func dealDamageToMob(payload MobDamagePayload) (bool, bool) {
	var damaged = false
	var destroyed = false

	mobs.Range(func(key, value any) bool {
		var mob = value.(Mob)
		if mob.Id == payload.TargetId {
			mob.dealDamage(payload.Damage)
			if mob.Destroyed == true {
				damaged = true
				destroyed = true
			} else {
				damaged = true
				destroyed = false
			}
			mobs.Store(mob.Id, mob)
			return false
		} else {
			return true
		}
	})

	return damaged, destroyed
}

func deleteChest(id string) bool {
	var deleted = false

	chests.Range(func(key, value any) bool {
		var chest = value.(Chest)
		if chest.Id == id {
			chest.Destroyed = true
			deleted = true
			return false
		} else {
			return true
		}
	})

	return deleted
}

func containsGameMaster() bool {
	var containsMaster = false
	players.Range(func(key, value any) bool {
		var player = value.(Player)
		if player.Master == true {
			containsMaster = true
		}
		return true
	})

	return containsMaster
}

func findNonMaster() *Player {
	var master *Player = nil
	players.Range(func(key, value any) bool {
		var player = value.(Player)
		if player.Master == false && player.Destroyed == false {
			master = &player
			return false
		} else {
			return true
		}
	})

	return master
}
