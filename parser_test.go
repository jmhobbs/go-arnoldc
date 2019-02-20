package arnoldc

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	src := `
IT'S SHOWTIME
	HEY CHRISTMAS TREE myVar
		YOU SET US UP 10

	TALK TO THE HAND "hello world"

	GET TO THE CHOPPER a
		HERE IS MY INVITATION 4
		GET UP b
		YOU'RE FIRED 2
	ENOUGH TALK

	GET TO THE CHOPPER a
		HERE IS MY INVITATION b
		GET UP 5
		LET OFF SOME STEAM BENNET c
	ENOUGH TALK
YOU HAVE BEEN TERMINATED
`

	f, err := fileFromString(src)
	if err != nil {
		t.Fatal(err)
	}
	defer func(f *os.File) {
		f.Close()
		os.Remove(f.Name())
	}(f)

	p := ArnoldC{input: f}
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("error parsing: %v", err)
	}

	expect := Program{
		Main: Function{
			Statements: []Statement{
				Expression{
					Instruction: "HEY CHRISTMAS TREE",
					Args:        []Value{VariableValue{"myVar"}, IntegerValue{10}},
				},
				Expression{
					Instruction: "TALK TO THE HAND",
					Args:        []Value{StringValue{"hello world"}},
				},
				Block{
					Instruction: "GET TO THE CHOPPER",
					Args:        []Value{VariableValue{"a"}},
					Statements: []Statement{
						Expression{
							Instruction: "HERE IS MY INVITATION",
							Args:        []Value{IntegerValue{4}},
						},
						Expression{
							Instruction: "GET UP",
							Args:        []Value{VariableValue{"b"}},
						},
						Expression{
							Instruction: "YOU'RE FIRED",
							Args:        []Value{IntegerValue{2}},
						},
					},
				},
				Block{
					Instruction: "GET TO THE CHOPPER",
					Args:        []Value{VariableValue{"a"}},
					Statements: []Statement{
						Expression{
							Instruction: "HERE IS MY INVITATION",
							Args:        []Value{VariableValue{"b"}},
						},
						Expression{
							Instruction: "GET UP",
							Args:        []Value{IntegerValue{5}},
						},
						Expression{
							Instruction: "LET OFF SOME STEAM BENNET",
							Args:        []Value{VariableValue{"c"}},
						},
					},
				},
			},
		},
	}

	if expect.String() != program.String() {
		t.Errorf("Program does not match expectations.\n  Expected:\n%s\n  Got:\n%s", expect, program)
	}
}
