package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	chatroom string
	//egress to avoid concurrent writes
	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		//cleanup connection when connection breaks
		c.manager.removeClient(c)
	}()

	err := c.connection.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Println(err)
		return
	}

	c.connection.SetReadLimit(512)

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()
		// error when conn closed
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error read message: %v", err)
			}
			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			fmt.Printf("error marshalling event: %v", err)
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			fmt.Println("error handeling message: ", err)
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				//Server write message to client when server has egress issues that connection has to be closed
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					fmt.Println("connection closed: ", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				fmt.Printf("failed to send message: %v", err)
			}
			fmt.Println("message sent")

		case <-ticker.C:
			fmt.Println("ping")

			err := c.connection.WriteMessage(websocket.PingMessage, []byte(``))
			if err != nil {
				fmt.Println("write message err: ", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	fmt.Println("pong")
	//reset timer when pong recieved
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
