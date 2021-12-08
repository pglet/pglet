package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-filemutex"
	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/config"
	"github.com/pglet/pglet/internal/proxy"
	"github.com/pglet/pglet/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	defaultPort  int = 5000
	lockFilename string
)

func init() {
	lockFilename = filepath.Join(os.TempDir(), "pglet.lock")
}

func newServerCommand() *cobra.Command {

	var serverPort int
	var background bool

	var cmd = &cobra.Command{
		Use:   "server",
		Short: "Start server service",
		Long:  `Server is for ...`,
		Run: func(cmd *cobra.Command, args []string) {

			if background {
				startServerService()
				return
			}

			// ensure one executable instance is running
			m, err := filemutex.New(lockFilename)
			if err != nil {
				log.Fatalln("Cannot create mutex - directory did not exist or file could not be created")
			}

			err = m.TryLock()
			if err != nil {
				log.Fatalln("Another Pglet Server process has already started")
			}

			defer m.Unlock()

			// init cache
			cache.Init()

			waitGroup := sync.WaitGroup{}
			waitGroup.Add(1)
			go server.Start(cmd.Context(), &waitGroup, serverPort)
			waitGroup.Wait()
		},
	}

	cmd.Flags().IntVarP(&serverPort, "port", "p", config.ServerPort(), "port on which the server will listen")
	cmd.Flags().BoolVarP(&background, "background", "b", false, "run server in background")

	return cmd
}

func startServerService() {
	log.Traceln("Starting Pglet Server")

	// run server
	execPath, _ := os.Executable()

	cmd := proxy.GetDetachedCmd(execPath, "server")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=true", config.LogToFileFlag))

	err := cmd.Start()

	if err != nil {
		log.Fatalln(err)
	}

	log.Traceln("Server process started with PID:", cmd.Process.Pid)
}
