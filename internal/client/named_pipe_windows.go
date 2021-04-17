// +build windows

package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/npipe"
)

const (
	readsize = 64 << 10
)

type namedPipe struct {
	conn            net.Conn
	id              string
	commandPipeName string
	commandListener *npipe.PipeListener
	eventPipeName   string
	eventListener   *npipe.PipeListener
	commands        chan string
	events          chan string
}

func newNamedPipe(id string) (*namedPipe, error) {
	pipeName := fmt.Sprintf("pglet_pipe_%s", id)

	pc := &namedPipe{
		id:              id,
		commandPipeName: pipeName,
		eventPipeName:   pipeName + ".events",
		commands:        make(chan string),
		events:          make(chan string, 2),
	}

	var err error
	pc.commandListener, err = npipe.Listen(`\\.\pipe\` + pc.commandPipeName)
	if err != nil {
		return nil, err
	}

	pc.eventListener, err = npipe.Listen(`\\.\pipe\` + pc.eventPipeName)
	if err != nil {
		return nil, err
	}

	go pc.commandLoop()
	go pc.eventLoop()

	return pc, nil
}

func (pc *namedPipe) getCommandPipeName() string {
	return pc.commandPipeName
}

func (pc *namedPipe) nextCommand() string {
	return <-pc.commands
}

func (pc *namedPipe) commandLoop() {
	log.Println("Starting command loop:", pc.commandPipeName)

	for {
		var err error
		pc.conn, err = pc.commandListener.Accept()
		if err != nil {
			log.Println("Command listener connection error:", err)
			return
		}

		log.Println("Connected to command pipe...")

		go func() {

			for {
				// read next command from pipeline
				cmdText := pc.read()

				if cmdText == "" {
					log.Println("Disconnected from command pipe")
					return
					//continue
				}

				pc.commands <- cmdText
			}

		}()
	}
}

func (pc *namedPipe) read() string {

	var bytesRead int
	var err error
	buf := make([]byte, readsize)

	r := bufio.NewReader(pc.conn)

	log.Debugln("Before read")

	for {
		var result []byte

		for {

			bytesRead, err = r.Read(buf)

			if err == io.EOF {
				return ""
			}

			result = append(result, buf[0:bytesRead]...)

			if bytesRead < readsize {
				break
			}
		}
		return strings.TrimSuffix(strings.TrimSuffix(string(result), "\n"), "\r")
	}
}

func (pc *namedPipe) writeResult(result string) {
	log.Debugln("Waiting for result to consume...")

	w := bufio.NewWriter(pc.conn)

	log.Debugln("Write result:", result)

	w.WriteString(fmt.Sprintf("%s\n", result))
	w.Flush()
}

func (pc *namedPipe) emitEvent(evt string) {

	//log.Debugln("Emit event:", evt)

	select {
	case pc.events <- evt:
		//log.Debugln("Event sent to queue:", evt)
	default:
		//log.Debugln("No event listeners:", evt)
	}
}

func (pc *namedPipe) eventLoop() {

	log.Println("Starting event loop:", pc.eventPipeName)

	for {
		conn, err := pc.eventListener.Accept()
		if err != nil {
			log.Println("Event listener connection error:", err)
			return
		}

		log.Println("Connected to event pipe...")

		go func(conn net.Conn) {

			defer log.Println("Disconnected from event pipe")

			for {
				select {
				case evt := <-pc.events:

					w := bufio.NewWriter(conn)

					_, err = w.WriteString(evt + "\n")
					if err != nil {
						if strings.Contains(err.Error(), "Pipe IO timed out waiting") {
							continue
						}
						log.Errorln("write error:", err)
						return
					}

					//conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
					err = w.Flush()
					if err != nil {
						if strings.Contains(err.Error(), "Pipe IO timed out waiting") {
							continue
						}
						log.Errorln("flush error:", err)
						return
					}

					log.Debugln("event written:", evt)

					// if !more {
					// 	return
					// }
				}
			}

		}(conn)
	}
}

func (pc *namedPipe) close() {
	log.Println("Closing Windows pipe...")

	pc.commandListener.Close()
	pc.eventListener.Close()
}
