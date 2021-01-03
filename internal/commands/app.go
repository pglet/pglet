package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/pglet/pglet/internal/utils"
	"github.com/spf13/cobra"
)

func newAppCommand() *cobra.Command {

	var web bool
	var server string
	var token string
	var uds bool
	var noWindow bool
	var window string

	var cmd = &cobra.Command{
		Use:   "app [[namespace/]<page_name>]",
		Short: "Connect to an app",
		Long:  `App command creates a new multi-user app and waits for new web user connections.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := &proxy.Client{}
			client.Start()

			pageName := "*" // auto-generated
			if len(args) > 0 {
				pageName = args[0]
			}

			connectArgs := &proxy.ConnectPageArgs{
				PageName: pageName,
				Web:      web,
				Server:   server,
				Token:    token,
				Uds:      uds,
			}

			results, err := client.ConnectAppPage(cmd.Context(), connectArgs)
			if err != nil {
				log.Fatalln("Connect app error:", err)
			}

			connectArgs.PageName = results.PageName

			if !noWindow {
				utils.OpenBrowser(results.PageURL, window)
			}

			fmt.Println(results.PageURL)

			// continuously wait for new client connections
			for {
				results, err := client.WaitAppSession(cmd.Context(), connectArgs)
				if err != nil {
					log.Fatalln("Error waiting for a new session:", err)
				}
				fmt.Println(results.PipeName)
			}
		},
	}

	cmd.Flags().BoolVarP(&web, "web", "", false, "makes the app available as public at pglet.io service or a self-hosted Pglet server")
	cmd.Flags().StringVarP(&server, "server", "s", "", "connects to the app on a self-hosted Pglet server")
	cmd.Flags().StringVarP(&token, "token", "t", "", "authentication token for pglet.io service or a self-hosted Pglet server")
	cmd.Flags().BoolVarP(&uds, "uds", "", false, "force Unix domain sockets to connect from PowerShell on Linux/macOS")
	cmd.Flags().StringVarP(&window, "window", "", "", "open app in a window with specified dimensions and position: [x,y,]width,height")
	cmd.Flags().BoolVarP(&noWindow, "no-window", "", false, "do not open browser window")

	return cmd
}
