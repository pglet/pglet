package page

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/pglet/pglet/utils"
)

const (
	commandRegexPattern string = `(?:(\w+(?:\:\w+)*)[\s]*(?:=[\s]*((?:[^"'\s]+)|'(?:[^']*)'|"(?:[^"]*)"))?)`
)

const (
	Add    string = "add"
	Addr          = "addr"
	Set           = "set"
	Get           = "get"
	Clean         = "clean"
	Remove        = "remove"
	Insert        = "insert"
	Quit          = "quit"
)

var (
	supportedCommands = []string{
		Add,
		Addr,
		Set,
		Get,
		Clean,
		Remove,
		Insert,
		Quit,
	}
)

type Command struct {
	Name  string // mandatory command name
	Attrs map[string]string
}

func ParseCommand(cmdText string) (*Command, error) {
	re := regexp.MustCompile(commandRegexPattern)
	matches := re.FindAllSubmatch([]byte(cmdText), -1)

	command := &Command{
		Attrs: make(map[string]string),
	}

	order := 1
	for _, m := range matches {
		key := string(m[1])
		value := string(m[2])

		if order == 1 {
			// verb
			if key == "" || value != "" {
				return nil, errors.New("The first argument of the command must be a verb")
			}
			command.Name = strings.ToLower(key)
			if !utils.ContainsString(supportedCommands, command.Name) {
				return nil, fmt.Errorf("Unknown command: %s", command.Name)
			}
		} else {
			// attr
			if strings.HasPrefix(value, "\"") {
				value = strings.Trim(value, "\"")
			} else if strings.HasPrefix(value, "'") {
				value = strings.Trim(value, "'")
			}
			command.Attrs[strings.ToLower(key)] = value
		}
		order++
	}

	return command, nil
}
