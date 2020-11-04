package commands

import (
	"sync"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/pglet/pglet/internal/server"
	"github.com/spf13/cobra"
)

func newServerCommand() *cobra.Command {

	var serverPort int

	var cmd = &cobra.Command{
		Use:   "server",
		Short: "Start server service",
		Long:  `Server is for ...`,
		Run: func(cmd *cobra.Command, args []string) {

			waitGroup := sync.WaitGroup{}

			waitGroup.Add(2)
			go proxy.Start(cmd.Context(), &waitGroup)
			go server.Start(cmd.Context(), &waitGroup, serverPort)

			waitGroup.Wait()
		},
	}

	cmd.Flags().IntVarP(&serverPort, "port", "p", 5000, "port on which the server will listen")

	return cmd
}
