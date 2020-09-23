// +build !windows

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"syscall"
)

const (
	readsize = 64 << 10
)

type pipeImpl struct {
	id              string
	pageName        string
	sessionID       string
	commandPipeName string
	eventPipeName   string
	commands        chan string
	events          chan string
}

func newPipeImpl(id string) (*pipeImpl, error) {
	pipeName := path.Join(os.TempDir(), fmt.Sprintf("pglet_pipe_%s", id))

	pc := &pipeClient{
		id:              id,
		commandPipeName: pipeName,
		eventPipeName:   pipeName + ".events",
		commands:        make(chan string),
		events:          make(chan string),
	}

	return pc, pc.start()
}

func (pc *pipeClient) start() error {
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

	go pc.commandLoop()
	go pc.eventLoop()

	return nil
}

func (pc *pipeClient) commandLoop() {
	log.Println("Starting command loop...")

	defer os.Remove(pc.commandPipeName)

	for {
		// read next command from pipeline
		cmdText := pc.read()

		if cmdText == "" {
			log.Println("Disconnected from command pipe")
			return
		}

		pc.commands <- cmdText
	}
}

func (pc *pipeClient) read() string {
	var bytesRead int
	var err error
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

			//fmt.Printf("read: %d\n", bytesRead)
		}
		input.Close()
		return string(result)
	}
	return ""
}

func (pc *pipeClient) writeResult(result string) {
	log.Println("Waiting for result to consume...")
	output, err := openFifo(pc.commandPipeName, os.O_WRONLY)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Write result:", result)

	output.WriteString(fmt.Sprintf("%s\n", result))
	output.Close()
}

func (pc *pipeClient) emitEvent(evt string) {
	select {
	case pc.events <- evt:
		// Event sent to queue
	default:
		// No event listeners
	}
}

func (pc *pipeClient) eventLoop() {

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

func (pc *pipeClient) close() {
	log.Println("Closing Unix pipe...")

	// TODO: delete temp files
}

func createFifo(filename string) (err error) {
	err = syscall.Mkfifo(filename, 0660)
	return
}

func openFifo(path string, oflag int) (f *os.File, err error) {
	f, err = os.OpenFile(path, oflag, os.ModeNamedPipe)
	return
}
