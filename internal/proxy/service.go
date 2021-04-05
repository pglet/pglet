package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/alexflint/go-filemutex"
	"github.com/keegancsmith/rpc"
	"github.com/pglet/pglet/internal/client"
	"github.com/pglet/pglet/internal/page"
	"github.com/pglet/pglet/internal/server"
)

const (
	pgletIoURL = "https://app.pglet.io"
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

func (ps *Service) getHostClient(serverURL string) (*client.HostClient, error) {
	ps.hcMutex.Lock()
	defer ps.hcMutex.Unlock()

	wsURL := buildWSEndPointURL(serverURL)

	hc, ok := ps.hostClients[wsURL]
	if !ok {
		hc = client.NewHostClient(wsURL)
		err := hc.Start()

		if err != nil {
			return nil, fmt.Errorf("Cannot connect to %s: %v", wsURL, err)
		}
		ps.hostClients[wsURL] = hc
	}
	return hc, nil
}

// ConnectSharedPage establishes a new connection to the specified shared page and returns file name of control pipe.
func (ps *Service) ConnectSharedPage(ctx context.Context, args *ConnectPageArgs, results *ConnectPageResults) error {

	pageName := args.PageName
	serverURL, err := getServerURL(args.Web, args.Server)

	if err != nil {
		return err
	}

	hc, err := ps.getHostClient(serverURL)
	if err != nil {
		log.Errorln(err)
		return err
	}

	log.Println("Connecting to shared page:", pageName)

	// call server
	result := hc.Call(ctx, page.RegisterHostClientAction, &page.RegisterHostClientRequestPayload{
		PageName:  pageName,
		IsApp:     false,
		AuthToken: args.Token,
	})

	// parse response
	payload := &page.RegisterHostClientResponsePayload{}
	err = json.Unmarshal(*result, payload)

	if err != nil {
		log.Errorln("Error parsing ConnectSharedPage response:", err)
		return err
	}

	if payload.Error != "" {
		log.Errorln("Error calling ConnectSharedPage:", payload.Error)
		return errors.New(payload.Error)
	}

	results.PageName = payload.PageName
	results.PageURL = getPageURL(serverURL, payload.PageName)

	// create new pipeClient
	pc, err := client.NewPipeClient(payload.PageName, payload.SessionID, hc, args.Uds, args.EmitAllEvents, args.TickerDuration)
	if err != nil {
		return err
	}

	pc.Start()

	// register pipe client, so it can receive events from pages/sessions
	hc.RegisterPipeClient(pc)

	results.PipeName = pc.CommandPipeName()

	return nil
}

// ConnectAppPage waits for new web clients connecting specified page, creates a new session and returns file name of control pipe.
func (ps *Service) ConnectAppPage(ctx context.Context, args *ConnectPageArgs, results *ConnectPageResults) error {

	pageName := args.PageName
	serverURL, err := getServerURL(args.Web, args.Server)

	if err != nil {
		return err
	}

	hc, err := ps.getHostClient(serverURL)
	if err != nil {
		log.Errorln(err)
		return err
	}

	log.Println("Connecting to app page:", pageName)

	// call server
	result := hc.Call(ctx, page.RegisterHostClientAction, &page.RegisterHostClientRequestPayload{
		PageName:  pageName,
		IsApp:     true,
		AuthToken: args.Token,
	})

	// parse response
	payload := &page.RegisterHostClientResponsePayload{}
	err = json.Unmarshal(*result, payload)

	if err != nil {
		log.Errorln("Error parsing ConnectAppPage response:", err)
		return err
	}

	if payload.Error != "" {
		log.Errorln("Error calling ConnectAppPage:", payload.Error)
		return errors.New(payload.Error)
	}

	log.Println("Connected to app page:", payload.PageName)

	results.PageName = payload.PageName
	results.PageURL = getPageURL(serverURL, payload.PageName)

	return nil
}

func (ps *Service) WaitAppSession(ctx context.Context, args *ConnectPageArgs, results *ConnectPageResults) error {

	pageName := args.PageName
	serverURL, err := getServerURL(args.Web, args.Server)

	if err != nil {
		return err
	}

	hc, err := ps.getHostClient(serverURL)
	if err != nil {
		log.Errorln(err)
		return err
	}

	log.Println("Waiting for a new app session:", pageName)

	var sessionID string

	// wait for new session
	select {
	case sessionID = <-hc.PageNewSessions(pageName):
		break
	case <-ctx.Done():
		return errors.New("abort waiting for new session")
	}

	// create new pipeClient
	pc, err := client.NewPipeClient(pageName, sessionID, hc, args.Uds, args.EmitAllEvents, args.TickerDuration)
	if err != nil {
		return err
	}

	pc.Start()

	// register pipe client, so it can receive events from pages/sessions
	hc.RegisterPipeClient(pc)

	results.PageName = pageName
	results.PageURL = getPageURL(serverURL, pageName)
	results.PipeName = pc.CommandPipeName()

	return nil
}

func Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Starting Proxy service...")

	m, err := filemutex.New(lockFilename)
	if err != nil {
		log.Fatalln("Directory did not exist or file could not be created")
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

	lc := net.ListenConfig{}
	l, e := lc.Listen(ctx, "unix", sockAddr)

	if e != nil {
		log.Fatal("listen error:", e)
	}

	srv := &http.Server{}

	go func() {
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Println("Serve error:", err)
		}
	}()

	log.Println("Waiting for connections...")

	<-ctx.Done()

	log.Println("Stopping proxy service...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Proxy service shutdown failed:%+s", err)
	}

	proxySvc.close()

	log.Println("Proxy service exited")
}

func (ps *Service) close() {
	// close all host clients
	for _, hc := range ps.hostClients {
		hc.Close()
	}
}

func buildWSEndPointURL(serverURL string) string {

	if serverURL == "" {
		return ""
	}

	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatalln("Cannot parse server URL:", err)
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	u.Path = "ws"
	u.RawQuery = ""

	return u.String()
}

func getServerURL(web bool, server string) (string, error) {

	if server == "" && web {
		return pgletIoURL, nil
	} else if server == "" {
		return "", nil
	}

	serverURL := strings.Trim(server, "/")

	if !strings.Contains(serverURL, "://") {
		// scheme is specified
		serverURL = "http://" + serverURL
	}

	return serverURL, nil
}

func getPageURL(serverURL string, pageName string) string {
	if serverURL == "" {
		serverURL = fmt.Sprintf("http://localhost:%d", server.Port)
	}
	return fmt.Sprintf("%s/%s", serverURL, pageName)
}
