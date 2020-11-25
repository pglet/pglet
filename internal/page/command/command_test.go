package command

import (
	"encoding/json"
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	cmd, _ := Parse2(`Add value1 c=3.1 TextField Text=aaa value="Hello,\n 'wor\"ld!" aaa='bbb' cmd2=1`)

	json, _ := json.MarshalIndent(cmd, "", "  ")

	log.Printf("%s", json)
}
