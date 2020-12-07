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

func TestOpenBrowser2(t *testing.T) {

	OpenBrowser("http://localhost:3000", "100, 200, 600, 600")
	//t.Errorf("ddd")
}

func TestFindChrome(t *testing.T) {
	path := findChrome()
	fmt.Println("path:", path)
}
