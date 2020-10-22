package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pglet/pglet/internal/commands"

	"github.com/pglet/pglet/internal/proxy"

	"github.com/pglet/pglet/internal/server"
)

var (
	isServer   bool
	isProxy    bool
	serverPort int
	serverAddr string
	pageName   string
	appName    string
	sessionID  string
)

func main() {

	commands.PrintVersion()

	flag.StringVar(&pageName, "page", "", "Shared page name to create and connect to.")
	flag.StringVar(&appName, "app", "", "App page name to create and connect to.")
	flag.StringVar(&serverAddr, "server", "", "Pglet server address.")
	flag.StringVar(&sessionID, "session-id", "", "Client session ID.")
	flag.BoolVar(&isProxy, "proxy", false, "Start Proxy service.")
	flag.IntVar(&serverPort, "port", server.DefaultServerPort, "The port number to run pglet server on.")
	flag.Parse()

	if !isProxy && pageName == "" && appName == "" {
		isServer = true

		if serverPort < 0 || serverPort > 65535 {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if isProxy {
		proxy.RunService()
	} else if isServer {
		server.Start(serverPort)
	} else {
		runClient()
	}
}

func runClient() {
	client := &proxy.Client{}
	client.Start()

	if pageName != "" {
		pipeName, err := client.ConnectSharedPage(pageName)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(pipeName)
	} else if appName != "" {
		// continuously wait for new client connections
		for {
			pipeName, err := client.ConnectAppPage(appName)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(pipeName)
		}
	}
}
