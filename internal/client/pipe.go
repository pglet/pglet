//go:build !windows

package client

type pipe interface {
	getCommandPipeName() string
	nextCommand() string
	writeResult(result string)
	emitEvent(evt string)
	close()
}
