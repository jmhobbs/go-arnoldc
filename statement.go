package arnoldc

import (
	"fmt"
)

// A Statement is either an Expression, or a Block
type Statement interface {
	Type() StatementType
}

type StatementType int

const (
	ExpressionType StatementType = iota
	BlockType
)

// An Expression is a single line of an ArnoldC program.  It is an instruction and any parameters to that instruction.
type Expression struct {
	Instruction string // Or Token?
	Args        []Value
}

func (e Expression) Type() StatementType {
	return ExpressionType
}

func (e Expression) String() string {
	return fmt.Sprintf("Expression(%q, %v)", e.Instruction, e.Args)
}

// A block is a group of expressions, like an if statement or an assignment
type Block struct {
	Instruction string
	Args        []Value
	Statements  []Statement
}

func (b Block) Type() StatementType {
	return BlockType
}
