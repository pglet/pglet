package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/spf13/cobra"
)

func newPageCommand() *cobra.Command {

	var uds bool

	var cmd = &cobra.Command{
		Use:   "page <namespace/page>",
		Short: "Connect to a shared page",
		Long:  `Page command is ...`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := &proxy.Client{}
			client.Start()

			pipeName, err := client.ConnectSharedPage(cmd.Context(), &proxy.ConnectPageArgs{
				PageURI: args[0],
				Uds:     uds,
			})
			if err != nil {
				log.Fatalln("Connect page error:", err)
			}
			fmt.Println(pipeName)
		},
	}

	cmd.Flags().BoolVarP(&uds, "uds", "", false, "force Unix domain sockets to connect from PowerShell on Linux/macOS")

	return cmd
}
