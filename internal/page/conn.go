package page

type ReadMessageHandler func(message []byte) error

type Conn interface {
	Start(handler ReadMessageHandler)
	Send(message []byte)
}
