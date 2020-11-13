package utils

import (
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Source: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenBrowser(url string) {
	var err error

	goos := runtime.GOOS

	if goos == "linux" {
		// check if it's WSL
		content, err := ioutil.ReadFile("/proc/version")
		if err == nil {
			version := strings.ToLower(string(content))
			if strings.Contains(version, "microsoft") {
				goos = "wsl"
			}
		}
	}

	switch goos {
	//case "linux":
	//	err = exec.Command("xdg-open", url).Start()
	case "windows", "wsl":
		err = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
	}

	if err != nil {
		log.Warnln("Error opening browser window:", err)
	}
}
