package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pglet/pglet/internal/page"
	"github.com/pglet/pglet/internal/utils"
)

type pipeClient struct {
	id         string
	pageName   string
	sessionID  string
	pipe       *pipeImpl
	hostClient *HostClient
}

func NewPipeClient(pageName string, sessionID string, hc *HostClient) (*pipeClient, error) {
	id, _ := utils.GenerateRandomString(10)

	pipe, err := newPipeImpl(id)

	if err != nil {
		return nil, err
	}

	pc := &pipeClient{
		id:         id,
		pageName:   pageName,
		sessionID:  sessionID,
		pipe:       pipe,
		hostClient: hc,
	}

	return pc, nil
}

func (pc *pipeClient) CommandPipeName() string {
	return pc.pipe.commandPipeName
}

func (pc *pipeClient) Start() error {

	go pc.commandLoop()

	return nil
}

func (pc *pipeClient) commandLoop() {
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

func (pc *pipeClient) emitEvent(evt string) {
	pc.pipe.emitEvent(evt)
}

func (pc *pipeClient) close() {
	log.Println("Closing pipe client...")

	pc.pipe.close()

	pc.hostClient.UnregisterPipeClient(pc)
}
