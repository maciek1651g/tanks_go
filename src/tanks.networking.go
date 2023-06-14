package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
)

var connector = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients []websocket.Conn
var metadata = sync.Map{}

func migrateAddUser(w http.ResponseWriter, r *http.Request) {
	// Unmarshal the response body into a User struct
	fmt.Println("Initializing user")
	var user Player
	var err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal("Error decoding response:", err)
	}

	players.Store(user.Id, user)
}

func migrateUser(w http.ResponseWriter, r *http.Request) {
	var requestedId = r.URL.Query().Get("id")
	var requestedUrl = r.URL.Query().Get("url")
	metadata.Range(func(key, value any) bool {
		var idMapped = value.(string)
		var addressMapped = key.(net.Addr)
		if idMapped == requestedId {
			for _, client := range clients {
				if client.RemoteAddr() == addressMapped {
					var player, _ = players.Load(idMapped)
					handleUserDisconnection(&client)

					var response, _ = json.Marshal(player)

					var clientP = http.Client{}
					// Create a new HTTP request
					req, _ := http.NewRequest("POST", requestedUrl+"/api/users:migrate-add", bytes.NewBuffer(response))

					clientP.Do(req)
				}
			}

			return false
		} else {
			w.Write([]byte("User not found"))
			w.WriteHeader(404)
		}

		return true
	})
}

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
	sendCurrentMobs(conn)
	sendCurrentChests(conn)

	for true {
		_, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			return
		}

		fmt.Printf("%s : Received payload = '%s'\n", conn.RemoteAddr(), string(message))
		go handleActionPayload(conn, message)
	}
}

func saveUserStatus(client *websocket.Conn, id string, player Player) {
	players.Store(id, player)
	if player.Master {
		if err := client.WriteJSON(createGameMasterPayload(id)); err != nil {
			fmt.Printf("Error occurred when sending 'GameMasterPayload' : %s\n", err.Error())
		}
	}
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
	var data, found = players.Load(payload.Id)
	var player Player

	if found {
		player = data.(Player)
	} else {
		player = Player{Id: payload.Id, Destroyed: false, Coordinates: Coordinates{X: 200, Y: 600}, Health: 100, Master: !containsGameMaster()}
	}

	var initializationPayload = createPlayerCreatePayload(player)
	saveUserStatus(client, payload.Id, player)
	broadcastPayload(client, initializationPayload)
}

func sendCurrentUserStatuses(client *websocket.Conn) {
	var id, _ = metadata.Load(client.RemoteAddr())
	players.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %s \n value: %s \n", key, value)
		if id != key {
			fmt.Printf("Sending historical data to %s : %s\n", client.RemoteAddr(), value)
			var err = client.WriteJSON(createPlayerCreatePayload(value.(Player)))

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
			verifyMaster(id.(string))
			metadata.Delete(address)
			players.Delete(id)
		}
		return true
	})
}

func sendCurrentMobs(client *websocket.Conn) {
	for _, mob := range mobs {
		if mob.Destroyed == false {
			var payload = createMobCreatedPayload(mob)
			if err := client.WriteJSON(payload); err != nil {
				fmt.Printf("Error occurred when sending 'MobCreatedPayload' %s\n", payload)
			}
		}
	}
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

func verifyMaster(id string) {
	var player, _ = players.Load(id)
	if player.(Player).Master == true {
		var player = findNonMaster()
		if player != nil {
			var client = findClient(player.Id)
			player.Master = true
			saveUserStatus(client, player.Id, *player)
		}
	}
}

func sendUserDisconnection(client *websocket.Conn, id string) {
	var payload = UserDisconnectedPayload{Id: id, MessageType: "user_disconnected"}
	fmt.Printf("Broadcasting 'user_disconnected' payload for id = %s\n", id)
	broadcastPayload(client, payload)
}

func broadcastPayload(except *websocket.Conn, payload any) {
	for _, client := range clients {
		if except == nil || client.RemoteAddr() != except.RemoteAddr() {
			fmt.Printf("%s : Sending payload : '%s'\n", client.RemoteAddr(), payload)
			if err := client.WriteJSON(payload); err != nil {
				fmt.Printf(err.Error() + "\n")
				return
			}
		}
	}
}

func broadcastPayloadToAll(payload any) {
	broadcastPayload(nil, payload)
}

func findClient(id string) *websocket.Conn {

	for _, client := range clients {
		var playerId, _ = metadata.Load(client.RemoteAddr())
		if playerId == id {
			return &client
		}
	}

	return nil
}
