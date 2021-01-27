package command

import (
	"log"
	"testing"
)

func TestParse1(t *testing.T) {
	cmd, err := Parse(`  Add value:1 c=3.1 TextField Text=aaa value="Hello,\n 'wor\"ld!" aaa='bbb' cmd2=1`, true)

	if err != nil {
		t.Fatal("Error parsing command", err)
	}

	// visualize command
	log.Printf("%s", cmd)

	if cmd.Name != "Add" {
		t.Errorf("command name is %s, want %s", cmd.Name, "Add")
	}

	if cmd.Indent != 2 {
		t.Errorf("command indent is %d, want %d", cmd.Indent, 2)
	}
}

func TestParse2(t *testing.T) {
	cmd, err := Parse(`set body:form:fullName value='John Smith' another_prop=value`, true)

	if err != nil {
		t.Fatal(err)
	}

	// visualize command
	log.Printf("%s", cmd)

	if len(cmd.Values) != 1 {
		t.Errorf("the number of values is %d, want %d", len(cmd.Values), 1)
	}

	expValue := "body:form:fullName"
	if cmd.Values[0] != expValue {
		t.Errorf("command values[0] is %s, want %s", cmd.Values[0], expValue)
	}

	if cmd.Indent != 0 {
		t.Errorf("command indent is %d, want %d", cmd.Indent, 0)
	}
}

func TestParseSingleCommand(t *testing.T) {
	cmd, err := Parse(`set`, true)

	if err != nil {
		t.Fatal(err)
	}

	// visualize command
	log.Printf("%s", cmd)

	if cmd.Name != "set" {
		t.Errorf("command name is %s, want %s", cmd.Name, "set")
	}

	if len(cmd.Values) != 0 {
		t.Errorf("the number of values is %d, want %d", len(cmd.Values), 0)
	}
}

func TestParseClean(t *testing.T) {
	cmd, err := Parse(`clean page`, true)

	if err != nil {
		t.Fatal(err)
	}

	// visualize command
	log.Printf("%s", cmd)

	if len(cmd.Values) != 1 {
		t.Errorf("the number of values is %d, want %d", len(cmd.Values), 1)
	}

	expValue := "page"
	if cmd.Values[0] != expValue {
		t.Errorf("command values[0] is %s, want %s", cmd.Values[0], expValue)
	}
}

func TestParseSlashesAndNewLine(t *testing.T) {
	cmd, err := Parse(`  Add value="C:\\Program Files\\Node\\node.exe" aaa='bbb' cmd2=1`, true)

	if err != nil {
		t.Fatal("Error parsing command", err)
	}

	// visualize command
	log.Printf("%s", cmd)

	if cmd.Name != "Add" {
		t.Errorf("command name is %s, want %s", cmd.Name, "Add")
	}

	expAttr := "C:\\Program Files\\Node\\node.exe"
	if cmd.Attrs["value"] != expAttr {
		t.Errorf("command value attribute is %s, want %s", cmd.Attrs["value"], expAttr)
	}
}

func TestParseMultilineCommand(t *testing.T) {
	cmd, err := Parse(`

	  add to=footer
	    stack
	      text value="Hello, world!"
	    stack
		  textbox id=txt1
		  button id=ok`, true)

	if err != nil {
		t.Fatal(err)
	}

	// visualize command
	log.Printf("%s", cmd)

	expName := "add"
	if cmd.Name != expName {
		t.Errorf("command name is %s, want %s", cmd.Name, expName)
	}

	if len(cmd.Values) != 0 {
		t.Errorf("the number of values is %d, want %d", len(cmd.Values), 0)
	}

	if len(cmd.Lines) != 5 {
		t.Errorf("the number of lines is %d, want %d", len(cmd.Lines), 5)
	}

	if cmd.Indent != 6 {
		t.Errorf("command indent is %d, want %d", cmd.Indent, 6)
	}
}
