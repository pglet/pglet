//go:build !windows

package connection

type ReadMessageHandler func(message []byte) error
type ReconnectHandler func(success bool)

type Conn interface {
	Start(readHandler ReadMessageHandler, reconnectHandler ReconnectHandler) (err error)
	Send(message []byte)
}
