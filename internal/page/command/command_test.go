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
		t.Fatal(err)
	}

	// visualize command
	log.Printf("%s", utils.ToJSON(cmd))

	if len(cmd.Values) != 1 {
		t.Errorf("the number of values is %d, want %d", len(cmd.Values), 1)
	}

	expValue := "body:form:fullName"
	if cmd.Values[0] != expValue {
		t.Errorf("command values[0] is %s, want %s", cmd.Values[0], expValue)
	}
}
