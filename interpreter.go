package arnoldc

import (
	"fmt"
	"io"
)

type interpreter struct {
	variables map[string]Value
}

func (i interpreter) Run(program *Program, stdout, stderr io.Writer) error {
	for _, expression := range program.Main.Expressions {
		switch expression.Instruction {
		case "TALK TO THE HAND":
			fmt.Fprintf(stdout, "%s", i.resolveValue(expression.Args[0]))
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
		return value
	default:
		return v.Value()
	}
}
