package main

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

// naming is from the client viewpoint
// send_message are messages sent by client, recieved by server
// recieve_message are messages recieved by client, sent by server
const (
	EventSendMessage    = "send_message"
	EventRecieveMessage = "recieve_message"
	EventChangeRoom     = "change_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json"from"`
}

type RecieveMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json"sent"`
}

type ChangeRoomEvent struct {
	Name string `json"name"`
}
