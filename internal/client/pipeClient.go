package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pglet/pglet/internal/page"
	"github.com/pglet/pglet/internal/utils"
)

type PipeClient struct {
	id         string
	pageName   string
	sessionID  string
	pipe       protocol
	hostClient *HostClient
}

func NewPipeClient(pageName string, sessionID string, hc *HostClient) (*PipeClient, error) {
	id, _ := utils.GenerateRandomString(10)

	pipe, err := newNamedPipes(id)

	if err != nil {
		return nil, err
	}

	pc := &PipeClient{
		id:         id,
		pageName:   pageName,
		sessionID:  sessionID,
		pipe:       pipe,
		hostClient: hc,
	}

	return pc, nil
}

func (pc *PipeClient) CommandPipeName() string {
	return pc.pipe.getCommandPipeName()
}

func (pc *PipeClient) Start() error {

	go pc.commandLoop()

	return nil
}

func (pc *PipeClient) commandLoop() {
	log.Println("Starting command loop...")

	for {
		// read next command from pipeline
		cmdText := pc.pipe.nextCommand()

		// parse command
		command, err := page.ParseCommand(cmdText)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Send command: %+v", command)

		if command.Name == page.Quit {
			pc.close()
			return
		}

		rawResult := pc.hostClient.Call(page.PageCommandFromHostAction, &page.PageCommandRequestPayload{
			PageName:  pc.pageName,
			SessionID: pc.sessionID,
			Command:   *command,
		})

		// parse response
		payload := &page.PageCommandResponsePayload{}
		err = json.Unmarshal(*rawResult, payload)

		if err != nil {
			log.Fatalln("Error parsing response from PageCommandFromHostAction:", err)
		}

		// save command results
		result := payload.Result
		if payload.Error != "" {
			result = fmt.Sprintf("error %s", payload.Error)
		}

		pc.pipe.writeResult("aaa" + result)
	}
}

func (pc *PipeClient) emitEvent(evt string) {
	pc.pipe.emitEvent(evt)
}

func (pc *PipeClient) close() {
	log.Println("Closing pipe client...")

	pc.pipe.close()

	pc.hostClient.UnregisterPipeClient(pc)
}
