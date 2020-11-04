package client

type Conn interface {
	Start()
	ReadMessage() []byte
	Send(message []byte)
}
