//go:build !windows

package connection

import (
	"github.com/pglet/pglet/internal/page"
	page_connection "github.com/pglet/pglet/internal/page/connection"
	log "github.com/sirupsen/logrus"
)

type Local struct {
	readCh  chan []byte
	writeCh chan []byte
}

func NewLocal() *Local {
	cws := &Local{
		readCh:  make(chan []byte),
		writeCh: make(chan []byte, 10),
	}
	return cws
}

func (c *Local) Start(readHandler ReadMessageHandler, reconnectHandler ReconnectHandler) (err error) {

	log.Println("Connecting to local Pglet Server")

	// create page client
	cl := page_connection.NewLocal(c.writeCh, c.readCh)
	page.NewClient(cl, "", nil)

	// start read loop
	go c.readLoop(readHandler)

	return
}

func (c *Local) readLoop(readHandler ReadMessageHandler) {
	for {
		message := <-c.readCh
		err := readHandler(message)
		if err != nil {
			log.Errorf("error processing message: %v", err)
			break
		}
	}
}

func (c *Local) Send(message []byte) {
	c.writeCh <- message
}
