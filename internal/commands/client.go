package commands

import (
	"sync"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/spf13/cobra"
)

func newClientCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "client",
		Short: "Start client service",
		Long:  `Client service is for ...`,
		Run: func(cmd *cobra.Command, args []string) {

			waitGroup := sync.WaitGroup{}
			waitGroup.Add(1)
			go proxy.Start(cmd.Context(), &waitGroup)
			waitGroup.Wait()
		},
	}

	return cmd
}
