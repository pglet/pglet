package proxy

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/keegancsmith/rpc"
	log "github.com/sirupsen/logrus"
)

const (
	connectAttempts = 20
)

type Client struct {
	client *rpc.Client
}

func (proxy *Client) Start() {
	var err error

	for i := 1; i <= connectAttempts; i++ {
		//log.Printf("Connecting to Proxy service (attempt %d of %d)\n", i, connectAttempts)
		proxy.client, err = rpc.DialHTTP("unix", sockAddr)
		if err != nil {
			//log.Println("Error connecting to Proxy service:", err)

			// start Proxy service
			startProxyService()
			time.Sleep(1 * time.Second)
		} else {
			//log.Println("Connected to Proxy service")
			return
		}
	}

	log.Fatalf("Gave up connecting to Proxy service after %d attemps\n", connectAttempts)
}

func (proxy *Client) ConnectSharedPage(ctx context.Context, pageName string) (pipeFilename string, err error) {
	err = proxy.client.Call(ctx, "Service.ConnectSharedPage", &pageName, &pipeFilename)
	return
}

func (proxy *Client) ConnectAppPage(ctx context.Context, pageName string) (pipeFilename string, err error) {
	err = proxy.client.Call(ctx, "Service.ConnectAppPage", &pageName, &pipeFilename)
	return
}

func startProxyService() {
	log.Println("Starting Proxy service")
	// run proxy
	execPath, _ := os.Executable()

	cmd := exec.Command(execPath, "server")
	err := cmd.Start()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Proxy service process started with PID:", cmd.Process.Pid)
}
