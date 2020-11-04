package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/pglet/pglet/internal/commands"
)

func main() {

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-interruptCh:
			cancel()
		}
	}()

	if err := commands.NewRootCmd().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
