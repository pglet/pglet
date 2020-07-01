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

// MultiplexerService manages connections to a shared page or app.
type MultiplexerService struct {
}

// ConnectSharedPage establishes a new connection to the specified shared page and returns file name of control pipe.
func (g MultiplexerService) ConnectSharedPage(name *string, pipeName *string) error {
	*pipeName = fmt.Sprintf("/tmp/%s-aaa", *name)
	return nil
}

// ConnectAppPage waits for new web clients connecting specified page, creates a new session and returns file name of control pipe.
func (g MultiplexerService) ConnectAppPage(name *string, pipeName *string) error {
	*pipeName = fmt.Sprintf("/tmp/%s-bbb", *name)
	return nil
}

func runProxyServer() {

	m, err := filemutex.New("/tmp/foo.lock")
	if err != nil {
		log.Fatalln("Directory did not exist or file could not created")
	}

	err = m.TryLock()
	if err != nil {
		fmt.Println("Another server process has started")
		os.Exit(1)
	}

	defer m.Unlock()

	//time.Sleep(30 * time.Second)

	if err := os.RemoveAll(sockAddr); err != nil {
		log.Fatal(err)
	}

	greeter := new(MultiplexerService)
	rpc.Register(greeter)
	rpc.HandleHTTP()
	l, e := net.Listen("unix", sockAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("Serving...")
	err = http.Serve(l, nil)
	fmt.Println(err)
}

func runProxyClient() {
	client, err := rpc.DialHTTP("unix", "/tmp/rpc.sock")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	name := "Joe"
	var reply string
	call := client.Go("MultiplexerService.Greet", &name, &reply, nil)
	replyCall := <-call.Done

	if replyCall.Error != nil {
		log.Fatal("MultiplexerService error:", replyCall.Error)
	}
	fmt.Printf("Got '%s'\n", reply)
}
