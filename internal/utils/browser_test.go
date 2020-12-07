package utils

import (
	"fmt"
	"os"
	"testing"
)

// skip the tests in CI environment (https://stackoverflow.com/questions/24030059/skip-some-tests-with-go-test)
func TestMain(t *testing.M) {
	if os.Getenv("CI") == "" {
		t.Run()
	}
}

func TestOpenBrowserDefault(t *testing.T) {

	OpenBrowser("http://localhost:3000", "")
	//t.Errorf("ddd")
}

func TestOpenBrowserWidthHeightOnly(t *testing.T) {

	OpenBrowser("http://localhost:3000", "400, 400")
	//t.Errorf("ddd")
}

func TestOpenBrowser(t *testing.T) {

	OpenBrowser("http://localhost:3000", "100, 200, 800, 600")
	//t.Errorf("ddd")
}

func TestOpenSafari(t *testing.T) {

	width, height := getMonitorSize()
	openSafari("http://google.com", 100, 100, width, height)
	//t.Errorf("ddd")
}

func TestResolution(t *testing.T) {
	width, height := getMonitorSize()
	fmt.Println("width:", width, "height:", height)
}

func TestFindChrome(t *testing.T) {
	path := findChrome()
	fmt.Println("path:", path)
}
