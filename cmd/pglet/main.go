package main

import (
	"flag"
	"fmt"
	"os"
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

func init() {
	flag.StringVar(&pageName, "page", "", "Shared page name to create and connect to.")
	flag.StringVar(&appName, "app", "", "App page name to create and connect to.")
	flag.StringVar(&serverAddr, "server", "", "Pglet server address.")
	flag.StringVar(&sessionID, "session-id", "", "Client session ID.")
	flag.BoolVar(&isProxy, "proxy", false, "Start Proxy service.")
	flag.IntVar(&serverPort, "port", 5000, "The port number to run pglet server on.")
	flag.Parse()

	if !isProxy && pageName == "" && appName == "" {
		isServer = true

		if serverPort < 0 || serverPort > 65535 {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}
}

func main() {

	if isProxy {
		runProxyService()
	} else if isServer {
		runServer()
	} else {
		runClient()
	}
}

func runClient() {
	client := &proxyClient{}
	client.start()

	if pageName != "" {
		pipeName, _ := client.connectSharedPage(pageName)
		fmt.Println(pipeName)
	} else if appName != "" {
		// continuously wait for new client connections
		for {
			pipeName, _ := client.connectAppPage(appName)
			fmt.Println(pipeName)
		}
	}
}
