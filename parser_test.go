package arnoldc

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	src := `IT'S SHOWTIME
HEY CHRISTMAS TREE myVar
YOU SET US UP 10
TALK TO THE HAND "hello world"
YOU HAVE BEEN TERMINATED`

	f, err := fileFromString(src)
	if err != nil {
		t.Fatal(err)
	}
	defer func(f *os.File) {
		f.Close()
		os.Remove(f.Name())
	}(f)

	p := ArnoldC{input: f}
	p.Debug = true
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("error parsing: %v", err)
	}

	expect := Program{
		Main: Function{
			Expressions: []Expression{
				Expression{
					Instruction: "HEY CHRISTMAS TREE",
					Args:        []Value{VariableValue{"myVar"}, IntegerValue{10}},
				},
				Expression{
					Instruction: "TALK TO THE HAND",
					Args:        []Value{StringValue{"hello world"}},
				},
			},
		},
	}

	if expect.String() != program.String() {
		t.Errorf("Program does not match expectations.\n  Expected:\n%s\n  Got:\n%s", expect, program)
	}
}
