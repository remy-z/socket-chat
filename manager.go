package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex

	otps RetentionMap

	handlers map[string]EventHandler
}

type ClientList map[*Client]bool

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
		otps:     NewRetentionMap(ctx, 5*time.Second),
	}

	m.setupEventhandlers()
	return m
}

// handlers stored in a map to route events to the correct handler
// based on event type (using m.routeEvent)
func (m *Manager) setupEventhandlers() {

	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventChangeRoom] = ChangeRoomHandler

}

// Route events triggered from the client sending messages through websocket
func (m *Manager) routeEvent(event Event, c *Client) error {
	//check if event type is included in handlers
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("event type not found")
	}
}

// Envoked when a client sends an event with type EventSendMessage
// TODO write message to DB and save in recent queue
// Broadcast to other clients in the same room
func SendMessageHandler(event Event, c *Client) error {
	var chatevent SendMessageEvent

	err := json.Unmarshal(event.Payload, &chatevent)
	if err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	var broadMessage RecieveMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatevent.Message
	broadMessage.From = chatevent.From

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	outgoingEvent := Event{
		Payload: data,
		Type:    EventRecieveMessage,
	}

	//write messages to egress channel on clients
	for client := range c.manager.clients {
		if client.chatroom == c.chatroom {
			client.egress <- outgoingEvent
		}

	}
	return nil
}

// Envoked when a client sends an event with type EventChangeRoom
func ChangeRoomHandler(event Event, c *Client) error {
	var changeRoomEvent ChangeRoomEvent
	err := json.Unmarshal(event.Payload, &changeRoomEvent)
	if err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	c.chatroom = changeRoomEvent.Name
	return nil
}

// http request handler for /ws, upgrades connection to websocket
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	//verify OTP
	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("new connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	client := NewClient(conn, m)

	m.addClient(client)

	// Start read and write go routines
	go client.readMessages()
	go client.writeMessages()
}

// http handler that is set in main.go to authenticate login requests form client
func (m *Manager) loginHandler(w http.ResponseWriter, r *http.Request) {

	var req userLoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO real authentification stuff here
	if req.Password == "123" {
		otp := m.otps.NewOTP()

		resp := otpResponse{
			OTP: otp.Key,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

// add client to manager's ClientList
// add client to chatroom map
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

// close the client's connection and remove client from memory
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

// return true to allow connection, false to dismiss
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "https://localhost:8080":
		return true
	default:
		return false
	}
}
