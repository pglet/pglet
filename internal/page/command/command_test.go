package command

import (
	"log"
	"testing"

	"github.com/pglet/pglet/internal/utils"
)

func TestParse(t *testing.T) {
	cmd, err := Parse(`Add value:1 c=3.1 TextField Text=aaa value="Hello,\n 'wor\"ld!" aaa='bbb' cmd2=1`)

	if err != nil {
		t.Fatal("Error parsing command", err)
	}

	// visualize command
	log.Printf("%s", utils.ToJSON(cmd))

	if cmd.Name != "add" {
		t.Errorf("command name is %s, want %s", cmd.Name, "add")
	}
}

func TestParse2(t *testing.T) {
	cmd, err := Parse(`set body:form:fullName value='John Smith' another_prop=value`)

	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%s", utils.ToJSON(cmd))
	}
}
