package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)
type room struct {
	// holds the clients that are connected to the room
	clients map[*client]bool

	// join is a channel for clients to join the room
	join chan *client

	// leave is a channel for clients to leave the room
	leave chan *client

	// forward is a channel for messages to be forwarded to all clients in the room
	forward chan []byte

	// room name
	name string
}

func newRoom(name string) *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		name: name,
	}
}

//function that will handle the user joining and leaving the room
func (r *room) run() {
	for{
		select{
		// adding a user to a channel
		case client := <-r.join:
			r.clients[client] = true

		// removing a user from a channel
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.recieve)
		// sending a message to all clients in the room
		case msg := <-r.forward:
			for client := range r.clients {
				client.recieve <- msg
			}
		}
	}
}

// upgrade the http connection to a websocket connection

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	
	log.Printf("Client connected to room: %s", r.name)
	
	client := &client{
		socket: socket,
		recieve: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client

	defer func() { 
		log.Printf("Client disconnected from room: %s", r.name)
		r.leave <- client 
	}()
	go client.read()
	go client.write()
}