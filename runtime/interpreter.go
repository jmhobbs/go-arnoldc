package runtime

import (
	"fmt"
	"io"

	arnoldc "github.com/jmhobbs/go-arnoldc"
)

const (
	TRUE  int = 1
	FALSE int = 0
)

type Interpreter struct {
	variables map[string]arnoldc.Value
	stdout    io.Writer
	stderr    io.Writer
}

func New(stdout, stderr io.Writer) *Interpreter {
	return &Interpreter{
		variables: make(map[string]arnoldc.Value),
		stdout:    stdout,
		stderr:    stderr,
	}
}

func (i *Interpreter) Run(program *arnoldc.Program) error {
	i.executeStatements(program.Main.Statements)
	return nil
}

func (i *Interpreter) executeStatements(statements []arnoldc.Statement) error {
	for _, statement := range statements {
		if arnoldc.ExpressionType == statement.Type() {
			expression := statement.(arnoldc.Expression)

			switch expression.Instruction {
			case "TALK TO THE HAND":
				value, err := i.resolveValue(expression.Args[0])
				if err != nil {
					return fmt.Errorf("runtime error; %v", err)
				}
				fmt.Fprintln(i.stdout, value)
			case "HEY CHRISTMAS TREE":
				i.variables[expression.Args[0].Value().(string)] = expression.Args[1]
			default:
				return fmt.Errorf("runtime error; unknown instruction %q", expression.Instruction)
			}
		} else {
			block := statement.(arnoldc.Block)

			switch block.Instruction {
			case "GET TO THE CHOPPER":
				if err := i.assigmentBlock(block); err != nil {
					return err
				}
			case "BECAUSE I'M GOING TO SAY PLEASE":
				if err := i.ifElseBlock(block); err != nil {
					return err
				}
			default:
				return fmt.Errorf("runtime error; unknown block instruction %q", block.Instruction)
			}
		}
	}
	return nil
}

// Resolve a value to it's underlying type, following variable references.
func (i *Interpreter) resolveValue(v arnoldc.Value) (interface{}, error) {
	if v.Type() == arnoldc.VariableType {
		var varName string = v.Value().(string)
		value, ok := i.variables[varName]
		if !ok {
			return nil, fmt.Errorf("undefined variable %q", varName)
		}
		return value.Value(), nil
	}
	return v.Value(), nil
}

// Resolve a value to it's underlying integer, following variable references.
func (i *Interpreter) resolveNumber(v arnoldc.Value) (int, error) {
	switch v.Type() {
	case arnoldc.VariableType:
		var varName string = v.Value().(string)
		value, ok := i.variables[varName]
		if !ok {
			return 0, fmt.Errorf("undefined variable %q", varName)
		}
		return value.Value().(int), nil
	case arnoldc.IntegerType:
		return v.Value().(int), nil
	default:
		return 0, fmt.Errorf("invalid value for number type")
	}
}

// Execute and return an Assignment Block
func (i *Interpreter) assigmentBlock(block arnoldc.Block) error {
	v := block.Args[0]

	if v.Type() != arnoldc.VariableType {
		return fmt.Errorf("can not assign results to a non-variable")
	}

	// TODO: Setup Instruction
	statement := block.Statements[0]
	if arnoldc.ExpressionType != statement.Type() {
		return fmt.Errorf("illegal block inside assignment")
	}
	expression := statement.(arnoldc.Expression)
	if expression.Instruction != "HERE IS MY INVITATION" {
		return fmt.Errorf("variable assignment must start with a first operand")
	}

	x, err := i.resolveNumber(expression.Args[0])
	if err != nil {
		return err
	}

	for _, statement := range block.Statements[1:] {
		if arnoldc.ExpressionType != statement.Type() {
			return fmt.Errorf("illegal block inside assignment: %q", statement.(arnoldc.Block).Instruction)
		}

		expression := statement.(arnoldc.Expression)

		// All of these should have exactly one argument.
		arg, err := i.resolveNumber(expression.Args[0])
		if err != nil {
			return err
		}

		switch expression.Instruction {
		case "GET UP":
			x = x + arg
		case "GET DOWN":
			x = x - arg
		case "YOU'RE FIRED":
			x = x * arg
		case "HE HAD TO SPLIT":
			x = x / arg
		case "YOU ARE NOT YOU YOU ARE ME":
			if x == arg {
				x = TRUE
			} else {
				x = FALSE
			}
		case "LET OFF SOME STEAM BENNET":
			if x > arg {
				x = TRUE
			} else {
				x = FALSE
			}
		case "CONSIDER THAT A DIVORCE":
			if x != FALSE || arg != FALSE {
				x = TRUE
			} else {
				x = FALSE
			}
		case "KNOCK KNOCK":
			if x != FALSE && arg != FALSE {
				x = TRUE
			} else {
				x = FALSE
			}
		default:
			return fmt.Errorf("illegal statement inside assignment: %q", expression.Instruction)
		}
	}

	i.variables[v.Value().(string)] = arnoldc.NewIntegerValue(x)

	return nil
}

func (i *Interpreter) ifElseBlock(block arnoldc.Block) error {
	v, err := i.resolveNumber(block.Args[0])
	if err != nil {
		return err
	}

	if v != FALSE {
		ifBlock := block.Statements[0].(arnoldc.Block)
		return i.executeStatements(ifBlock.Statements)
	} else if len(block.Statements) > 1 {
		elseBlock := block.Statements[1].(arnoldc.Block)
		return i.executeStatements(elseBlock.Statements)
	}

	return nil
}
