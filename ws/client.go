package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
)

type Client struct {
	ID         string
	Connection *websocket.Conn
}

type Message struct {
	Action  string          `json:"action"`
	Channel string          `json:"channel"`
	Message json.RawMessage `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func autoID() string {

	return uuid.New().String()
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		ID:         autoID(),
		Connection: conn,
	}

	// add this client into the list
	//ps.AddClient(client)

	fmt.Printf("New Client %s is connected, total: %d", client.ID, 0)

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Something went wrong", err)

			//ps.RemoveClient(client)
			log.Println("total clients and subscriptions ", 0, 0)

			return
		}

		fmt.Println(messageType, msg)

		conn.WriteMessage(messageType, msg)

		//ps.HandleReceiveMessage(client, messageType, p)
	}
}
