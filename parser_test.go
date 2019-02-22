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
				// TODO: If & if/else should probably be a special type, since it's a compound block statement.
				Block{
					Instruction: "BECAUSE I'M GOING TO SAY PLEASE",
					Args:        []Value{IntegerValue{0}},
					Statements: []Statement{
						Block{
							Instruction: "__TRUE",
							Args:        []Value{},
							Statements: []Statement{
								Expression{
									Instruction: "TALK TO THE HAND",
									Args:        []Value{StringValue{"false is true?!"}},
								},
							},
						},
						Block{
							Instruction: "__FALSE",
							Args:        []Value{},
							Statements: []Statement{
								Expression{
									Instruction: "TALK TO THE HAND",
									Args:        []Value{StringValue{"false is not true"}},
								},
							},
						},
					},
				},
			},
		},
		Methods: []Function{
			Function{
				Name:      "hello",
				Arguments: []string{},
				Statements: []Statement{
					Expression{
						Instruction: "TALK TO THE HAND",
						Args:        []Value{StringValue{"hello"}},
					},
				},
			},
			Function{
				Name:      "double",
				Arguments: []string{"number"},
				Statements: []Statement{
					Block{
						Instruction: "GET TO THE CHOPPER",
						Args:        []Value{VariableValue{"result"}},
						Statements: []Statement{
							Expression{
								Instruction: "HERE IS MY INVITATION",
								Args:        []Value{VariableValue{"double"}},
							},
							Expression{
								Instruction: "YOU'RE FIRED",
								Args:        []Value{IntegerValue{2}},
							},
						},
					},
					Expression{
						Instruction: "I'LL BE BACK",
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
