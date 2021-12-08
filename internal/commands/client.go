//go:build !windows

package commands

import (
	"sync"

	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/config"
	"github.com/pglet/pglet/internal/proxy"
	"github.com/pglet/pglet/internal/server"
	"github.com/spf13/cobra"
)

func newClientCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "client",
		Short: "Start client service",
		Long:  `Client service is for ...`,
		Run: func(cmd *cobra.Command, args []string) {

			// ensure one executable instance is running
			m, err := filemutex.New(lockFilename)
			if err != nil {
				log.Fatalln("Cannot create mutex - directory did not exist or file could not be created")
			}
		
			err = m.TryLock()
			if err != nil {
				log.Fatalln("Another Pglet Server process has already started")
			}

			defer m.Unlock()

			// init cache
			cache.Init()

			waitGroup := sync.WaitGroup{}
			waitGroup.Add(2)
			go server.Start(cmd.Context(), &waitGroup, serverPort)
			go proxy.Start(cmd.Context(), &waitGroup)
			waitGroup.Wait()
		},
	}

	return cmd
}
