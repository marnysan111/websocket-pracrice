package models

import "github.com/gorilla/websocket"

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan Message)

type Message struct {
	Type    int
	Message []byte
}
