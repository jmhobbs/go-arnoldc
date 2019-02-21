package arnoldc

import (
	"fmt"
)

// Program represents the ArnoldC program as a collection of Functions.
type Program struct {
	Main    Function
	Methods []Function
}

func (p Program) String() string {
	return fmt.Sprintf("ArnoldC(Main: %v, Methods: %v)", p.Main, p.Methods)
}

// A Function is either the Main function, or a user defined method.
type Function struct {
	// The name of the function, which is empty for Main
	Name string
	// Names of any arguments to this function
	Arguments []string
	// The actual lines to be executed in this function
	Statements []Statement
}

func (f Function) String() string {
	return fmt.Sprintf("Function(%q, Arguments: %v, Statements: %v)", f.Name, f.Arguments, f.Statements)
}
