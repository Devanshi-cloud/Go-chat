package main

import (
	"encoding/json"
	"log"
	"github.com/gorilla/websocket"
)

type message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Room    string `json:"room"`
}

type client struct {
	socket *websocket.Conn
	recieve chan []byte
	room *room
	name string
}

//close the connection when the client is disconnected
func (c *client) read() {
	//as long as the client is connected, read the message from the socket and forward it to the room
	for{
		_, msg,err := c.socket.ReadMessage()

		if err != nil {
			log.Printf("Client read error: %v", err)
			return
		}

		log.Printf("Received message from client: %s", string(msg))

		// Parse the message as JSON
		var messageData message
		if err := json.Unmarshal(msg, &messageData); err != nil {
			// If not JSON, treat as plain text
			messageData = message{
				Name:    c.name,
				Message: string(msg),
				Room:    c.room.name,
			}
		}

		// Set the client name if not set
		if c.name == "" {
			c.name = messageData.Name
		}

		// Forward the message to the room
		c.room.forward <- msg
	}
}

// used to recieved messages
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.recieve {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}