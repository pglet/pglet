package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/pglet/pglet/internal/utils"
	"github.com/spf13/cobra"
)

func newPageCommand() *cobra.Command {

	var public bool
	var private bool
	var server string
	var token string
	var uds bool

	var cmd = &cobra.Command{
		Use:   "page [[namespace/]<page_name>]",
		Short: "Connect to a shared page",
		Long:  `Page command is ...`,
		Run: func(cmd *cobra.Command, args []string) {
			client := &proxy.Client{}
			client.Start()

			pageName := "*" // auto-generated
			if len(args) > 0 {
				pageName = args[0]
			}

			results, err := client.ConnectSharedPage(cmd.Context(), &proxy.ConnectPageArgs{
				PageName: pageName,
				Private:  private,
				Public:   public,
				Server:   server,
				Token:    token,
				Uds:      uds,
			})

			if err != nil {
				log.Fatalln("Connect page error:", err)
			}

			utils.OpenBrowser(results.PageURL)

			// output connection ID and page URL to be consumed by a client
			fmt.Println(results.PipeName, results.PageURL)
		},
	}

	cmd.Flags().BoolVarP(&public, "public", "", false, "makes the page available as public at pglet.io service or a self-hosted Pglet server")
	cmd.Flags().BoolVarP(&private, "private", "", false, "makes the page available as private at pglet.io service or a self-hosted Pglet server")
	cmd.Flags().StringVarP(&server, "server", "s", "", "connects to the page on a self-hosted Pglet server")
	cmd.Flags().StringVarP(&token, "token", "t", "", "authentication token for pglet.io service or a self-hosted Pglet server")
	cmd.Flags().BoolVarP(&uds, "uds", "", false, "force Unix domain sockets to connect from PowerShell on Linux/macOS")

	return cmd
}
