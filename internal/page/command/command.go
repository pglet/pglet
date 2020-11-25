package command

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"log"
	"regexp"
	"strings"

	"github.com/pglet/pglet/internal/utils"
)

const (
	commandRegexPattern string = `(?:(\w+(?:\:\w+)*)[\s]*(?:=[\s]*((?:[^"'\s]+)|'(?:[^']*)'|"(?:[^"]*)"))?)`
)

const (
	Add     string = "add"
	Addf           = "addf"
	Set            = "set"
	Setf           = "setf"
	Get            = "get"
	Clean          = "clean"
	Cleanf         = "cleanf"
	Remove         = "remove"
	Removef        = "removef"
	Quit           = "quit"
)

var (
	supportedCommands = map[string]*CommandMetadata{
		Add:     &CommandMetadata{Name: Add, ShouldReturn: true},
		Addf:    &CommandMetadata{Name: Addf, ShouldReturn: false},
		Set:     &CommandMetadata{Name: Set, ShouldReturn: true},
		Setf:    &CommandMetadata{Name: Setf, ShouldReturn: false},
		Get:     &CommandMetadata{Name: Get, ShouldReturn: true},
		Clean:   &CommandMetadata{Name: Clean, ShouldReturn: true},
		Cleanf:  &CommandMetadata{Name: Cleanf, ShouldReturn: false},
		Remove:  &CommandMetadata{Name: Remove, ShouldReturn: true},
		Removef: &CommandMetadata{Name: Removef, ShouldReturn: false},
		Quit:    &CommandMetadata{Name: Quit, ShouldReturn: false},
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

func Parse2(cmdText string) (*Command, error) {

	command := &Command{
		Attrs:  make(map[string]string),
		Values: make([]string, 0),
	}

	var errs scanner.ErrorList
	errorHandler := func(pos token.Position, msg string) {
		if msg != "illegal rune literal" {
			errs.Add(pos, msg)
		}
	}

	var src = []byte(cmdText)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, errorHandler, scanner.ScanComments)

	prevToken := token.ILLEGAL
	prevLit := ""
	for {
		pos, tok, lit := s.Scan()
		p := fset.Position(pos)

		fmt.Printf("%s\t%s\t%s\n", p, tok, lit)

		if tok == token.EOF {
			break
		} else if tok == token.ASSIGN {
			if prevToken == token.ILLEGAL || prevToken == token.ASSIGN {
				return nil, fmt.Errorf("Unexpected = at %d", p.Column)
			}
		} else if tok != token.ASSIGN && prevToken == token.ASSIGN && prevLit != "" {
			// name=value
			command.Attrs[strings.ToLower(utils.TrimQuotes(prevLit))] = utils.TrimQuotes(lit)
			prevLit = ""
		} else if tok != token.ASSIGN && prevToken != token.ASSIGN && prevLit != "" {
			command.Values = append(command.Values, strings.ToLower(utils.TrimQuotes(prevLit)))
			prevLit = lit
		} else {
			prevLit = lit
		}
		prevToken = tok
	}

	for _, e := range errs {
		log.Printf("error: %d - %s", e.Pos.Column, e.Msg)
	}

	return command, nil
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

func (cmd *Command) ShouldReturn() bool {
	cmdMeta, _ := supportedCommands[strings.ToLower(cmd.Name)]
	return cmdMeta.ShouldReturn
}
