package commands

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/pglet/pglet/internal/config"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	version  = "unknown"
	LogLevel string
)

func configLogging() {

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
}
