package main

type pipeClient struct {
	id              string
	pageName        string
	commandPipeName string
	eventPipeName   string
	events          chan string
	hostClient      *hostClient
}
