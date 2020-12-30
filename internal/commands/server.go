package commands

import (
	"os"
	"strconv"
	"sync"

	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/proxy"
	"github.com/pglet/pglet/internal/server"
	"github.com/spf13/cobra"
)

var (
	defaultPort int = 5000
)

func newServerCommand() *cobra.Command {

	var serverPort int

	envPort, err := strconv.Atoi(os.Getenv("PGLET_SERVER_PORT"))
	if err == nil {
		defaultPort = envPort
	}

	var cmd = &cobra.Command{
		Use:   "server",
		Short: "Start server service",
		Long:  `Server is for ...`,
		Run: func(cmd *cobra.Command, args []string) {

			// init cache
			cache.Init()

			waitGroup := sync.WaitGroup{}

			waitGroup.Add(2)
			go server.Start(cmd.Context(), &waitGroup, serverPort)
			go proxy.Start(cmd.Context(), &waitGroup)

			waitGroup.Wait()
		},
	}

	cmd.Flags().IntVarP(&serverPort, "port", "p", defaultPort, "port on which the server will listen")

	return cmd
}
