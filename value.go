package arnoldc

import "fmt"

type ValueType int

const (
	VariableType ValueType = iota
	StringType
	IntegerType
)

type Value interface {
	Type() ValueType
	Value() interface{}
}

type StringValue struct {
	str string
}

func (v StringValue) Type() ValueType       { return StringType }
func (v StringValue) String() string        { return fmt.Sprintf("String(%q)", v.str) }
func (v StringValue) Value() interface{}    { return v.str }
func NewStringValue(str string) StringValue { return StringValue{str} }

type VariableValue struct {
	name string
}

func (v VariableValue) Type() ValueType          { return VariableType }
func (v VariableValue) String() string           { return fmt.Sprintf("Variable(%q)", v.name) }
func (v VariableValue) Value() interface{}       { return v.name }
func NewVariableValue(name string) VariableValue { return VariableValue{name} }

type IntegerValue struct {
	d int
}

func (v IntegerValue) Type() ValueType    { return IntegerType }
func (v IntegerValue) String() string     { return fmt.Sprintf("Integer(%d)", v.d) }
func (v IntegerValue) Value() interface{} { return v.d }
func NewIntegerValue(d int) IntegerValue  { return IntegerValue{d} }
