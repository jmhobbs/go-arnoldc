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
			name:        "Ignore Spaces",
			source:      " \t",
			token_types: []int{},
		},
		{
			name:        "Passthrough Non-Trigger Characters",
			source:      "TALK TO THE HAND {",
			token_types: []int{PRINT, int('{')},
		},
		{
			name:        "Invalid Statement",
			source:      `THIS IS NOT A STATEMENT`,
			token_types: []int{LexError},
		},
		{
			name:        "Main Open",
			source:      `IT'S SHOWTIME`,
			token_types: []int{MAIN_OPEN},
		},
		{
			name:        "Print String",
			source:      `TALK TO THE HAND "hello world"`,
			token_types: []int{PRINT, String},
		},
		{
			name:        "Incomplete String",
			source:      `TALK TO THE HAND "hello world`,
			token_types: []int{PRINT, LexError},
		},
		{
			name:        "Print Variable",
			source:      `TALK TO THE HAND myVar`,
			token_types: []int{PRINT, Variable},
		},
		{
			name:        "Variable with number",
			source:      `TALK TO THE HAND myVar10`,
			token_types: []int{PRINT, Variable},
		},
		{
			name:        "Main Close",
			source:      `YOU HAVE BEEN TERMINATED`,
			token_types: []int{MAIN_CLOSE},
		},
		{
			name: "Full Main Function",
			source: `
IT'S SHOWTIME
TALK TO THE HAND "hello world"
YOU HAVE BEEN TERMINATED`,
			token_types: []int{MAIN_OPEN, PRINT, String, MAIN_CLOSE},
		},
		{
			name: "Ignore Multiple Newlines",
			source: `
IT'S SHOWTIME


TALK TO THE HAND "hello world"


YOU HAVE BEEN TERMINATED`,
			token_types: []int{MAIN_OPEN, PRINT, String, MAIN_CLOSE},
		},
		{
			name:        "Declare Variable",
			source:      "HEY CHRISTMAS TREE myVar",
			token_types: []int{DECLARE, Variable},
		},
		{
			name:        "False Bool Macro",
			source:      "YOU SET US UP @I LIED",
			token_types: []int{INITIALIZE, FALSE},
		},
		{
			name:        "True Bool Macro",
			source:      "YOU SET US UP @NO PROBLEMO",
			token_types: []int{INITIALIZE, TRUE},
		},
		{
			name:        "Invalid Bool Macro",
			source:      "YOU SET US UP @NOT REAL",
			token_types: []int{INITIALIZE, LexError},
		},
		{
			name: "Void Method",
			source: `
LISTEN TO ME VERY CAREFULLY methodName
TALK TO THE HAND "hello world"
HASTA LA VISTA, BABY`,
			// TODO: Using "Variable" for the method name and args is...not great.
			token_types: []int{METHOD_OPEN, Variable, PRINT, String, METHOD_CLOSE},
		},
		{
			name: "Parameterized Method",
			source: `
LISTEN TO ME VERY CAREFULLY methodName
I NEED YOUR CLOTHES YOUR BOOTS AND YOUR MOTORCYCLE arg1
GIVE THESE PEOPLE AIR
TALK TO THE HAND "hello world"
HASTA LA VISTA, BABY`,
			token_types: []int{METHOD_OPEN, Variable, DECLARE_PARAMETER, Variable, END_PARAMETER_DECLARATION, PRINT, String, METHOD_CLOSE},
		},
		{
			name:        "Negative Integers",
			source:      "TALK TO THE HAND -150",
			token_types: []int{PRINT, Integer},
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
	case MAIN_OPEN:
		return "main() Open"
	case MAIN_CLOSE:
		return "main() Close"
	case METHOD_OPEN:
		return "method() Open"
	case METHOD_CLOSE:
		return "method() Close"
	case PRINT:
		return "print()"
	}
	if typ < 128 {
		return fmt.Sprintf("%q", typ)
	}
	return "UNKNOWN"
}
