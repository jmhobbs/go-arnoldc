package arnoldc

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	cases := []struct {
		name        string
		source      string
		token_types []int
	}{
		{
			name:        "Invalid Statement",
			source:      `THIS IS NOT A STATEMENT`,
			token_types: []int{LexError},
		},
		{
			name:        "Main Open",
			source:      `IT'S SHOWTIME`,
			token_types: []int{TK_MAIN_OPEN},
		},
		{
			name:        "Print String",
			source:      `TALK TO THE HAND "hello world"`,
			token_types: []int{TK_PRINT, String},
		},
		{
			name:        "Print Variable",
			source:      `TALK TO THE HAND myVar`,
			token_types: []int{TK_PRINT, Variable},
		},
		{
			name:        "Variable with number",
			source:      `TALK TO THE HAND myVar10`,
			token_types: []int{TK_PRINT, Variable},
		},
		{
			name:        "Main Close",
			source:      `YOU HAVE BEEN TERMINATED`,
			token_types: []int{TK_MAIN_CLOSE},
		},
		{
			name: "Full Main Function",
			source: `
IT'S SHOWTIME
TALK TO THE HAND "hello world"
YOU HAVE BEEN TERMINATED`,
			token_types: []int{TK_MAIN_OPEN, TK_PRINT, String, TK_MAIN_CLOSE},
		},
		{
			name: "Ignore Multiple Newlines",
			source: `
IT'S SHOWTIME


TALK TO THE HAND "hello world"


YOU HAVE BEEN TERMINATED`,
			token_types: []int{TK_MAIN_OPEN, TK_PRINT, String, TK_MAIN_CLOSE},
		},
		{
			name:        "Declare Variable",
			source:      "HEY CHRISTMAS TREE myVar",
			token_types: []int{TK_DECLARE, Variable},
		},
		{
			name:        "Bool Macro",
			source:      "YOU SET US UP @I LIED",
			token_types: []int{TK_INITIALIZE, TK_FALSE},
		},
		{
			name: "Void Method",
			source: `
LISTEN TO ME VERY CAREFULLY methodName
TALK TO THE HAND "hello world"
HASTA LA VISTA, BABY`,
			// TODO: Using "Variable" for the method name and args is...not great.
			token_types: []int{TK_METHOD_OPEN, Variable, TK_PRINT, String, TK_METHOD_CLOSE},
		},
		{
			name: "Parameterized Method",
			source: `
LISTEN TO ME VERY CAREFULLY methodName
I NEED YOUR CLOTHES YOUR BOOTS AND YOUR MOTORCYCLE arg1
GIVE THESE PEOPLE AIR
TALK TO THE HAND "hello world"
HASTA LA VISTA, BABY`,
			token_types: []int{TK_METHOD_OPEN, Variable, TK_DECLARE_PARAMETER, Variable, TK_END_PARAMETER_DECLARATION, TK_PRINT, String, TK_METHOD_CLOSE},
		},
	}

	var (
		lval          yySymType
		lexer         *ArnoldC
		token_types   []int
		i, token_type int
	)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := fileFromString(tc.source)
			if err != nil {
				t.Fatal(err)
			}
			defer func(f *os.File) {
				f.Close()
				os.Remove(f.Name())
			}(f)

			lexer = &ArnoldC{input: f}
			token_types = []int{}
			for i = 0; ; i++ {
				token_type = lexer.Lex(&lval)
				if token_type == 0 {
					break
				}
				token_types = append(token_types, token_type)
				if token_type == LexError {
					break
				}
			}

			if !reflect.DeepEqual(token_types, tc.token_types) {
				t.Errorf("Unexpected token types for %q\n  Expect: %v\n     Got: %v", tc.source, stringifyTypes(tc.token_types), stringifyTypes(token_types))
			}
		})
	}
}

func fileFromString(src string) (*os.File, error) {
	tmp, err := ioutil.TempFile("", "test")
	if err != nil {
		return nil, err
	}

	if _, err := tmp.Write([]byte(src)); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return nil, err
	}

	tmp.Seek(0, io.SeekStart)

	return tmp, nil
}

func stringifyTypes(types []int) []string {
	sTypes := []string{}
	for _, t := range types {
		sTypes = append(sTypes, tokenTypeToString(t))
	}
	return sTypes
}

func tokenTypeToString(typ int) string {
	switch typ {
	case LexError:
		return "LexError"
	case String:
		return "String"
	case Variable:
		return "Variable"
	case TK_MAIN_OPEN:
		return "main() Open"
	case TK_MAIN_CLOSE:
		return "main() Close"
	case TK_METHOD_OPEN:
		return "method() Open"
	case TK_METHOD_CLOSE:
		return "method() Close"
	case TK_PRINT:
		return "print()"
	}
	if typ < 128 {
		return fmt.Sprintf("%q", typ)
	}
	return "UNKNOWN"
}
