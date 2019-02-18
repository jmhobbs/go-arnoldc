package arnoldc

import (
	"fmt"
	"io"
)

type interpreter struct {
	variables map[string]Value
}

func (i interpreter) Run(program *Program, stdout, stderr io.Writer) error {
	i.variables = make(map[string]Value)

	for _, expression := range program.Main.Expressions {
		switch expression.Instruction {
		case "TALK TO THE HAND":
			fmt.Fprintln(stdout, i.resolveValue(expression.Args[0]))
		case "HEY CHRISTMAS TREE":
			i.variables[expression.Args[0].Value().(string)] = expression.Args[1]
		default:
			return fmt.Errorf("runtime error; unknown instruction %q", expression.Instruction)
		}
	}
	return nil
}

func (i interpreter) resolveValue(v Value) interface{} {
	switch v.Type() {
	case VariableType:
		var varName string = v.Value().(string)
		value, ok := i.variables[varName]
		if !ok {
			// TODO
			panic("Well that isn't good.")
		}
		return value.Value()
	default:
		return v.Value()
	}
}
