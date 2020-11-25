package command

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"

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

func Parse(cmdText string) (*Command, error) {

	command := &Command{
		Attrs:  make(map[string]string),
		Values: make([]string, 0),
	}

	var err error
	var s scanner.Scanner
	s.Init(strings.NewReader(cmdText))
	s.Filename = "command"
	s.Error = func(s *scanner.Scanner, msg string) {
		if msg != "invalid char literal" {
			err = fmt.Errorf("error parsing command at position %d: %s", s.Column, msg)
		}
	}

	// treat ':' as part of an identifier
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == ':' || ch == '_' || ch == '-' || ch == '.' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
	}

	prevToken := ""
	prevLit := ""
	for r := s.Scan(); r != scanner.EOF; r = s.Scan() {

		if err != nil {
			return nil, err
		}

		tok := s.TokenText()

		fmt.Printf("%s: %s\n", s.Position, tok)

		if tok == "=" {
			if prevLit == "" || prevToken == "=" {
				return nil, fmt.Errorf("unexpected '=' at position %d", s.Column)
			}
		} else if tok != "=" && prevToken == "=" && prevLit != "" {
			// name=value
			command.Attrs[strings.ToLower(utils.TrimQuotes(prevLit))] = utils.TrimQuotes(tok)
			prevLit = ""
		} else if tok != "=" && prevToken != "=" && prevLit != "" {
			v := utils.TrimQuotes(prevLit)
			if command.Name == "" {
				command.Name = strings.ToLower(v)
				_, commandExists := supportedCommands[command.Name]
				if !commandExists {
					return nil, fmt.Errorf("Unknown command: %s", command.Name)
				}
			} else {
				command.Values = append(command.Values, v)
			}
			prevLit = tok
		} else {
			prevLit = tok
		}
		prevToken = tok
	}

	// consume last token collected
	if prevLit != "" {
		command.Values = append(command.Values, prevLit)
	}

	return command, nil
}

func (cmd *Command) ShouldReturn() bool {
	cmdMeta, _ := supportedCommands[strings.ToLower(cmd.Name)]
	return cmdMeta.ShouldReturn
}
