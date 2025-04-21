package webconnect

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Id  string
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

func NewClient(newId string, newHub *Hub, newConn *websocket.Conn) *Client {
	return &Client{
		Id:   newId,
		hub:  newHub,
		conn: newConn,
		send: make(chan []byte, 10),
	}
}

func (c *Client) Run() {
	// Write new message permanently
	go func() {
		for {
			select {
			case message := <-c.send:
				err := c.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					return
				}
			case <-time.After(time.Second * 2):
				err := c.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello, %s!", c.Id)))
				if err != nil {
					return
				}
			}
		}
	}()
	// Read messsage
	// and close on disconnect
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("|> Socket read error:", err)
			c.conn.Close()
			break
		}
		c.send <- message
	}

}
