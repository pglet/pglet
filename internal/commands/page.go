package commands

import (
	"fmt"
	"log"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/spf13/cobra"
)

func newPageCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "page <namespace/page>",
		Short: "Connect to a shared page",
		Long:  `Page command is ...`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := &proxy.Client{}
			client.Start()

			pageName := args[0]
			pipeName, err := client.ConnectSharedPage(pageName)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(pipeName)
		},
	}

	return cmd
}
