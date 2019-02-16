package arnoldc

import "fmt"

type ValueType int

// TODO: Or use the int types defined in parser.go?
const (
	LiteralType ValueType = iota
	VariableType
	StringType
	IntegerType
	BoolType
)

type Value interface {
	Type() ValueType
	Value() interface{}
}

type StringValue struct {
	str string
}

func (v StringValue) Type() ValueType    { return StringType }
func (v StringValue) String() string     { return fmt.Sprintf("String(%q)", v.str) }
func (v StringValue) Value() interface{} { return v.str }

type VariableValue struct {
	name string
}

func (v VariableValue) Type() ValueType    { return VariableType }
func (v VariableValue) String() string     { return fmt.Sprintf("Variable(%q)", v.name) }
func (v VariableValue) Value() interface{} { return v.name }

type IntegerValue struct {
	v int
}

func (v IntegerValue) Type() ValueType    { return IntegerType }
func (v IntegerValue) String() string     { return fmt.Sprintf("Integer(%d)", v.v) }
func (v IntegerValue) Value() interface{} { return v.v }

type BoolValue struct {
	b bool
}

func (v BoolValue) Type() ValueType { return BoolType }
func (v BoolValue) String() string {
	if v.b {
		return "TRUE"
	}
	return "FALSE"
}
func (v BoolValue) Value() interface{} { return v.b }
