package commands

import (
	"sync"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/spf13/cobra"
)

func newProxyCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "proxy",
		Short: "Start proxy service",
		Long:  `Proxy service is for ...`,
		Run: func(cmd *cobra.Command, args []string) {

			waitGroup := sync.WaitGroup{}
			waitGroup.Add(1)
			go proxy.Start(cmd.Context(), &waitGroup)
			waitGroup.Wait()
		},
	}

	return cmd
}
