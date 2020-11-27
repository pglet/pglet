package command

import (
	"log"
	"testing"

	"github.com/pglet/pglet/internal/utils"
)

func TestParse1(t *testing.T) {
	cmd, err := Parse(`Add value:1 c=3.1 TextField Text=aaa value="Hello,\n 'wor\"ld!" aaa='bbb' cmd2=1`)

	if err != nil {
		t.Fatal("Error parsing command", err)
	}

	// visualize command
	log.Printf("%s", utils.ToJSON(cmd))

	if cmd.Name != "Add" {
		t.Errorf("command name is %s, want %s", cmd.Name, "Add")
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

func TestParseClean(t *testing.T) {
	cmd, err := Parse(`clean page`)

	if err != nil {
		t.Fatal(err)
	}

	// visualize command
	log.Printf("%s", utils.ToJSON(cmd))

	if len(cmd.Values) != 1 {
		t.Errorf("the number of values is %d, want %d", len(cmd.Values), 1)
	}

	expValue := "page"
	if cmd.Values[0] != expValue {
		t.Errorf("command values[0] is %s, want %s", cmd.Values[0], expValue)
	}
}

func TestParseMultilineCommand(t *testing.T) {
	cmd, err := Parse(`
	  add to=footer
	    stack
	      text value="Hello, world!"
	    stack
		  textbox id=txt1
		  button id=ok`)

	if err != nil {
		t.Fatal(err)
	}

	// visualize command
	log.Printf("%s", utils.ToJSON(cmd))

	// if len(cmd.Values) != 1 {
	// 	t.Errorf("the number of values is %d, want %d", len(cmd.Values), 1)
	// }

	// expValue := "page"
	// if cmd.Values[0] != expValue {
	// 	t.Errorf("command values[0] is %s, want %s", cmd.Values[0], expValue)
	// }
}
