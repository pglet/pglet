package proxy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/keegancsmith/rpc"
	"github.com/pglet/pglet/internal/config"
	log "github.com/sirupsen/logrus"
)

const (
	connectAttempts = 20
)

var (
	sockAddr string
)

func init() {
	sockAddr = filepath.Join(os.TempDir(), "pglet.sock")
}

type Client struct {
	client *rpc.Client
}

type ConnectPageArgs struct {
	PageName       string
	Web            bool
	Server         string
	Token          string
	Permissions    string
	Uds            bool
	EmitAllEvents  bool
	TickerDuration int
}

type ConnectPageResults struct {
	PipeName string
	PageName string
	PageURL  string
}

func (proxy *Client) Start() {
	var err error

	for i := 1; i <= connectAttempts; i++ {
		proxy.client, err = rpc.DialHTTP("unix", sockAddr)
		if err != nil {
			// start Proxy service
			startProxyService()
			time.Sleep(200 * time.Millisecond)
		} else {
			return
		}
	}

	log.Fatalf("Gave up connecting to Client service after %d attemps\n", connectAttempts)
}

func (proxy *Client) ConnectSharedPage(ctx context.Context, args *ConnectPageArgs) (results *ConnectPageResults, err error) {
	err = proxy.client.Call(ctx, "Service.ConnectSharedPage", &args, &results)
	return
}

func (proxy *Client) ConnectAppPage(ctx context.Context, args *ConnectPageArgs) (results *ConnectPageResults, err error) {
	err = proxy.client.Call(ctx, "Service.ConnectAppPage", &args, &results)
	return
}

func (proxy *Client) WaitAppSession(ctx context.Context, args *ConnectPageArgs) (results *ConnectPageResults, err error) {
	err = proxy.client.Call(ctx, "Service.WaitAppSession", &args, &results)
	return
}

func startProxyService() {
	log.Traceln("Starting Pglet server")

	// run server
	execPath, _ := os.Executable()

	cmd := GetDetachedCmd(execPath, "client")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=true", config.LogToFileFlag))

	err := cmd.Start()

	if err != nil {
		log.Fatalln(err)
	}

	log.Traceln("Server process started with PID:", cmd.Process.Pid)
}
