package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/pglet/pglet/internal/utils"
	"github.com/spf13/cobra"
)

func newAppCommand() *cobra.Command {

	var public bool
	var private bool
	var server string
	var token string
	var uds bool

	var cmd = &cobra.Command{
		Use:   "app [[namespace/]<page_name>]",
		Short: "Connect to an app",
		Long:  `App command is ...`,
		Run: func(cmd *cobra.Command, args []string) {
			client := &proxy.Client{}
			client.Start()

			pageName := "*" // auto-generated
			if len(args) > 0 {
				pageName = args[0]
			}

			connectArgs := &proxy.ConnectPageArgs{
				PageName: pageName,
				Private:  private,
				Public:   public,
				Server:   server,
				Token:    token,
				Uds:      uds,
			}

			results, err := client.ConnectAppPage(cmd.Context(), connectArgs)
			if err != nil {
				log.Fatalln("Connect app error:", err)
			}

			connectArgs.PageName = results.PageName
			utils.OpenBrowser(results.PageURL)

			// continuously wait for new client connections
			for {
				results, err := client.WaitAppSession(cmd.Context(), connectArgs)
				if err != nil {
					log.Fatalln("Error waiting for a new session:", err)
				}
				fmt.Println(results.PipeName, results.PageURL)
			}
		},
	}

	cmd.Flags().BoolVarP(&public, "public", "", false, "makes the app available as public at pglet.io service or a self-hosted Pglet server")
	cmd.Flags().BoolVarP(&private, "private", "", false, "makes the app available as private at pglet.io service or a self-hosted Pglet server")
	cmd.Flags().StringVarP(&server, "server", "s", "", "connects to the app on a self-hosted Pglet server")
	cmd.Flags().StringVarP(&token, "token", "t", "", "authentication token for pglet.io service or a self-hosted Pglet server")
	cmd.Flags().BoolVarP(&uds, "uds", "", false, "force Unix domain sockets to connect from PowerShell on Linux/macOS")

	return cmd
}
