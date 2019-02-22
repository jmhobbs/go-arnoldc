package arnoldc

import (
	"fmt"
)

// Program represents the ArnoldC program as a collection of Methods.
type Program struct {
	Main    Method
	Methods []Method
}

func (p Program) String() string {
	return fmt.Sprintf("ArnoldC(Main: %v, Methods: %v)", p.Main, p.Methods)
}

// A Method is either the Main function, or a user defined method.
type Method struct {
	// The name of the function, which is empty for Main
	Name string
	// Names of any parameters to this function
	Parameters []string
	// The actual lines to be executed in this function
	Statements []Statement
}

func (f Method) String() string {
	return fmt.Sprintf("Method(%q, Parameters: %v, Statements: %v)", f.Name, f.Parameters, f.Statements)
}
