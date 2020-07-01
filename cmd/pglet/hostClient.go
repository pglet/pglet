package main

import "github.com/gorilla/websocket"

type hostClient struct {

	// active WebSocket connection
	conn *websocket.Conn

	// clients by client ID
	clients map[string]*client
}
