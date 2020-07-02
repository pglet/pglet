package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

const (
	apiRoutePrefix      string = "/api"
	contentRootFolder   string = "./client/build"
	siteDefaultDocument string = "index.html"
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
		pipeName, _ := client.connectAppPage(appName)
		fmt.Println(pipeName)
	}
}

func runClient2() {
	fmt.Printf("Running in client mode: %s...\n", pageName)
	u := url.URL{Scheme: "ws", Host: *&serverAddr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		defer func() {
			fmt.Println("Closing...")
		}()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from Go: %s", pageName)))

	time.Sleep(5 * time.Second)

	// run proxy
	execPath, _ := os.Executable()
	fmt.Println(execPath)

	cmd := exec.Command(execPath, "--session-id=12345")
	err = cmd.Start()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(cmd.Process.Pid)
}
