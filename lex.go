package arnoldc

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

//go:generate goyacc -l -o parser.go parser.y

type ArnoldC struct {
	input   io.ReadSeeker
	program Program // Our result object
	err     error   // The last error to come out of the parser
	// State
	scannedInstruction bool
	// Where are we in out input?
	line        int
	offset      int
	currentByte byte
}

func New(input io.ReadSeeker) *ArnoldC {
	return &ArnoldC{input: input}
}

func (a *ArnoldC) Parse() (*Program, error) {
	_ = yyParse(a)
	if a.err != nil {
		return nil, a.err
	}
	return &a.program, nil
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
		// TODO: Set the error contents somewhere?
		//if err != io.EOF {
		//	return byte(0)
		//}
		return 0
	}

	a.offset++
	a.currentByte = b[0]

	return b[0]
}

// Seek input backwards one byte.
func (a *ArnoldC) backup() error {
	_, err := a.input.Seek(-1, io.SeekCurrent)
	if err == nil {
		a.offset--
		a.currentByte = 0
	}
	return err
}

func (a *ArnoldC) scanInstruction(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := a.next(); ; b = a.next() {
		switch {
		case (unicode.IsUpper(rune(b)) || b == '\'' || unicode.IsSpace(rune(b))) && b != '\n':
			buf.WriteByte(b)
		default:
			lval.str = strings.TrimSpace(buf.String())
			if len(lval.str) == 0 {
				continue
			}
			if b != 0 {
				a.backup()
			}
			a.scannedInstruction = true
			tk, ok := instructions[lval.str]
			if !ok {
				return LexError
			}
			if '\n' == b {
				a.line++
				a.offset = 0
			}
			return tk
		}
	}
	return 0
}

func (a *ArnoldC) scanNormal(lval *yySymType) int {
	for b := a.next(); b != 0; b = a.next() {
		switch {
		case '\n' == b:
			a.line++
			a.offset = 0
			a.scannedInstruction = false
			return int(b)
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
	buf := bytes.NewBuffer(nil)
	for b := a.next(); b != 0; b = a.next() {
		switch b {
		case '"':
			lval.str = buf.String()
			return String
		default:
			buf.WriteByte(b)
		}
	}

	return 0
}

func (a *ArnoldC) scanVariable(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := a.next(); ; b = a.next() {
		switch {
		case unicode.IsLetter(rune(b)):
			buf.WriteByte(b)
		default:
			lval.str = buf.String()
			return Variable
		}
	}

	return 0
}

func (a *ArnoldC) scanInteger(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := a.next(); ; b = a.next() {
		if unicode.IsDigit(rune(b)) {
			buf.WriteByte(b)
		} else {
			a.backup()
			val, err := strconv.Atoi(buf.String())
			if err != nil {
				return LexError
			}
			lval.integer = val
			return Integer
		}
	}

	return 0
}
