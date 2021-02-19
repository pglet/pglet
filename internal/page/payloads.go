package page

import (
	"encoding/json"

	"github.com/pglet/pglet/internal/model"
	"github.com/pglet/pglet/internal/page/command"
)

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

type RegisterHostClientRequestPayload struct {
	PageName string `json:"pageName"`
	IsApp    bool   `json:"isApp"`
}

type RegisterHostClientResponsePayload struct {
	SessionID string `json:"sessionID"`
	PageName  string `json:"pageName"`
	Error     string `json:"error"`
}

type RegisterWebClientRequestPayload struct {
	PageName  string `json:"pageName"`
	PageHash  string `json:"pageHash"`
	SessionID string `json:"sessionID"`
}

type RegisterWebClientResponsePayload struct {
	Session *SessionPayload `json:"session"`
	Error   string          `json:"error"`
}

type SessionPayload struct {
	ID       string                    `json:"id"`
	Controls map[string]*model.Control `json:"controls"`
}

type SessionCreatedPayload struct {
	PageName  string `json:"pageName"`
	SessionID string `json:"sessionID"`
}

type PageCommandRequestPayload struct {
	PageName  string          `json:"pageName"`
	SessionID string          `json:"sessionID"`
	Command   command.Command `json:"command"`
}

type PageCommandResponsePayload struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

type PageEventPayload struct {
	PageName    string `json:"pageName"`
	SessionID   string `json:"sessionID"`
	EventTarget string `json:"eventTarget"`
	EventName   string `json:"eventName"`
	EventData   string `json:"eventData"`
}

type AddPageControlsPayload struct {
	Controls []*model.Control `json:"controls"`
	TrimIDs  []string         `json:"trimIDs"`
}

type ReplacePageControlsPayload struct {
	IDs      []string         `json:"ids"`
	Remove   bool             `json:"remove"`
	Controls []*model.Control `json:"controls"`
}

type UpdateControlPropsPayload struct {
	Props []map[string]interface{} `json:"props"`
}

type AppendControlPropsPayload struct {
	Props []map[string]string `json:"props"`
}

type RemoveControlPayload struct {
	IDs []string `json:"ids"`
}

type CleanControlPayload struct {
	IDs []string `json:"ids"`
}
