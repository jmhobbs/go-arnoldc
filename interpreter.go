package arnoldc

import (
	"fmt"
	"io"
)

type interpreter struct {
	variables map[string]Value
}

func (i *interpreter) Run(program *Program, stdout, stderr io.Writer) error {
	i.variables = make(map[string]Value)

	for _, statement := range program.Main.Statements {
		if ExpressionType == statement.Type() {
			expression := statement.(Expression)

			switch expression.Instruction {
			case "TALK TO THE HAND":
				value, err := i.resolveValue(expression.Args[0])
				if err != nil {
					return fmt.Errorf("runtime error; %v", err)
				}
				fmt.Fprintln(stdout, value)
			case "HEY CHRISTMAS TREE":
				i.variables[expression.Args[0].Value().(string)] = expression.Args[1]
			default:
				return fmt.Errorf("runtime error; unknown instruction %q", expression.Instruction)
			}
		} else {
			// TODO: OH SHIT ITS A BLOCK BOYS
			block := statement.(Block)

			switch block.Instruction {
			case "GET TO THE CHOPPER":
				err := i.assigmentBlock(block)
				if err != nil {
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
func (i *interpreter) resolveValue(v Value) (interface{}, error) {
	if v.Type() == VariableType {
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
func (i *interpreter) resolveNumber(v Value) (int, error) {
	switch v.Type() {
	case VariableType:
		var varName string = v.Value().(string)
		value, ok := i.variables[varName]
		if !ok {
			return 0, fmt.Errorf("undefined variable %q", varName)
		}
		return value.Value().(int), nil
	case IntegerType:
		return v.Value().(int), nil
	default:
		return 0, fmt.Errorf("invalid value for number type")
	}
}

func (i *interpreter) assigmentBlock(block Block) error {
	v := block.Args[0]

	if v.Type() != VariableType {
		return fmt.Errorf("runtime error; can not assign results to a non-variable")
	}

	// TODO: Setup Instruction
	statement := block.Statements[0]
	if ExpressionType != statement.Type() {
		return fmt.Errorf("runtime error; illegal block inside assignment")
	}
	expression := statement.(Expression)
	if expression.Instruction != "HERE IS MY INVITATION" {
		return fmt.Errorf("runtime error; variable assignment must start with a first operand")
	}

	x, err := i.resolveNumber(expression.Args[0])
	if err != nil {
		return err
	}

	for _, statement := range block.Statements[1:] {
		if ExpressionType != statement.Type() {
			return fmt.Errorf("runtime error; illegal block inside assignment: %q", statement.(Block).Instruction)
		}

		expression := statement.(Expression)

		switch expression.Instruction {
		case "GET UP":
			arg, err := i.resolveNumber(expression.Args[0])
			if err != nil {
				return err
			}
			x = x + arg
		case "GET DOWN":
			arg, err := i.resolveNumber(expression.Args[0])
			if err != nil {
				return err
			}
			x = x - arg
		case "YOU'RE FIRED":
			arg, err := i.resolveNumber(expression.Args[0])
			if err != nil {
				return err
			}
			x = x * arg
		case "HE HAD TO SPLIT":
			arg, err := i.resolveNumber(expression.Args[0])
			if err != nil {
				return err
			}
			x = x / arg
		default:
			return fmt.Errorf("runtime error; illegal statement inside assignment: %q", expression.Instruction)
		}
	}

	i.variables[v.Value().(string)] = IntegerValue{x}

	return nil
}
