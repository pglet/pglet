package page

import "encoding/json"

type Message struct {
	ID      string          `json:"id"`
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

func NewMessage(action string, payload interface{}) *Message {
	msg := &Message{
		Action: action,
	}

	// serialize payload
	serializedPayload, _ := json.Marshal(payload)
	msg.Payload = serializedPayload

	return msg
}
