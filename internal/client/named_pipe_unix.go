// +build !windows

package client

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const (
	readsize = 64 << 10
)

type namedPipe struct {
	id              string
	pageName        string
	sessionID       string
	commandPipeName string
	eventPipeName   string
	events          chan string
}

func newNamedPipe(id string) (*namedPipe, error) {
	pipeName := filepath.Join(os.TempDir(), fmt.Sprintf("pglet_pipe_%s", id))

	pc := &namedPipe{
		id:              id,
		commandPipeName: pipeName,
		eventPipeName:   pipeName + ".events",
		events:          make(chan string, 2),
	}

	return pc, pc.start()
}

func (pc *namedPipe) getCommandPipeName() string {
	return pc.commandPipeName
}

func (pc *namedPipe) start() error {
	// create "command" named pipe
	err := createFifo(pc.commandPipeName)
	if err != nil {
		return err
	}

	// create "events" named pipe
	err = createFifo(pc.eventPipeName)
	if err != nil {
		return err
	}

	go pc.eventLoop()

	return nil
}

func (pc *namedPipe) nextCommand() string {
	return pc.read()
}

func (pc *namedPipe) read() string {
	var bytesRead int
	buf := make([]byte, readsize)
	for {
		var result []byte
		input, err := openFifo(pc.commandPipeName, os.O_RDONLY)
		if err != nil {
			break
		}
		for err == nil {
			bytesRead, err = input.Read(buf)
			result = append(result, buf[0:bytesRead]...)

			if err == io.EOF {
				break
			}

			//log.Printf("read: %d\n", bytesRead)
		}
		input.Close()
		return string(result)
	}
	return ""
}

func (pc *namedPipe) writeResult(result string) {
	log.Println("Waiting for result to consume...")
	output, err := openFifo(pc.commandPipeName, os.O_WRONLY)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Write result:", result)

	output.WriteString(fmt.Sprintf("%s\n", result))
	output.Close()
}

func (pc *namedPipe) emitEvent(evt string) {
	select {
	case pc.events <- evt:
		// Event sent to queue
	default:
		// No event listeners
	}
}

func (pc *namedPipe) eventLoop() {

	log.Println("Starting event loop...")

	defer os.Remove(pc.eventPipeName)

	for {
		output, err := openFifo(pc.eventPipeName, os.O_WRONLY)
		if err != nil {
			log.Fatal(err)
		}

		select {
		case evt, more := <-pc.events:
			output.WriteString(evt + "\n")
			output.Close()

			if !more {
				return
			}
		}
	}
}

func (pc *namedPipe) close() {
	log.Println("Closing Unix pipe...")

	os.Remove(pc.commandPipeName)
	os.Remove(pc.eventPipeName)
}

func createFifo(filename string) (err error) {
	err = syscall.Mkfifo(filename, 0660)
	return
}

func openFifo(path string, oflag int) (f *os.File, err error) {
	f, err = os.OpenFile(path, oflag, os.ModeNamedPipe)
	return
}
