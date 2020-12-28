package connection

import (
	"github.com/pglet/pglet/internal/page"
	page_connection "github.com/pglet/pglet/internal/page/connection"
	log "github.com/sirupsen/logrus"
)

type Local struct {
	readCh  chan []byte
	writeCh chan []byte
	done    chan bool
}

func NewLocal() *Local {
	cws := &Local{
		readCh:  make(chan []byte),
		writeCh: make(chan []byte, 10),
		done:    make(chan bool),
	}
	return cws
}

func (c *Local) Start(handler ReadMessageHandler) (err error) {

	log.Println("Connecting to local Pglet Server")

	// create page client
	cl := page_connection.NewLocal(c.writeCh, c.readCh)
	page.NewClient(cl)

	// start read loop
	go c.readLoop(handler)

	return
}

func (c *Local) readLoop(readHandler ReadMessageHandler) {
	for {
		message := <-c.readCh
		err := readHandler(message)
		if err != nil {
			log.Printf("error processing message: %v", err)
			break
		}
	}
}

func (c *Local) Send(message []byte) {
	c.writeCh <- message
}
