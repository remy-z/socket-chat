package main

import (
	"encoding/json"
	"time"
)

// constants for routing on both front and backend
// naming is from the client viewpoint
// send_message are messages sent by client, recieved by server
// recieve_message are messages recieved by client, sent by server
const (
	EventSendMessage    = "send_message"
	EventRecieveMessage = "recieve_message"
	EventChangeRoom     = "change_room"
	EventSendDir        = "send_dir"
	EventRecieveDir     = "recieve_dir"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type RecieveMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeRoomEvent struct {
	Name string `json:"name"`
}

type RecieveDirEvent struct {
	Chats bool `json:"chats"`
	Users bool `json:"users"`
}

type SendDirEvent struct {
	Users []string       `json:"users"`
	Chats map[string]int `json:"chats"`
}
