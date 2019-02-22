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
	_, err := i.invokeFunction(&program.Main, []arnoldc.Value{}, nil)
	return err
}

func (i *Interpreter) invokeFunction(f *arnoldc.Function, arguments []arnoldc.Value, parentScope *scope) (int, error) {
	if len(arguments) != len(f.Arguments) {
		return 0, fmt.Errorf("incorrect number of arguments for %q", f.Name)
	}

	vars := newScope(parentScope)
	for i, name := range f.Arguments {
		v, err := parentScope.Get(arguments[i])
		if err != nil {
			return 0, fmt.Errorf("invalid argument; %v", arguments[i])
		}
		vars.Set(name, v)
	}

	err := i.executeStatements(f.Statements, vars)
	if err != nil {
		return 0, err
	}

	// TODO: I do not like doing it like this.
	lastStatement := f.Statements[len(f.Statements)-1]
	if arnoldc.ExpressionType == lastStatement.Type() && lastStatement.(arnoldc.Expression).Instruction == "I'LL BE BACK" {
		expression := lastStatement.(arnoldc.Expression)

		v, err := vars.Get(expression.Args[0])
		if err != nil {
			return 0, fmt.Errorf("invalid return value; %v", err)
		}

		return v, nil
	}

	return 0, nil
}

func (i *Interpreter) executeStatements(statements []arnoldc.Statement, vars *scope) error {
	for _, statement := range statements {
		if arnoldc.ExpressionType == statement.Type() {
			expression := statement.(arnoldc.Expression)

			switch expression.Instruction {
			case "TALK TO THE HAND":
				value, err := vars.resolveValue(expression.Args[0])
				if err != nil {
					return fmt.Errorf("runtime error; %v", err)
				}
				fmt.Fprintln(i.stdout, value)
			case "HEY CHRISTMAS TREE":
				vars.Set(expression.Args[0].Value().(string), expression.Args[1].Value().(int))
			case "GET YOUR ASS TO MARS":
				returnName := expression.Args[0].Value().(string)
				methodName := expression.Args[1].Value().(string)
				function, ok := i.method(methodName)
				if !ok {
					return fmt.Errorf("unknown method; %q", methodName)
				}
				ret, err := i.invokeFunction(function, expression.Args[2:], newScope(vars))
				if err != nil {
					return fmt.Errorf("runtime err; %v", err)
				}
				vars.Set(returnName, ret)
			case "DO IT NOW":
				methodName := expression.Args[0].Value().(string)
				function, ok := i.method(methodName)
				if !ok {
					return fmt.Errorf("unknown method; %q", methodName)
				}
				_, err := i.invokeFunction(function, expression.Args[1:], newScope(vars))
				if err != nil {
					return fmt.Errorf("runtime err; %v", err)
				}
			case "I'LL BE BACK":
				// NO-OP Handled in invokeFunction
				continue
			default:
				return fmt.Errorf("runtime error; unknown instruction %q", expression.Instruction)
			}
		} else {
			block := statement.(arnoldc.Block)

			switch block.Instruction {
			case "GET TO THE CHOPPER":
				if err := i.assigmentBlock(block, vars); err != nil {
					return err
				}
			case "BECAUSE I'M GOING TO SAY PLEASE":
				if err := i.ifElseBlock(block, vars); err != nil {
					return err
				}
			case "STICK AROUND":
				if err := i.whileBlock(block, vars); err != nil {
					return err
				}
			default:
				return fmt.Errorf("runtime error; unknown block instruction %q", block.Instruction)
			}
		}
	}
	return nil
}

// Execute and return an Assignment Block
func (i *Interpreter) assigmentBlock(block arnoldc.Block, parentScope *scope) error {
	vars := newScope(parentScope)

	v := block.Args[0]

	if v.Type() != arnoldc.VariableType {
		return fmt.Errorf("can not assign results to a non-variable")
	}

	// TODO: Setup Instruction
	statement := block.Statements[0]
	if arnoldc.ExpressionType != statement.Type() {
		return fmt.Errorf("illegal block inside assignment")
	}
	expression := statement.(arnoldc.Expression)
	if expression.Instruction != "HERE IS MY INVITATION" {
		return fmt.Errorf("variable assignment must start with a first operand")
	}

	x, err := vars.Get(expression.Args[0])
	if err != nil {
		return err
	}

	for _, statement := range block.Statements[1:] {
		if arnoldc.ExpressionType != statement.Type() {
			return fmt.Errorf("illegal block inside assignment: %q", statement.(arnoldc.Block).Instruction)
		}

		expression := statement.(arnoldc.Expression)

		// All of these should have exactly one argument.
		arg, err := vars.Get(expression.Args[0])
		if err != nil {
			return err
		}

		switch expression.Instruction {
		case "GET UP":
			x = x + arg
		case "GET DOWN":
			x = x - arg
		case "YOU'RE FIRED":
			x = x * arg
		case "HE HAD TO SPLIT":
			x = x / arg
		case "YOU ARE NOT YOU YOU ARE ME":
			if x == arg {
				x = TRUE
			} else {
				x = FALSE
			}
		case "LET OFF SOME STEAM BENNET":
			if x > arg {
				x = TRUE
			} else {
				x = FALSE
			}
		case "CONSIDER THAT A DIVORCE":
			if x != FALSE || arg != FALSE {
				x = TRUE
			} else {
				x = FALSE
			}
		case "KNOCK KNOCK":
			if x != FALSE && arg != FALSE {
				x = TRUE
			} else {
				x = FALSE
			}
		default:
			return fmt.Errorf("illegal statement inside assignment: %q", expression.Instruction)
		}
	}

	parentScope.Set(v.Value().(string), x)

	return nil
}

func (i *Interpreter) ifElseBlock(block arnoldc.Block, vars *scope) error {
	v, err := vars.Get(block.Args[0])
	if err != nil {
		return err
	}

	if v != FALSE {
		ifBlock := block.Statements[0].(arnoldc.Block)
		return i.executeStatements(ifBlock.Statements, vars)
	} else if len(block.Statements) > 1 {
		elseBlock := block.Statements[1].(arnoldc.Block)
		return i.executeStatements(elseBlock.Statements, vars)
	}

	return nil
}

func (i *Interpreter) whileBlock(block arnoldc.Block, parentScope *scope) error {
	vars := newScope(parentScope)

	for {
		v, err := vars.Get(block.Args[0])
		if err != nil {
			return err
		}
		if v == FALSE {
			break
		}

		err = i.executeStatements(block.Statements, vars)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) method(name string) (*arnoldc.Function, bool) {
	for _, f := range i.program.Methods {
		if f.Name == name {
			return &f, true
		}
	}
	return nil, false
}
