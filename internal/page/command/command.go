package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	commandRegexPattern string = `(?:(\w+(?:\:\w+)*)[\s]*(?:=[\s]*((?:[^"'\s]+)|'(?:[^']*)'|"(?:[^"]*)"))?)`
)

const (
	Add    string = "add"
	Addf          = "addf"
	Set           = "set"
	Get           = "get"
	Clean         = "clean"
	Remove        = "remove"
	Insert        = "insert"
	Quit          = "quit"
)

var (
	supportedCommands = map[string]*CommandMetadata{
		Add:    &CommandMetadata{Name: Add, ShouldReturn: true},
		Addf:   &CommandMetadata{Name: Addf, ShouldReturn: false},
		Set:    &CommandMetadata{Name: Set, ShouldReturn: false},
		Get:    &CommandMetadata{Name: Get, ShouldReturn: true},
		Clean:  &CommandMetadata{Name: Clean, ShouldReturn: false},
		Remove: &CommandMetadata{Name: Remove, ShouldReturn: false},
		Insert: &CommandMetadata{Name: Insert, ShouldReturn: false},
		Quit:   &CommandMetadata{Name: Quit, ShouldReturn: false},
	}
)

type Command struct {
	Name   string // mandatory command name
	Values []string
	Attrs  map[string]string
}

type CommandMetadata struct {
	Name         string
	ShouldReturn bool
}

func Parse(cmdText string) (*Command, error) {
	re := regexp.MustCompile(commandRegexPattern)
	matches := re.FindAllSubmatch([]byte(cmdText), -1)

	command := &Command{
		Attrs:  make(map[string]string),
		Values: make([]string, 0),
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
			_, commandExists := supportedCommands[command.Name]
			if !commandExists {
				return nil, fmt.Errorf("Unknown command: %s", command.Name)
			}
		} else if value == "" {
			// bare value
			command.Values = append(command.Values, key)
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
