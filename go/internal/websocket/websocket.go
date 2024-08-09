package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // すべてのオリジンを許可
	},
}

var Broadcast = make(chan Message)

// Message represents the message structure
type Message struct {
	Type    int    `json:"type"`
	Message []byte `json:"message"`
	RoomID  string
}

type Room struct {
	RoomID   string
	UUID     string
	Clients  map[*websocket.Conn]bool
	Messages []Message
}

type RoomDetail struct {
	RoomID string
	UUID   string
}

var Rooms = make(map[string]*Room)

func NewRoom(roomID string) *Room {
	roomUUID := uuid.New().String()
	room := &Room{
		UUID:     roomUUID,
		RoomID:   roomID,
		Clients:  make(map[*websocket.Conn]bool),
		Messages: make([]Message, 0),
	}
	Rooms[roomID] = room
	return room
}

func GetRoom(roomID string) *Room {
	room, exists := Rooms[roomID]
	if !exists {
		room = NewRoom(roomID)
		Rooms[roomID] = room
	}
	return room
}

func GetRoomList() []RoomDetail {
	for _, r := range Rooms {
		fmt.Println(r.Messages)
	}
	roomList := []RoomDetail{}
	for _, room := range Rooms {
		roomDetail := RoomDetail{RoomID: room.RoomID, UUID: room.UUID}
		roomList = append(roomList, roomDetail)
	}
	return roomList
}

// ConnHandler handles WebSocket connections
func ConnHandler(c *gin.Context, roomID string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}
	room := GetRoom(roomID)
	room.Clients[conn] = true

	go HandleMessages() // Start handling messages
	for _, msg := range room.Messages {
		err := conn.WriteMessage(msg.Type, msg.Message)
		if err != nil {
			log.Printf("Error sending past messages: %v", err)
			conn.Close()
			delete(room.Clients, conn)
			return
		}
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Unexpected WebSocket close error: %v", err)
			delete(room.Clients, conn)
			conn.Close()
			break
		}
		message := Message{Type: t, Message: msg, RoomID: roomID}
		room.Messages = append(room.Messages, message)
		Broadcast <- message
	}
}

// handleMessages handles incoming messages and broadcasts them to all clients
func HandleMessages() {
	for {
		message := <-Broadcast
		room, exists := Rooms[message.RoomID]
		if !exists {
			continue
		}
		for client := range room.Clients {
			err := client.WriteMessage(message.Type, message.Message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(room.Clients, client)
			}
		}
	}
}
