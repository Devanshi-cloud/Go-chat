package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type room struct {

	//hold track of all connected clients
	clients map[*client]bool

	//join is a channel for clients to join the room
	join chan *client

	//leave is a channel for clients to leave the room
	leave chan *client

	//forward is a channel that holds incoming messages that should be sent to clients
	forward chan []byte
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

//each room is a separate thread that should run independently (but as long as the server is running)
func (r *room) run() {
	for {
		select{
			// adding a user to a channel
			case client := <-r.join:
				r.clients[client] = true //add the client to the room
			case client := <-r.leave:
				delete(r.clients, client) //remove the client from the room
				close(client.recieve) //close the client's receive channel

			//send a message to all clients in the room
			case msg := <-r.forward:
				for client := range r.clients {
					select {
					case client.recieve <- msg:
						// message sent successfully
					default:
						// client's receive channel is full or closed, remove client
						delete(r.clients, client)
						close(client.recieve)
					}
				}
		}
	}
}

// getRoom retrieves a room by name or creates a new one if it doesn't exist
var rooms = make(map[string]*room)
var mu sync.Mutex // mutex to protect the rooms map

func getRoom(name string) *room {

	// prevent creating a room with an same name when multiple users do that at the same time

	mu.Lock()
	defer mu.Unlock()
	if r, ok := rooms[name]; ok {
		return r //return the existing room
	}

	//create a new room if it doesn't exist
	r := newRoom()
	rooms[name] = r
	go r.run() // Start the room's goroutine
	return r
}

// upgrade a basic http connection to a websocket connection
const(
	// The maximum size of a message that can be received from the client
	socketBufferSize = 1024 // 4 // 256 bytes
	// The maximum size of a message that can be sent to the client
	messageBufferSize = 256 // 1024 bytes
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  socketBufferSize, 
	WriteBufferSize: messageBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (adjust for production)
	},
}

// ServeHTTP implements the http.Handler interface
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	roomName := req.URL.Query().Get("room")
	if roomName == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}
	realRoom := getRoom(roomName) //create a new room if it doesn't exist

	//upgrade the connection to a websocket connection
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Printf("ServeHTTP upgrade error: %v", err)
		return
	}

	//create a new client and add it to the room
	client := &client{
		socket:  socket,
		recieve: make(chan []byte, messageBufferSize),
		room:    realRoom, // Use realRoom instead of r
		name: fmt.Sprintf("User%d", rand.Intn(1000)), //assign a name to the client
	}

	realRoom.join <- client //add the client to the room

	defer func() {
		realRoom.leave <- client //remove the client from the room when the function exits
	}()

	go client.write() //start writing messages to the client
	client.read() //start reading messages from the client
}