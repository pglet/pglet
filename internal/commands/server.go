package commands

import (
	"sync"

	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/config"
	"github.com/pglet/pglet/internal/server"
	"github.com/spf13/cobra"
)

var (
	defaultPort int = 5000
)

func NewServerCommand() *cobra.Command {

	var serverPort int

	var cmd = &cobra.Command{
		Use:     "pglet-server",
		Short:   "Start Pglet server service",
		Long:    `Pglet Server is ...`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configLogging()
		},
		Run: func(cmd *cobra.Command, args []string) {

			// init cache
			cache.Init()

			waitGroup := sync.WaitGroup{}

			waitGroup.Add(1)
			go server.Start(cmd.Context(), &waitGroup, serverPort)
			waitGroup.Wait()
		},
	}

	cmd.SetVersionTemplate("{{.Version}}")

	cmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "verbosity level for logs")

	cmd.Flags().IntVarP(&serverPort, "port", "p", config.ServerPort(), "port on which the server will listen")

	return cmd
}
