// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorilla

import (
	"Gorilla_Websocket_Exercise/database"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Online struct {
	User    string
	Message string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	//connect to db
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	for {
		select {
		case client := <-h.register:
			h.Clients[client] = true
			if _, ok := h.Clients[client]; ok {
				var online1 Online
				online1.User = ""
				online1.Message = ""
				for cl := range h.Clients {
					w, err := cl.Conn.NextWriter(websocket.TextMessage)
					if err != nil {
						return
					}

					on, _ := json.Marshal(online1)
					w.Write(on)
				}
			}
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
					content := string(message)
					fmt.Println("message", message)

					m := strings.Split(content, " ")
					big := len(m)
					var online1 Online
					online1.User = strings.Trim(m[0], ":")
					online1.Message = strings.Join(m[1:big], " ")

					//upload message to db
					_, err = db.Exec("INSERT INTO messages (sender, content) VALUES (?,?)", online1.User, online1.Message)
					if err != nil {
						fmt.Println("error inserting message into db", err)
					}

				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
