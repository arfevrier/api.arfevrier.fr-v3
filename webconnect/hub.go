package webconnect

import (
	"fmt"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Registered clients.
	rtcs map[*Rtc]bool

	// Register requests from the clients.
	RegisterRtc chan *Rtc

	// Unregister requests from clients.
	unregisterRtc chan *Rtc
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:     make(chan []byte),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Clients:       make(map[*Client]bool),
		RegisterRtc:   make(chan *Rtc),
		unregisterRtc: make(chan *Rtc),
		rtcs:          make(map[*Rtc]bool),
	}
}

func (h *Hub) Run() {
	// Default broadcast message every 2 seconds
	go func() {
		for {
			time.Sleep(time.Second * 5)
			h.Broadcast <- []byte("Broadcast message")
		}

	}()
	// Loop for managing the hub
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
				client.conn.Close()
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					fmt.Println("!! Warning closing Socket because of broadcast")
					close(client.send)
				}
			}
			for rtc := range h.rtcs {
				select {
				case rtc.send <- message:
				default:
					fmt.Println("!! Warning closing RTC because of broadcast")
					close(rtc.send)
				}
			}
		case rtc := <-h.RegisterRtc:
			h.rtcs[rtc] = true
		case rtc := <-h.unregisterRtc:
			if _, ok := h.rtcs[rtc]; ok {
				delete(h.rtcs, rtc)
				close(rtc.send)
				rtc.Conn.Close()
			}
		}
	}
}
