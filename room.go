package main

import (
	"log"
	"net/http"

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
					//send the message to the client's receive channel
					client.recieve <- msg
				}
		}
	}
}

// upgrade a basic http connection to a websocket connection
const(
	// The maximum size of a message that can be received from the client
	socketBufferSize = 1024 // 4 // 256 bytes
	// The maximum size of a message that can be sent to the client
	messageBufferSize = 256 // 1024 bytes
)

var upgrader = websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func (r *room)ServerHttp(w http.ResponseWriter, req *http.Request) {
	//upgrade the connection to a websocket connection
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServerHTTP:", err)
		return
	}

	//create a new client and add it to the room
	client := &client{
		socket:  socket,
		recieve: make(chan []byte, messageBufferSize),
		room:    r,
	}

	r.join <- client //add the client to the room

	defer func() {
		r.leave <- client //remove the client from the room when the function exits
	}()

	go client.write() //start writing messages to the client
	client.read() //start reading messages from the client
}