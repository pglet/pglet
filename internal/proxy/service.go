package proxy

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alexflint/go-filemutex"
	"github.com/pglet/pglet/internal/client"
	"github.com/pglet/pglet/internal/page"
)

var (
	sockAddr     string
	lockFilename string
)

func init() {
	sockAddr = filepath.Join(os.TempDir(), "pglet.sock")
	lockFilename = filepath.Join(os.TempDir(), "pglet.lock")
}

// Service manages connections to a shared page or app.
type Service struct {
	hcMutex     sync.RWMutex
	hostClients map[string]*client.HostClient
}

func newService() *Service {
	ps := &Service{}
	ps.hostClients = make(map[string]*client.HostClient)
	return ps
}

func (ps *Service) getHostClient(pageURI string) *client.HostClient {
	ps.hcMutex.Lock()
	defer ps.hcMutex.Unlock()

	wsURL := buildWSEndPointURL(pageURI)

	hc, ok := ps.hostClients[wsURL]
	if !ok {
		hc = client.NewHostClient(wsURL)
		err := hc.Start()

		if err != nil {
			log.Fatalf("Cannot connect to %s: %v\n", wsURL, err)
		}
		ps.hostClients[wsURL] = hc
	}
	return hc
}

// ConnectSharedPage establishes a new connection to the specified shared page and returns file name of control pipe.
func (ps *Service) ConnectSharedPage(pageURI *string, pipeName *string) error {

	hc := ps.getHostClient(*pageURI)
	pageName := getPageNameFromURI(*pageURI)

	log.Println("Connecting to shared page:", pageName)

	// call server
	result := hc.Call(page.RegisterHostClientAction, &page.RegisterHostClientRequestPayload{
		PageName: pageName,
		IsApp:    false,
	})

	// parse response
	payload := &page.RegisterHostClientResponsePayload{}
	err := json.Unmarshal(*result, payload)

	if err != nil {
		log.Fatalln("Error calling ConnectSharedPage:", err)
	}

	// create new pipeClient
	pc, err := client.NewPipeClient(pageName, payload.SessionID, hc)
	if err != nil {
		return err
	}

	pc.Start()

	// register pipe client, so it can receive events from pages/sessions
	hc.RegisterPipeClient(pc)

	*pipeName = pc.CommandPipeName()

	return nil
}

// ConnectAppPage waits for new web clients connecting specified page, creates a new session and returns file name of control pipe.
func (ps *Service) ConnectAppPage(pageURI *string, pipeName *string) error {

	hc := ps.getHostClient(*pageURI)
	pageName := getPageNameFromURI(*pageURI)

	log.Println("Connecting to app page:", pageName)

	// call server
	result := hc.Call(page.RegisterHostClientAction, &page.RegisterHostClientRequestPayload{
		PageName: pageName,
		IsApp:    true,
	})

	// parse response
	payload := &page.RegisterHostClientResponsePayload{}
	err := json.Unmarshal(*result, payload)

	if err != nil {
		log.Fatalln("Error calling ConnectSharedPage:", err)
	}

	// wait for new session
	sessionID := <-hc.PageNewSessions(pageName)

	// create new pipeClient
	pc, err := client.NewPipeClient(pageName, sessionID, hc)
	if err != nil {
		return err
	}

	pc.Start()

	// register pipe client, so it can receive events from pages/sessions
	hc.RegisterPipeClient(pc)

	*pipeName = pc.CommandPipeName()

	return nil
}

func RunService() {

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

	proxySvc := newService()
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
		log.Fatalln("Cannot parse page URI:", err)
	}

	u.Scheme = "ws"
	u.Path = "ws"
	u.RawQuery = ""

	return u.String()
}

func getPageNameFromURI(pageURI string) string {
	u, err := url.Parse(pageURI)
	if err != nil {
		log.Fatalln("Cannot parse page URI:", err)
	}

	return strings.ToLower(strings.Trim(u.Path, "/"))
}
