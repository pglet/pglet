package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

// Source: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenBrowser(url string) {
	var err error

	goos := runtime.GOOS

	if goos == "linux" {
		// check if it's WSL
		content, err := ioutil.ReadFile("/cat/version")
		if err == nil {
			version := string(content)
			if strings.Contains(version, "Microsoft") {
				goos = "wsl"
			}
		}
	}

	switch goos {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows", "wsl":
		err = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
