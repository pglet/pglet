package client

type ReadMessageHandler func(message []byte) error

type Conn interface {
	Start(handler ReadMessageHandler) (err error)
	Send(message []byte)
}
