package model

import (
	"bytes"
	"encoding/json"
)

type Message struct {
	Sender   int  `json:"sender"`
	Receiver int  `json:"receiver"`
	Message  string `json:"message"`
	Cmd      string `json:"cmd"`
	Todo     string `json:"todo"`
}

func (m Message) GetCmd() string { return m.Cmd }
func (m Message) GetSender() int   { return m.Sender }
func (m Message) GetReceiver() int { return m.Receiver }
func (m Message) GetMessage() string { return m.Message }
func (m Message) GetTodo() string    { return m.Todo }

func ParseToModel(raw_message []byte) Message {
	var m Message

	err := json.NewDecoder(bytes.NewReader(raw_message)).Decode(&m)
	if err != nil {
		panic("Exception during parsing to Message model")
	}
	return m
}

func TransformToJson(message Message) []byte {
	transform, error := json.Marshal(message)
	if error != nil {
		panic("Cannot Transform Message To String")
	}
		
	return transform
}