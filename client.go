package main

import "github.com/gorilla/websocket"

// client represents a single connected client in the chat room.
type client struct {
	//a websocket connection for this client
	socket *websocket.Conn

	//recieve is a channel to receive messages from the client
	recieve chan []byte

	room *room //the room this client is in
}

//send messages function sends messages to the client.
func (c *client) read(){

	//close the connection when the function exits
	defer c.socket.Close()

	for{
		//read message from the client as long as their is an input, forward it to the room
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return //if there's an error, exit the loop
		}

		//forward the message to the room's forward channel
		c.room.forward <- msg
	}
}

//used to recieve messages from the room and send them to the client
func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.recieve {
		err:=c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return //if there's an error, exit the loop
		}
	}
}