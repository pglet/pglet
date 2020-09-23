package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pglet/pglet/page"
	"github.com/pglet/pglet/utils"
)

type pipeClient struct {
	id         string
	pageName   string
	sessionID  string
	pipe       *pipeImpl
	hostClient *hostClient
}

func newPipeClient(pageName string, sessionID string, hc *hostClient) (*pipeClient, error) {
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

func (pc *pipeClient) commandPipeName() string {
	return pc.pipe.commandPipeName
}

func (pc *pipeClient) start() error {

	go pc.commandLoop()

	return nil
}

func (pc *pipeClient) commandLoop() {
	log.Println("Starting command loop...")

	for {
		// read next command from pipeline
		cmdText := pc.pipe.read()

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

		rawResult := pc.hostClient.call(page.PageCommandFromHostAction, &page.PageCommandRequestPayload{
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

		pc.pipe.writeResult(result)
	}
}

func (pc *pipeClient) emitEvent(evt string) {
	pc.pipe.emitEvent(evt)
}

func (pc *pipeClient) close() {
	log.Println("Closing pipe client...")

	pc.pipe.close()

	pc.hostClient.unregisterPipeClient(pc)
}
