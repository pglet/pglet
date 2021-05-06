package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/proxy"
	"github.com/pglet/pglet/internal/utils"
	"github.com/spf13/cobra"
)

func newPageCommand() *cobra.Command {

	var web bool
	var server string
	var token string
	var permissions string
	var uds bool
	var tickerDuration int
	var noWindow bool
	var allEvents bool
	var window string

	var cmd = &cobra.Command{
		Use:   "page [[namespace/]<page_name>]",
		Short: "Connect to a shared page",
		Long:  `Page command creates a new shared page and opens connection to it.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := &proxy.Client{}
			client.Start()

			pageName := "*" // auto-generated
			if len(args) > 0 {
				pageName = args[0]
			}

			results, err := client.ConnectSharedPage(cmd.Context(), &proxy.ConnectPageArgs{
				PageName:       pageName,
				Web:            web,
				Server:         server,
				Token:          token,
				Permissions:    permissions,
				Uds:            uds,
				EmitAllEvents:  allEvents,
				TickerDuration: tickerDuration,
			})

			if err != nil {
				log.Fatalln("Connect page error:", err)
			}

			if !noWindow {
				utils.OpenBrowser(results.PageURL, window)
			}

			// output connection ID and page URL to be consumed by a client
			fmt.Println(results.PipeName, results.PageURL)
		},
	}

	cmd.Flags().BoolVarP(&web, "web", "", false, "makes the page available as public at pglet.io service or a self-hosted Pglet server")
	cmd.Flags().StringVarP(&server, "server", "s", "", "connects to the page on a self-hosted Pglet server")
	cmd.Flags().StringVarP(&token, "token", "t", "", "authentication token for pglet.io service or a self-hosted Pglet server")
	cmd.Flags().StringVarP(&permissions, "permissions", "", "", "comma-separated list of users and groups allowed to access this app")
	cmd.Flags().BoolVarP(&uds, "uds", "", false, "force Unix domain sockets to connect from PowerShell on Linux/macOS")
	cmd.Flags().IntVarP(&tickerDuration, "ticker", "", 0, "interval in milliseconds between 'tick' events; disabled if not specified.")
	cmd.Flags().StringVarP(&window, "window", "", "", "open page in a window with specified dimensions and position: [x,y,]width,height")
	cmd.Flags().BoolVarP(&noWindow, "no-window", "", false, "do not open browser window")
	cmd.Flags().BoolVarP(&allEvents, "all-events", "", false, "emit all page events")

	return cmd
}
