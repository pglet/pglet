package client

type protocol interface {
	getCommandPipeName() string
	nextCommand() string
	writeResult(result string)
	emitEvent(evt string)
	close()
}
