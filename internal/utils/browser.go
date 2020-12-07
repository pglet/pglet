package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/andybrewer/mack"
	"github.com/kbinani/screenshot"
	log "github.com/sirupsen/logrus"
)

// Source: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenBrowser(url string, posSize string) {

	if posSize == "" {
		openDefaultBrowser(url)
		return
	}

	x := -1
	y := -1
	width := -1
	height := -1

	defWidth, defHeight := getMonitorSize()

	sizeParts := strings.Split(posSize, ",")
	if len(sizeParts) > 3 {
		x, _ = strconv.Atoi(strings.TrimSpace(sizeParts[0]))
		y, _ = strconv.Atoi(strings.TrimSpace(sizeParts[1]))
		sizeParts = sizeParts[2:]
	}

	if len(sizeParts) > 1 {
		width, _ = strconv.Atoi(strings.TrimSpace(sizeParts[0]))
		height, _ = strconv.Atoi(strings.TrimSpace(sizeParts[1]))
	}

	if width == -1 {
		width = defWidth / 2
	}

	if height == -1 {
		height = defHeight / 2
	}

	if x == -1 {
		x = defWidth/2 - width/2
	}

	if y == -1 {
		y = defHeight/2 - height/2
	}

	chromePath := findChrome()
	if chromePath != "" {
		openChrome(chromePath, url, x, y, width, height)
	} else if runtime.GOOS == "darwin" {
		openSafari(url, x, y, width, height)
	} else {
		openDefaultBrowser(url)
	}
}

func openDefaultBrowser(url string) {
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

func openChrome(chromePath string, url string, x int, y int, width int, height int) {

	// create temp profile dir
	homeDir, _ := os.UserHomeDir()
	profilePath := filepath.Join(homeDir, ".pglet-chrome-profile")

	_, err := os.Stat(profilePath)
	if !os.IsNotExist(err) {
		os.RemoveAll(profilePath)
	}
	os.MkdirAll(profilePath, os.ModePerm)

	err = exec.Command(chromePath,
		"--chrome-frame",
		fmt.Sprintf("--user-data-dir=%s", profilePath),
		fmt.Sprintf("--window-position=%d,%d", x, y),
		fmt.Sprintf("--window-size=%d,%d", width, height),
		fmt.Sprintf("--app=%s", url),
		// "--inprivate",
		// "--incognito",
		"--disable-sync",
		"--no-first-run",
		"--noerrdialogs",
		"--no-default-browser-check",
	).Start()

	if err != nil {
		log.Warnln("Error opening Chrome window:", err)
	}
}

func openSafari(url string, x int, y int, width int, height int) {
	mack.Tell("Safari",
		fmt.Sprintf("open location \"%s\"", url),
		fmt.Sprintf("set bounds of front window to {%d, %d, %d, %d}", x, y, width, height),
		"activate")
}

func getMonitorSize() (width int, height int) {
	scale := getMonitorScale()
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		if bounds.Min.X == 0 && bounds.Min.Y == 0 {
			// primary monitor
			// calculate default window size and position
			width = bounds.Max.X / scale
			height = bounds.Max.Y / scale
		}
	}

	return
}

func findChrome() string {
	progFiles := os.Getenv("ProgramFiles")
	progFilesX86 := os.Getenv("ProgramFiles(x86)")
	paths := []string{
		//"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		//"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		fmt.Sprintf(`%s\Google\Chrome\Application\chrome.exe`, progFilesX86),
		fmt.Sprintf(`%s\Google\Chrome\Application\chrome.exe`, progFiles),
		fmt.Sprintf(`%s\Microsoft\Edge\Application\msedge.exe`, progFilesX86),
		fmt.Sprintf(`%s\Microsoft\Edge\Application\msedge.exe`, progFiles),
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			return path
		}
	}
	return ""
}
