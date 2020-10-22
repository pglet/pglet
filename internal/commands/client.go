package commands

import (
	"fmt"
)

var (
	version = "unknown"
	commit  = "unknown"
)

func PrintVersion() {
	fmt.Println(version, commit)
}
