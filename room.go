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
			log.Printf("Client joined room %s. Total clients: %d", r.name, len(r.clients))

		// removing a user from a channel
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.recieve)
			log.Printf("Client left room %s. Total clients: %d", r.name, len(r.clients))

		// sending a message to all clients in the room
		case msg := <-r.forward:
			log.Printf("Broadcasting message to %d clients in room %s", len(r.clients), r.name)
			for client := range r.clients {
				select {
				case client.recieve <- msg:
					// Message sent successfully
				default:
					// Channel is full or closed, remove the client
					delete(r.clients, client)
					close(client.recieve)
					log.Printf("Removed client due to channel issue in room %s", r.name)
				}
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