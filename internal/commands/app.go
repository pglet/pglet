package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/spf13/cobra"
)

func newAppCommand() *cobra.Command {

	var uds bool

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
				pipeName, err := client.ConnectAppPage(cmd.Context(), &proxy.ConnectPageArgs{
					PageURI: args[0],
					Uds:     uds,
				})
				if err != nil {
					log.Fatalln("Connect app error:", err)
				}
				fmt.Println(pipeName)
			}
		},
	}

	cmd.Flags().BoolVarP(&uds, "uds", "", false, "force Unix domain sockets to connect from PowerShell on Linux/macOS")

	return cmd
}
