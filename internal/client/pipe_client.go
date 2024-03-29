//go:build !windows

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/page"
	"github.com/pglet/pglet/internal/page/command"
	"github.com/pglet/pglet/internal/utils"
)

const (
	pageChangeEventPrefix = "page change"
	pageCloseEvent        = "page close "
	readsize              = 64 << 10
)

type PipeClient struct {
	id            string
	pageName      string
	sessionID     string
	pipe          pipe
	hostClient    *HostClient
	emitAllEvents bool
	done          chan bool
	commandBatch  []*command.Command
}

func NewPipeClient(pageName string, sessionID string, hc *HostClient, uds bool, emitAllEvents bool, tickerDuration int) (*PipeClient, error) {
	id, _ := utils.GenerateRandomString(10)

	var err error
	var p pipe
	if uds && runtime.GOOS != "windows" {
		p, err = newUnixDomainSocket(id)
	} else {
		p, err = newNamedPipe(id)
	}

	if err != nil {
		return nil, err
	}

	pc := &PipeClient{
		id:            id,
		pageName:      pageName,
		sessionID:     sessionID,
		pipe:          p,
		hostClient:    hc,
		emitAllEvents: emitAllEvents,
		done:          make(chan bool, 1),
	}

	if tickerDuration > 0 {
		go pc.timerTicker(tickerDuration)
	}

	return pc, nil
}

func (pc *PipeClient) timerTicker(tickerDuration int) {
	ticker := time.NewTicker(time.Duration(tickerDuration) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-pc.done:
			return
		case <-ticker.C:
			pc.emitEvent("timer tick")
		}
	}

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
		cmd, err := command.Parse(strings.Trim(cmdText, "\n"), true)
		if err != nil {
			log.Errorln(err)
			pc.pipe.writeResult(fmt.Sprintf("error %s", err))
			continue
		}

		log.Debugf("Send command: %+v", cmd)

		if cmd.Name == command.Close {
			log.Debugln("Close command")
			pc.emitEvent(pageCloseEvent)
			return
		} else if cmd.Name == command.Begin {
			// start new batch
			log.Debugln("Start new batch")
			pc.commandBatch = make([]*command.Command, 0)
			pc.writeResult("")
		} else if cmd.Name != command.End && pc.commandBatch != nil {
			// add command to batch
			log.Debugln("Add command to the batch")
			pc.commandBatch = append(pc.commandBatch, cmd)
			pc.writeResult("")
		} else if cmd.Name == command.End && pc.commandBatch != nil {
			// run batch
			log.Debugln("Run batch")
			rawResult := pc.hostClient.Call(context.Background(), page.PageCommandsBatchFromHostAction,
				&page.PageCommandsBatchRequestPayload{
					PageName:  pc.pageName,
					SessionID: pc.sessionID,
					Commands:  pc.commandBatch,
				})

			// parse response
			payload := &page.PageCommandsBatchResponsePayload{}
			err = json.Unmarshal(*rawResult, payload)

			if err != nil {
				log.Errorln("Error parsing response from PageCommandsBatchFromHostAction:", err)
			}

			log.Debugln("Response from PageCommandsBatchFromHostAction", payload.Results)

			// save command results
			if payload.Error != "" {
				pc.pipe.writeResult(fmt.Sprintf("error %s", payload.Error))
			} else {
				pc.writeResult(strings.Join(payload.Results, "\n"))
			}
			pc.commandBatch = nil
		} else {
			// run single command
			if cmd.ShouldReturn() {
				// call and wait for result
				rawResult := pc.hostClient.Call(context.Background(), page.PageCommandFromHostAction, &page.PageCommandRequestPayload{
					PageName:  pc.pageName,
					SessionID: pc.sessionID,
					Command:   cmd,
				})

				// parse response
				payload := &page.PageCommandResponsePayload{}
				err = json.Unmarshal(*rawResult, payload)

				if err != nil {
					log.Errorln("Error parsing response from PageCommandFromHostAction:", err)
				}

				// save command results
				if payload.Error != "" {
					pc.pipe.writeResult(fmt.Sprintf("error %s", payload.Error))
				} else {
					pc.writeResult(payload.Result)
				}

			} else {
				// fire and forget
				pc.hostClient.CallAndForget(page.PageCommandFromHostAction, &page.PageCommandRequestPayload{
					PageName:  pc.pageName,
					SessionID: pc.sessionID,
					Command:   cmd,
				})
			}

			if cmd.Name == command.Error {
				log.Debugln("Error command")
				pc.emitEvent(pageCloseEvent)
				return
			}
		}
	}
}

func (pc *PipeClient) writeResult(result string) {
	linesCount := utils.CountRune(result, '\n')
	pc.pipe.writeResult(fmt.Sprintf("%d %s", linesCount, result))
}

func (pc *PipeClient) emitEvent(evt string) {

	if strings.HasPrefix(evt, pageChangeEventPrefix) && !pc.emitAllEvents {
		return
	}

	pc.pipe.emitEvent(evt)

	if strings.HasPrefix(evt, pageCloseEvent) {
		pc.close()
	}
}

func (pc *PipeClient) close() {
	log.Println("Closing pipe client...")
	pc.done <- true
	pc.pipe.close()

	pc.hostClient.UnregisterPipeClient(pc)
}
