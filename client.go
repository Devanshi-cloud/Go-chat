package main

import "github.com/gorilla/websocket"
type client struct {
	socket *websocket.Conn

	recieve chan []byte

	room *room
}

//close the connection when the client is disconnected
func (c *client) read() {
	//as long as the client is connected, read the message from the socket and forward it to the room
	for{
		_, msg,err := c.socket.ReadMessage()

		if err != nil {
			return
	}

	c.room.forward <- msg
}
}

// used to recieved messages
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.recieve {
		err:=c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
	}
}
}