package commands

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/pglet/pglet/internal/config"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version  = "unknown"
	commit   = "unknown"
	LogLevel string
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pglet",
		Short:   "Pglet",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			level := log.FatalLevel // default logging level
			level, err := log.ParseLevel(LogLevel)

			if err != nil {
				log.Fatalln(err)
			}

			log.SetLevel(level)

			formatter := &log.TextFormatter{
				FullTimestamp: true,
			}

			if runtime.GOOS == "windows" {
				formatter.ForceColors = true
			}

			log.SetFormatter(formatter)

			if os.Getenv(config.LogToFileFlag) == "true" {
				logPath := "/var/log/pglet.log"
				if runtime.GOOS == "windows" {
					logPath = filepath.Join(os.TempDir(), "pglet.log")
				}
				pathMap := lfshook.PathMap{
					logrus.DebugLevel: logPath,
					logrus.InfoLevel:  logPath,
					logrus.ErrorLevel: logPath,
				}
				log.AddHook(lfshook.NewHook(
					pathMap,
					&log.TextFormatter{}))
			}
		},
	}

	cmd.SetVersionTemplate("{{.Version}}")

	cmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "verbosity level for logs")

	cmd.AddCommand(
		newPageCommand(),
		newAppCommand(),
		newServerCommand(),
		newClientCommand(),
	)

	return cmd
}
