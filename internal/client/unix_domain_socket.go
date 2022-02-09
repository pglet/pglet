//go:build !windows

package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	coreFxSocketPrefix = "CoreFxPipe_"
)

type unixDomainSocket struct {
	conn            net.Conn
	id              string
	commandPipeName string
	commandListener net.Listener
	eventPipeName   string
	eventListener   net.Listener
	commands        chan string
	events          chan string
}

func newUnixDomainSocket(id string) (*unixDomainSocket, error) {
	pipeName := fmt.Sprintf("pglet_pipe_%s", id)

	pc := &unixDomainSocket{
		id:              id,
		commandPipeName: pipeName,
		eventPipeName:   pipeName + ".events",
		commands:        make(chan string),
		events:          make(chan string, 2),
	}

	go pc.commandLoop()
	go pc.eventLoop()

	return pc, nil
}

func (pc *unixDomainSocket) getCommandPipeName() string {
	return pc.commandPipeName
}

func (pc *unixDomainSocket) nextCommand() string {
	return <-pc.commands
}

func (pc *unixDomainSocket) commandLoop() {
	log.Println("Starting command loop:", pc.commandPipeName)

	var err error
	sockAddr := filepath.Join(os.TempDir(), coreFxSocketPrefix+pc.commandPipeName)
	pc.commandListener, err = net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatalln("Error listening command loop Unix domain socket", sockAddr, err)
	}

	for {
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

func (pc *unixDomainSocket) read() string {

	var bytesRead int
	var err error
	buf := make([]byte, readsize)

	r := bufio.NewReader(pc.conn)

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

func (pc *unixDomainSocket) writeResult(result string) {
	log.Debugln("Waiting for result to consume...")

	w := bufio.NewWriter(pc.conn)

	log.Debugln("Write result:", result)

	w.WriteString(fmt.Sprintf("%s\n", result))
	w.Flush()
}

func (pc *unixDomainSocket) emitEvent(evt string) {
	//log.Debugln("Emit event:", evt)

	select {
	case pc.events <- evt:
		log.Debugln("Event sent to queue:", evt)
	default:
		log.Debugln("No event listeners:", evt)
	}
}

func (pc *unixDomainSocket) eventLoop() {

	log.Println("Starting event loop:", pc.eventPipeName)

	var err error
	sockAddr := filepath.Join(os.TempDir(), coreFxSocketPrefix+pc.eventPipeName)
	pc.eventListener, err = net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatalln("Error listening event loop Unix domain socket", sockAddr, err)
	}

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
				evt, more := <-pc.events

				w := bufio.NewWriter(conn)

				_, err = w.WriteString(evt + "\n")
				if err != nil {
					if strings.Contains(err.Error(), "Pipe IO timed out waiting") {
						continue
					}
					return
				}

				//conn.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
				err = w.Flush()
				if err != nil {
					if strings.Contains(err.Error(), "Pipe IO timed out waiting") {
						continue
					}
					return
				}

				log.Debugln("event written:", evt)

				if !more {
					return
				}
			}

		}(conn)
	}
}

func (pc *unixDomainSocket) close() {
	log.Println("Closing Unix domain socket...")

	pc.commandListener.Close()
	pc.eventListener.Close()
}
