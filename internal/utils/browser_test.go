package utils

import (
	"fmt"
	"testing"
)

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
