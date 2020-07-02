package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/alexflint/go-filemutex"
)

const sockAddr = "/tmp/pglet.sock"
const lockFilename = "/tmp/pglet.lock"

// ProxyService manages connections to a shared page or app.
type ProxyService struct {
}

// ConnectSharedPage establishes a new connection to the specified shared page and returns file name of control pipe.
func (g ProxyService) ConnectSharedPage(name *string, pipeName *string) error {
	*pipeName = fmt.Sprintf("/tmp/%s-aaa", *name)
	return nil
}

// ConnectAppPage waits for new web clients connecting specified page, creates a new session and returns file name of control pipe.
func (g ProxyService) ConnectAppPage(name *string, pipeName *string) error {
	*pipeName = fmt.Sprintf("/tmp/%s-bbb", *name)
	return nil
}

func runProxyService() {

	log.Println("Starting Proxy service...")

	m, err := filemutex.New(lockFilename)
	if err != nil {
		log.Fatalln("Directory did not exist or file could not created")
	}

	err = m.TryLock()
	if err != nil {
		fmt.Println("Another Proxy service process has started")
		os.Exit(1)
	}

	defer m.Unlock()

	if err := os.RemoveAll(sockAddr); err != nil {
		log.Fatal(err)
	}

	greeter := new(ProxyService)
	rpc.Register(greeter)
	rpc.HandleHTTP()
	l, e := net.Listen("unix", sockAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("Waiting for connections...")
	err = http.Serve(l, nil)
	fmt.Println(err)
}
