package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/url"
	"os"
	"sync"

	"github.com/alexflint/go-filemutex"
	"github.com/pglet/pglet/page"
)

const (
	sockAddr     = "/tmp/pglet.sock"
	lockFilename = "/tmp/pglet.lock"
)

// ProxyService manages connections to a shared page or app.
type ProxyService struct {
	hcMutex     sync.RWMutex
	hostClients map[string]*hostClient
}

func newProxyService() *ProxyService {
	ps := &ProxyService{}
	ps.hostClients = make(map[string]*hostClient)
	return ps
}

func (ps *ProxyService) getHostClient(pageURI string) *hostClient {
	ps.hcMutex.Lock()
	defer ps.hcMutex.Unlock()

	wsURL := buildWSEndPointURL(pageURI)

	hc, ok := ps.hostClients[wsURL]
	if !ok {
		hc = newHostClient(wsURL)
		err := hc.start()

		if err != nil {
			log.Fatalf("Cannot connect to %s: %v\n", wsURL, err)
		}
		ps.hostClients[wsURL] = hc
	}
	return hc
}

// ConnectSharedPage establishes a new connection to the specified shared page and returns file name of control pipe.
func (ps *ProxyService) ConnectSharedPage(pageURI *string, pipeName *string) error {

	hc := ps.getHostClient(*pageURI)

	// call server
	result := hc.call(page.RegisterHostClientAction, &page.RegisterClientActionRequestPayload{
		PageName: *pageURI,
	})

	// parse response
	payload := &page.RegisterClientActionResponsePayload{}
	err := json.Unmarshal(*result, payload)

	if err != nil {
		log.Fatalln("Error calling ConnectSharedPage:", err)
	}

	*pipeName = fmt.Sprintf("/tmp/%s-aaa", *pageURI)
	return nil
}

// ConnectAppPage waits for new web clients connecting specified page, creates a new session and returns file name of control pipe.
func (ps *ProxyService) ConnectAppPage(pageURI *string, pipeName *string) error {

	hc := ps.getHostClient(*pageURI)
	hc.send <- []byte("hello!")

	*pipeName = fmt.Sprintf("%s", *pageURI)
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
		log.Println("Another Proxy service process has started")
		os.Exit(1)
	}

	defer m.Unlock()

	if err := os.RemoveAll(sockAddr); err != nil {
		log.Fatal(err)
	}

	proxySvc := newProxyService()
	rpc.Register(proxySvc)
	rpc.HandleHTTP()
	l, e := net.Listen("unix", sockAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Println("Waiting for connections...")
	err = http.Serve(l, nil)
	log.Println(err)
}

func buildWSEndPointURL(pageURI string) string {
	u, err := url.Parse(pageURI)
	if err != nil {
		log.Fatalln("Cannot parse page URL:", err)
	}

	u.Scheme = "ws"
	u.Path = "ws"
	u.RawQuery = ""

	return u.String()
}
