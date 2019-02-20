package arnoldc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

//go:generate goyacc -l -o parser.go parser.y

type ArnoldC struct {
	input   io.ReadSeeker
	program Program // Our result object
	err     error   // The last error to come out of the parser
	Debug   bool    // Print debugging messages
	// State
	scannedInstruction bool
	// Where are we in out input?
	line        int
	offset      int
	currentByte byte
	// TODO: Keep the current line in memory so we can add it to errors?
}

func New(input io.ReadSeeker) *ArnoldC {
	return &ArnoldC{input: input, line: 1, Debug: false}
}

func (a *ArnoldC) Parse() (*Program, error) {
	_ = yyParse(a)
	if a.err != nil {
		return nil, a.err
	}
	return &a.program, nil
}

func (a *ArnoldC) log(format string, args ...interface{}) {
	if a.Debug {
		fmt.Fprintf(os.Stderr, "Lexer: "+format+"\n", args...)
	}
}

// Error satisfies yyLexer.
func (a *ArnoldC) Error(s string) {
	a.err = fmt.Errorf("%s on line %d at character %d, %q", s, a.line, a.offset, a.currentByte)
}

// Lex satisfies yyLexer.
func (a *ArnoldC) Lex(lval *yySymType) int {
	if !a.scannedInstruction {
		return a.scanInstruction(lval)
	}
	return a.scanNormal(lval)
}

// Gets the next byte in the input.
func (a *ArnoldC) next() byte {
	// TODO: Optimize this buffer creation.
	// TODO: We can also scan like, 10 chars and walk them to eat up space and stuff.
	var b []byte = make([]byte, 1)

	_, err := a.input.Read(b)
	if err != nil {
		if err == io.EOF {
			a.log("!! EOF")
		}
		// TODO: Set the non-EOF error contents somewhere?
		return 0
	}

	a.offset++
	a.currentByte = b[0]

	return b[0]
}

// Seek input backwards one byte.
func (a *ArnoldC) backup() error {
	a.log("<- backup")
	_, err := a.input.Seek(-1, io.SeekCurrent)
	if err == nil {
		a.offset--
		a.currentByte = 0
	}
	return err
}

func (a *ArnoldC) scanInstruction(lval *yySymType) int {
	a.log("> scanInstruction")
	buf := bytes.NewBuffer(nil)
	for b := a.next(); ; b = a.next() {
		a.log("%q  %d", b, b)
		switch {
		case (unicode.IsUpper(rune(b)) || b == '\'' || unicode.IsSpace(rune(b))) && b != '\n':
			buf.WriteByte(b)
		default:
			lval.str = strings.TrimSpace(buf.String())
			if len(lval.str) == 0 {
				if b == 0 {
					return 0
				}
				continue
			}
			a.log("str = %q", lval.str)
			if b != 0 {
				a.backup()
			}
			a.scannedInstruction = true
			tk, ok := instructions[lval.str]
			if !ok {
				a.log("!! Unknown Instruction")
				// TODO: Store error.
				return LexError
			}
			if '\n' == b {
				a.line++
				a.offset = 0
			}
			return tk
		}
	}
}

func (a *ArnoldC) scanNormal(lval *yySymType) int {
	a.log("> scanNormal")
	for b := a.next(); b != 0; b = a.next() {
		a.log("%q  %d", b, b)
		switch {
		case '\n' == b:
			a.line++
			a.offset = 0
			a.scannedInstruction = false
			return a.scanInstruction(lval)
		case unicode.IsSpace(rune(b)):
			continue
		case b == '"':
			return a.scanString(lval)
		case unicode.IsDigit(rune(b)):
			a.backup()
			return a.scanInteger(lval)
		case unicode.IsLetter(rune(b)):
			a.backup()
			return a.scanVariable(lval)
		default:
			return int(b)
		}
	}

	return 0
}

func (a *ArnoldC) scanString(lval *yySymType) int {
	a.log("> scanString")
	buf := bytes.NewBuffer(nil)
	for b := a.next(); b != 0; b = a.next() {
		a.log("%q  %d", b, b)
		switch b {
		case '"':
			lval.str = buf.String()
			a.log("str = %q", lval.str)
			return String
		default:
			buf.WriteByte(b)
		}
	}

	return 0
}

func (a *ArnoldC) scanVariable(lval *yySymType) int {
	a.log("> scanVariable")
	buf := bytes.NewBuffer(nil)
	for b := a.next(); ; b = a.next() {
		a.log("%q  %d", b, b)
		switch {
		case unicode.IsLetter(rune(b)):
			buf.WriteByte(b)
		case b == '\n':
			a.backup()
			fallthrough
		default:
			lval.str = buf.String()
			a.log("str = %q", lval.str)
			return Variable
		}
	}

	return 0
}

func (a *ArnoldC) scanInteger(lval *yySymType) int {
	a.log("> scanInteger")
	buf := bytes.NewBuffer(nil)
	for b := a.next(); ; b = a.next() {
		a.log("%q  %d", b, b)
		if unicode.IsDigit(rune(b)) {
			buf.WriteByte(b)
		} else {
			a.backup()
			val, err := strconv.Atoi(buf.String())
			if err != nil {
				a.log("!! Invalid Integer")
				// TODO: Store error.
				return LexError
			}
			lval.integer = val
			a.log("integer = %d", lval.integer)
			return Integer
		}
	}

	return 0
}
