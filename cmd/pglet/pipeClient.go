package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"syscall"

	"github.com/pglet/pglet/page"
	"github.com/pglet/pglet/utils"
)

const (
	readsize = 64 << 10
)

type pipeClient struct {
	id              string
	pageName        string
	sessionID       string
	commandPipeName string
	eventPipeName   string
	events          chan string
	hostClient      *hostClient
	done            chan bool
}

func newPipeClient(pageName string, sessionID string, hc *hostClient) (*pipeClient, error) {
	id, _ := utils.GenerateRandomString(10)
	pipeName := path.Join(os.TempDir(), fmt.Sprintf("pglet_pipe_%s", id))

	pc := &pipeClient{
		id:              id,
		pageName:        pageName,
		sessionID:       sessionID,
		commandPipeName: pipeName,
		eventPipeName:   pipeName + ".events",
		events:          make(chan string),
		hostClient:      hc,
	}

	return pc, nil
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

	go pc.startCommandLoop()
	go pc.startEventLoop()

	return nil
}

func (pc *pipeClient) startCommandLoop() {
	log.Println("Starting command loop...")

	defer os.Remove(pc.commandPipeName)

	for {
		// read next command fro pipeline
		command := pc.read()

		// TODO send command to hostClient
		fmt.Print(command)

		result := pc.hostClient.call(page.PageCommandFromHostAction, &page.PageCommandActionRequestPayload{
			PageName:  pc.pageName,
			SessionID: pc.sessionID,
			Command:   command,
		})

		// parse response
		payload := &page.PageCommandActionResponsePayload{}
		err := json.Unmarshal(*result, payload)

		if err != nil {
			log.Fatalln("Error calling PageCommandFromHostAction:", err)
		}

		// reply back to pipe with command results
		pc.write(fmt.Sprintf("%s\n", payload.Result))
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
	log.Fatal(err)
	return ""
}

func (pc *pipeClient) write(s string) {
	output, err := openFifo(pc.commandPipeName, os.O_WRONLY)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()
	output.Write([]byte(s))
}

func (pc *pipeClient) emitEvent(evt string) {
	select {
	case pc.events <- evt:
		// Event sent to queue
	default:
		// No event listeners
	}
}

func (pc *pipeClient) startEventLoop() {

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

func createFifo(filename string) (err error) {
	err = syscall.Mkfifo(filename, 0660)
	return
}

func openFifo(path string, oflag int) (f *os.File, err error) {
	f, err = os.OpenFile(path, oflag, os.ModeNamedPipe)
	return
}
