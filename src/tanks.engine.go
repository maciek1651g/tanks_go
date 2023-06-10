package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

var chests []Chest
var mobs []Mob
var players = sync.Map{}

func initializeEngine() {
	createMobs()
	createChests()
}

func createMobs() {
	mobs = []Mob{
		createStandardMob(uuid.New().String(), Coordinates{X: 200, Y: 200}),
		createStandardMob(uuid.New().String(), Coordinates{X: 300, Y: 300}),
		createStandardMob(uuid.New().String(), Coordinates{X: 400, Y: 400}),
		createStandardMob(uuid.New().String(), Coordinates{X: 500, Y: 500}),
		createStandardMob(uuid.New().String(), Coordinates{X: 600, Y: 600}),
		createStandardMob(uuid.New().String(), Coordinates{X: 700, Y: 700}),
		createStandardMob(uuid.New().String(), Coordinates{X: 800, Y: 800}),
		createStandardMob(uuid.New().String(), Coordinates{X: 900, Y: 900}),
		createStandardMob(uuid.New().String(), Coordinates{X: 1000, Y: 1000}),
		createStandardMob(uuid.New().String(), Coordinates{X: 1100, Y: 1100}),
		createStandardMob(uuid.New().String(), Coordinates{X: 550, Y: 600}),
	}
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

func handleActionPayload(conn *websocket.Conn, client *mongo.Client, bytes []byte) {
	var payload Payload
	var error = json.Unmarshal(bytes, &payload)

	if error != nil {
		fmt.Printf("Could not deserialize payload %s", error)
		return
	}

	switch payload.MessageType {
	case "status":
		handleStatusPayload(conn, client, payload)
	case "user_attack":
		handleUserAttackPayload(conn, client, payload)
	case "chest_grab":
		handleChestGrabPayload(conn, client, payload)
	case "user_damage":
		handleUserDamagePayload(conn, client, payload)
	case "mob_status":
		handleMobStatusPayload(conn, client, payload)
	case "mob_damage":
		handleMobDamagePayload(conn, client, payload)
	}
}

func handleStatusPayload(conn *websocket.Conn, client *mongo.Client, payload Payload) {
	var userStatusPayload, err = createUserStatusPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'UserStatusPayload' %s", err)
		return
	}

	updateCoordinates(conn, userStatusPayload.Coordinates)
	broadcastPayload(conn, userStatusPayload)
	go saveInMongo(userStatusPayload, client, "users_statuses")
}

func handleMobStatusPayload(conn *websocket.Conn, client *mongo.Client, payload Payload) {
	var mobStatusPayload, err = createMobStatusPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'MobStatusPayload' %s", err)
		return
	}

	if updateMobCoordinates(mobStatusPayload.Id, mobStatusPayload.Coordinates) {
		broadcastPayload(conn, mobStatusPayload)
	}

	go saveInMongo(mobStatusPayload, client, "mob_statuses")
}

func updateMobCoordinates(id string, coordinates Coordinates) bool {
	for index, mob := range mobs {
		if mob.Id == id {
			mobs[index].Coordinates = coordinates
			return true
		}
	}

	return false
}

func handleUserAttackPayload(conn *websocket.Conn, client *mongo.Client, payload Payload) {
	var userAttackPayload, err = createUserAttackPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'UserAttackPayload' %s", err)
		return
	}

	broadcastPayload(conn, userAttackPayload)
	go saveInMongo(userAttackPayload, client, "user_attacks")
}

func updateCoordinates(conn *websocket.Conn, coordinates Coordinates) {
	var id, _ = metadata.Load(conn.RemoteAddr())
	var status, _ = players.Load(id)
	var currentStatus = status.(Player)
	currentStatus.Coordinates = coordinates
	players.Store(id, currentStatus)
}

func handleChestGrabPayload(conn *websocket.Conn, client *mongo.Client, payload Payload) {
	var chestGrabPayload, err = createChestGrabPayload([]byte(payload.Data))

	if err != nil {
		fmt.Printf("Could not deserialize 'ChestGrabPayload' %s", err)
		return
	}

	if deleteChest(chestGrabPayload.Id) == true {
		resolveScore(conn, chestGrabPayload.PlayerId, 10)
		broadcastPayload(conn, createChestDestroyedPayload(chestGrabPayload.Id))
		go saveInMongo(createChestDestroyedPayload(chestGrabPayload.Id), client, "chest_destroyed")
	}

	go saveInMongo(chestGrabPayload, client, "chest_grabbed")
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

func handleUserDamagePayload(conn *websocket.Conn, client *mongo.Client, payload Payload) {
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
		go saveInMongo(createUserDestroyedPayload(userDamagePayload.TargetId), client, "user_destroyed")
	}
}

func handleMobDamagePayload(conn *websocket.Conn, client *mongo.Client, payload Payload) {
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
		go saveInMongo(createMobDestroyedPayload(mobDamagePayload.TargetId), client, "mob_destroyed")
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
	for index, mob := range mobs {
		if mob.Id == payload.TargetId {
			mobs[index].dealDamage(payload.Damage)
			if mobs[index].Destroyed == true {
				return true, true
			} else {
				return true, false
			}
		}
	}

	return false, false
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
