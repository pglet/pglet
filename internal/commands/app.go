package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/spf13/cobra"
)

func newAppCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "app",
		Short: "Connect to an app",
		Long:  `App command is ...`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := &proxy.Client{}
			client.Start()

			// continuously wait for new client connections
			for {
				pipeName, err := client.ConnectAppPage(args[0])
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(pipeName)
			}
		},
	}

	return cmd
}
