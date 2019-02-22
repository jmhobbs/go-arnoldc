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

	BECAUSE I'M GOING TO SAY PLEASE @I LIED
		TALK TO THE HAND "false is true?!"
	BULLSHIT
		TALK TO THE HAND "false is not true"
	YOU HAVE NO RESPECT FOR LOGIC
YOU HAVE BEEN TERMINATED

LISTEN TO ME VERY CAREFULLY hello
	TALK TO THE HAND "hello"
HASTA LA VISTA, BABY

LISTEN TO ME VERY CAREFULLY double
	I NEED YOUR CLOTHES YOUR BOOTS AND YOUR MOTORCYCLE number
	GIVE THESE PEOPLE AIR

	GET TO THE CHOPPER result
		HERE IS MY INVITATION double
		YOU'RE FIRED 2
	ENOUGH TALK

	I'LL BE BACK result
HASTA LA VISTA, BABY
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
		Main: Method{
			Statements: []Statement{
				Expression{
					Instruction: DECLARE,
					Args:        []Value{VariableValue{"myVar"}, IntegerValue{10}},
				},
				Expression{
					Instruction: PRINT,
					Args:        []Value{StringValue{"hello world"}},
				},
				Block{
					Instruction: ASSIGNMENT,
					Args:        []Value{VariableValue{"a"}},
					Statements: []Statement{
						Expression{
							Instruction: FIRST_OPERAND,
							Args:        []Value{IntegerValue{4}},
						},
						Expression{
							Instruction: ADD,
							Args:        []Value{VariableValue{"b"}},
						},
						Expression{
							Instruction: MULTIPLY,
							Args:        []Value{IntegerValue{2}},
						},
					},
				},
				Block{
					Instruction: ASSIGNMENT,
					Args:        []Value{VariableValue{"a"}},
					Statements: []Statement{
						Expression{
							Instruction: FIRST_OPERAND,
							Args:        []Value{VariableValue{"b"}},
						},
						Expression{
							Instruction: ADD,
							Args:        []Value{IntegerValue{5}},
						},
						Expression{
							Instruction: GREATER_THAN,
							Args:        []Value{VariableValue{"c"}},
						},
					},
				},
				// TODO: If & if/else should probably be a special type, since it's a compound block statement.
				Block{
					Instruction: IF,
					Args:        []Value{IntegerValue{0}},
					Statements: []Statement{
						Block{
							Instruction: IF_TRUE,
							Args:        []Value{},
							Statements: []Statement{
								Expression{
									Instruction: PRINT,
									Args:        []Value{StringValue{"false is true?!"}},
								},
							},
						},
						Block{
							Instruction: IF_FALSE,
							Args:        []Value{},
							Statements: []Statement{
								Expression{
									Instruction: PRINT,
									Args:        []Value{StringValue{"false is not true"}},
								},
							},
						},
					},
				},
			},
		},
		Methods: []Method{
			Method{
				Name: "hello",
				Statements: []Statement{
					Expression{
						Instruction: PRINT,
						Args:        []Value{StringValue{"hello"}},
					},
				},
			},
			Method{
				Name:       "double",
				Parameters: []string{"number"},
				Statements: []Statement{
					Block{
						Instruction: ASSIGNMENT,
						Args:        []Value{VariableValue{"result"}},
						Statements: []Statement{
							Expression{
								Instruction: FIRST_OPERAND,
								Args:        []Value{VariableValue{"double"}},
							},
							Expression{
								Instruction: MULTIPLY,
								Args:        []Value{IntegerValue{2}},
							},
						},
					},
					Expression{
						Instruction: RETURN,
						Args:        []Value{VariableValue{"result"}},
					},
				},
			},
		},
	}

	if expect.String() != program.String() {
		t.Errorf("Program does not match expectations.\n  Expected:\n%s\n  Got:\n%s", expect, program)
	}
}
