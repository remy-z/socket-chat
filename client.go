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
	startingRoom = "general"
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event //egress to avoid concurrent writes
	chatroom   string
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	c := &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		chatroom:   startingRoom,
	}

	c.setupClient()
	return c
}

func (c *Client) setupClient() {
	err := c.connection.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Println(err)
		return
	}

	//No Jumbo Frames!
	c.connection.SetReadLimit(512)

	c.connection.SetPongHandler(c.pongHandler)
}

// Event loop for reading incoming messages from client on websocket connection
func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

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
			fmt.Println("error handling message: ", err)
		}
	}
}

// Event loop for sending pings and messages to client
func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	for {
		select {

		//write message to websocket connection whenever a messages comes from egress channel
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

			//Send a ping to client whenever there is an incoming tick
		case <-ticker.C:
			err := c.connection.WriteMessage(websocket.PingMessage, []byte(``))
			if err != nil {
				fmt.Println("write message err: ", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	//fmt.Println("pong")
	//reset timer when pong recieved
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
