package command

import (
	"fmt"
	"go/scanner"
	"go/token"
	"log"
	"testing"
)

func TestScanner(t *testing.T) {

	// var src = `add textbox text=aaa height=3 value="Hello,\n 'wor\"ld!" width = d\'d"dd`

	// var s scanner.Scanner
	// s.Init(strings.NewReader(src))
	// s.Filename = "command"
	// // s.Error = func(s *scanner.Scanner, msg string) {
	// // 	if msg == "invalid char literal" {
	// // 		return
	// // 	}
	// // }
	// for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
	// 	log.Printf("%s (%s): %s\n", s.Position, strconv.FormatBool(s.IsValid()), s.TokenText())
	// }

	var src = []byte(`add textbox text=aaa height=3 value="Hello,\n 'wor\"ld!" width = 'd\'d"dd'`)

	errorHandler := func(pos token.Position, msg string) {
		log.Println("error:", msg)
	}

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, errorHandler, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%s\n", fset.Position(pos), tok, lit)
	}
}
