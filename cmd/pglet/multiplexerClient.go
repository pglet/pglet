package main

import (
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

type multiplexerClient struct {
	client *rpc.Client
}

func start(proxy *multiplexerClient) {
	log.Println("Connecting to multiplexer service")
	var err error

	for {
		proxy.client, err = rpc.DialHTTP("unix", sockAddr)
		if err != nil {
			log.Println("Error connecting to multiplexer service:", err)

			// start multiplexer service
			startMultiplexerService()
			time.Sleep(1 * time.Second)
		} else {
			log.Println("Connected to multiplexer service")
			break
		}
	}
}

func (proxy *multiplexerClient) connectSharedPage(pageName string) (pipeFilename string, err error) {
	err = proxy.client.Call("MultiplexerService.ConnectSharedPage", &pageName, &pipeFilename)
	return
}

func (proxy *multiplexerClient) connectAppPage(pageName string) (pipeFilename string, err error) {
	err = proxy.client.Call("MultiplexerService.ConnectAppPage", &pageName, &pipeFilename)
	return
}

func startMultiplexerService() {
	log.Println("Starting multiplexer service")
	// run proxy
	execPath, _ := os.Executable()

	cmd := exec.Command(execPath, "--multiplexer")
	err := cmd.Start()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Multiplexer service process started with PID:", cmd.Process.Pid)
}
