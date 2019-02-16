package arnoldc

import (
	"fmt"
	"io"
)

// Program represents the ArnoldC program as a collection of Functions.
type Program struct {
	Main    Function
	Methods []Function
}

func (p Program) String() string {
	return fmt.Sprintf("ArnoldC(Main: %v, Methods: %v)", p.Main, p.Methods)
}

func (p Program) Run(stdout, stderr io.Writer) error {
	i := interpreter{}
	return i.Run(&p, stdout, stderr)
}

// A Function is either the Main function, or a user defined method.
type Function struct {
	// The name of the function, which is empty for Main
	Name string
	// Names of any arguments to this function
	Arguments []string
	// The actual lines to be executed in this function
	Expressions []Expression
}

func (f Function) String() string {
	return fmt.Sprintf("Function(%q, Arguments: %v, Expressions: %v)", f.Name, f.Arguments, f.Expressions)
}

// An Expression is a single line of an ArnoldC program.  It is an instruction and any parameters to that instruction.
type Expression struct {
	Instruction string // Or Token?
	Args        []Value
}

func (e Expression) String() string {
	return fmt.Sprintf("Expression(%q, %v)", e.Instruction, e.Args)
}
