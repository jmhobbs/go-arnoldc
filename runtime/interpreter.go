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
	stdout  io.Writer
	stderr  io.Writer
	program *arnoldc.Program
}

func New(stdout, stderr io.Writer) *Interpreter {
	return &Interpreter{
		stdout: stdout,
		stderr: stderr,
	}
}

func (i *Interpreter) Run(program *arnoldc.Program) error {
	i.program = program
	_, err := i.invokeMethod(&program.Main, []arnoldc.Value{}, nil)
	return err
}

func (i *Interpreter) invokeMethod(f *arnoldc.Method, arguments []arnoldc.Value, parentScope *scope) (int, error) {
	if len(arguments) != len(f.Parameters) {
		return 0, fmt.Errorf("incorrect number of arguments for %q", f.Name)
	}

	vars := newScope(nil)
	for i, name := range f.Parameters {
		v, err := parentScope.Get(arguments[i])
		if err != nil {
			return 0, fmt.Errorf("invalid argument; %v", arguments[i])
		}
		vars.Set(name, v)
	}

	return i.executeStatements(f.Statements, vars)
}

func (i *Interpreter) executeStatements(statements []arnoldc.Statement, vars *scope) (int, error) {
	for _, statement := range statements {
		if arnoldc.ExpressionType == statement.Type() {
			expression := statement.(arnoldc.Expression)

			switch expression.Instruction {
			case arnoldc.PRINT:
				value, err := vars.resolveValue(expression.Args[0])
				if err != nil {
					return 0, fmt.Errorf("runtime error; %v", err)
				}
				fmt.Fprintln(i.stdout, value)
			case arnoldc.DECLARE:
				vars.Set(expression.Args[0].Value().(string), expression.Args[1].Value().(int))
			case arnoldc.ASSIGN_FROM_CALL:
				returnName := expression.Args[0].Value().(string)
				methodName := expression.Args[1].Value().(string)
				function, ok := i.method(methodName)
				if !ok {
					return 0, fmt.Errorf("unknown method; %q", methodName)
				}
				ret, err := i.invokeMethod(function, expression.Args[2:], vars)
				if err != nil {
					return 0, fmt.Errorf("runtime err; %v", err)
				}
				vars.Set(returnName, ret)
			case arnoldc.CALL_METHOD:
				methodName := expression.Args[0].Value().(string)
				function, ok := i.method(methodName)
				if !ok {
					return 0, fmt.Errorf("unknown method; %q", methodName)
				}
				_, err := i.invokeMethod(function, expression.Args[1:], vars)
				if err != nil {
					return 0, fmt.Errorf("runtime error; %v", err)
				}
			case arnoldc.RETURN:
				v, err := vars.Get(expression.Args[0])
				if err != nil {
					return 0, fmt.Errorf("invalid return value; %v", err)
				}
				return v, nil
			default:
				return 0, fmt.Errorf("runtime error; unknown instruction %q", expression.Instruction)
			}
		} else {
			block := statement.(arnoldc.Block)

			switch block.Instruction {
			case arnoldc.ASSIGNMENT:
				if _, err := i.assigmentBlock(block, vars); err != nil {
					return 0, err
				}
			case arnoldc.IF:
				if _, err := i.ifElseBlock(block, vars); err != nil {
					return 0, err
				}
			case arnoldc.WHILE:
				if _, err := i.whileBlock(block, vars); err != nil {
					return 0, err
				}
			default:
				return 0, fmt.Errorf("runtime error; unknown block instruction %q", block.Instruction)
			}
		}
	}
	return 0, nil
}

// Execute and return an Assignment Block
func (i *Interpreter) assigmentBlock(block arnoldc.Block, parentScope *scope) (int, error) {
	vars := newScope(parentScope)

	v := block.Args[0]

	if v.Type() != arnoldc.VariableType {
		return 0, fmt.Errorf("can not assign results to a non-variable")
	}

	// TODO: Setup Instruction
	statement := block.Statements[0]
	if arnoldc.ExpressionType != statement.Type() {
		return 0, fmt.Errorf("illegal block inside assignment")
	}
	expression := statement.(arnoldc.Expression)
	if expression.Instruction != arnoldc.FIRST_OPERAND {
		return 0, fmt.Errorf("variable assignment must start with a first operand")
	}

	x, err := vars.Get(expression.Args[0])
	if err != nil {
		return 0, err
	}

	for _, statement := range block.Statements[1:] {
		if arnoldc.ExpressionType != statement.Type() {
			return 0, fmt.Errorf("illegal block inside assignment: %q", statement.(arnoldc.Block).Instruction)
		}

		expression := statement.(arnoldc.Expression)

		// All of these should have exactly one argument.
		arg, err := vars.Get(expression.Args[0])
		if err != nil {
			return 0, err
		}

		switch expression.Instruction {
		case arnoldc.ADD:
			x = x + arg
		case arnoldc.SUBTRACT:
			x = x - arg
		case arnoldc.MULTIPLY:
			x = x * arg
		case arnoldc.DIVIDE:
			x = x / arg
		case arnoldc.EQUAL_TO:
			if x == arg {
				x = TRUE
			} else {
				x = FALSE
			}
		case arnoldc.GREATER_THAN:
			if x > arg {
				x = TRUE
			} else {
				x = FALSE
			}
		case arnoldc.OR:
			if x != FALSE || arg != FALSE {
				x = TRUE
			} else {
				x = FALSE
			}
		case arnoldc.AND:
			if x != FALSE && arg != FALSE {
				x = TRUE
			} else {
				x = FALSE
			}
		default:
			return 0, fmt.Errorf("illegal statement inside assignment: %q", expression.Instruction)
		}
	}

	parentScope.Set(v.Value().(string), x)

	return 0, nil
}

func (i *Interpreter) ifElseBlock(block arnoldc.Block, vars *scope) (int, error) {
	v, err := vars.Get(block.Args[0])
	if err != nil {
		return 0, err
	}

	if v != FALSE {
		ifBlock := block.Statements[0].(arnoldc.Block)
		return i.executeStatements(ifBlock.Statements, vars)
	} else if len(block.Statements) > 1 {
		elseBlock := block.Statements[1].(arnoldc.Block)
		return i.executeStatements(elseBlock.Statements, vars)
	}

	return 0, nil
}

func (i *Interpreter) whileBlock(block arnoldc.Block, parentScope *scope) (int, error) {
	vars := newScope(parentScope)

	for {
		v, err := vars.Get(block.Args[0])
		if err != nil {
			return 0, err
		}
		if v == FALSE {
			break
		}

		_, err = i.executeStatements(block.Statements, vars)

		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

func (i *Interpreter) method(name string) (*arnoldc.Method, bool) {
	for _, f := range i.program.Methods {
		if f.Name == name {
			return &f, true
		}
	}
	return nil, false
}
