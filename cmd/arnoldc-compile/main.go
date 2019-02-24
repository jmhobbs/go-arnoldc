package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	arnoldc "github.com/jmhobbs/go-arnoldc"
)

func init() {
	log.SetFlags(0)
}

func main() {
	var (
		binOut  string
		goOut   string
		verbose bool
	)

	flag.BoolVar(&verbose, "v", false, "Lots of debugging output.")
	flag.StringVar(&goOut, "g", "", "Intermediate go source file output. Defaults to a tempfile that is automatically removed.")
	flag.StringVar(&binOut, "o", "", "Output filename. Defaults to input file without arnoldc extension.")
	flag.Parse()

	sourcefile := flag.Arg(0)

	if binOut == "" {
		binOut = strings.TrimSuffix(filepath.Base(sourcefile), ".arnoldc")
	}

	if verbose {
		log.Println("Source ArnoldC file", sourcefile)
		log.Println("Writing binary to", binOut)
		if goOut != "" {
			log.Println("Writing intermediate go source to", goOut)
		} else {
			log.Println("Writing intermediate go source to a temp file.")
		}
	}

	f, err := os.Open(sourcefile)
	if err != nil {
		log.Fatalf("error opening program: %v", err)
	}
	defer f.Close()

	a := arnoldc.New(f)
	program, err := a.Parse()
	if err != nil {
		log.Fatalf("error parsing program: %v", err)
	}

	var out *os.File

	if goOut != "" {
		out, err = os.OpenFile(goOut, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("error opening go source output file: %v", err)
		}
		defer out.Close()
	} else {
		out, err = ioutil.TempFile("", fmt.Sprintf("*_%s.arnoldc.go", filepath.Base(sourcefile)))
		if err != nil {
			log.Fatalf("error opening tempfile: %v", err)
		}
		defer func() {
			out.Close()
			os.Remove(out.Name())
		}()
	}

	err = writePreamble(out)
	if err != nil {
		log.Fatal(err)
	}

	err = writeMain(out, program.Main)
	if err != nil {
		log.Fatal(err)
	}

	for _, method := range program.Methods {
		err = writeMethod(out, method)
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd := exec.Command("go", "build", "-o", binOut, out.Name())
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error compiling with go: %v", err)
		log.Fatalf("command output: %q", stdoutStderr)
	}
}

func writePreamble(out io.Writer) error {
	_, err := out.Write([]byte(`package main

// Generated by go-arnoldc, https://github.com/jmhobbs/go-arnoldc
// It's not meant to be pretty, it's meant to compile.

import(
	"fmt"
)
`))

	return err
}

func writeMain(out io.Writer, main arnoldc.Method) error {
	var err error
	if _, err = out.Write([]byte("func main() {\n")); err != nil {
		return err
	}

	if err = writeStatements(out, main.Statements); err != nil {
		return err
	}

	if _, err = out.Write([]byte("}\n\n")); err != nil {
		return err
	}

	return nil
}

func writeMethod(out io.Writer, method arnoldc.Method) error {
	var err error

	if _, err = fmt.Fprintf(out, "func %s(", method.Name); err != nil {
		return err
	}
	if len(method.Parameters) > 0 {
		if _, err = fmt.Fprintf(out, "%s int", strings.Join(method.Parameters, ", ")); err != nil {
			return err
		}
	}
	if _, err = fmt.Fprint(out, ") int {\n"); err != nil {
		return err
	}

	if err = writeStatements(out, method.Statements); err != nil {
		return err
	}

	// Always add a return to make things simpler.
	if _, err = fmt.Fprint(out, "return 0\n}\n\n"); err != nil {
		return err
	}
	return nil
}

func writeStatements(out io.Writer, statements []arnoldc.Statement) error {
	var err error

	for _, statement := range statements {
		if arnoldc.ExpressionType == statement.Type() {
			expression := statement.(arnoldc.Expression)
			switch expression.Instruction {
			case arnoldc.DECLARE:
				_, err = fmt.Fprintf(out, "var %s int = %d\n", expression.Args[0].Value().(string), expression.Args[1].Value().(int))
				if err != nil {
					return err
				}
			case arnoldc.PRINT:
				if _, err = fmt.Fprintf(out, "fmt.Println(%s)\n", goStringForValue(expression.Args[0])); err != nil {
					return err
				}
			case arnoldc.ASSIGN_FROM_CALL:
				returnName := expression.Args[0].Value().(string)
				methodName := expression.Args[1].Value().(string)
				if _, err := fmt.Fprintf(out, "%s = %s(", returnName, methodName); err != nil {
					return err
				}
				if len(expression.Args) > 2 {
					argStrings := []string{}
					for _, arg := range expression.Args[2:] {
						argStrings = append(argStrings, goStringForValue(arg))
					}
					if _, err := fmt.Fprintf(out, strings.Join(argStrings, ", ")); err != nil {
						return err
					}
				}
				if _, err := fmt.Fprintln(out, ")"); err != nil {
					return err
				}
			case arnoldc.CALL_METHOD:
				methodName := expression.Args[0].Value().(string)
				if _, err := fmt.Fprintf(out, "%s()\n", methodName); err != nil {
					return err
				}
			case arnoldc.RETURN:
				if _, err := fmt.Fprintf(out, "return %s\n", goStringForValue(expression.Args[0])); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown instruction %q", expression.Instruction)
			}
		} else {
			block := statement.(arnoldc.Block)

			switch block.Instruction {
			case arnoldc.ASSIGNMENT:
				v := block.Args[0]
				var operand string = goStringForValue(v)

				if v.Type() != arnoldc.VariableType {
					return fmt.Errorf("can not assign results to a non-variable")
				}
				statement := block.Statements[0]
				if arnoldc.ExpressionType != statement.Type() {
					return fmt.Errorf("illegal block inside assignment")
				}
				expression := statement.(arnoldc.Expression)
				if expression.Instruction != arnoldc.FIRST_OPERAND {
					return fmt.Errorf("variable assignment must start with a first operand")
				}
				if _, err = fmt.Fprintf(out, "%s = %s\n", operand, goStringForValue(expression.Args[0])); err != nil {
					return err
				}

				for _, statement := range block.Statements[1:] {
					if arnoldc.ExpressionType != statement.Type() {
						return fmt.Errorf("illegal block inside assignment: %q", statement.(arnoldc.Block).Instruction)
					}

					expression := statement.(arnoldc.Expression)

					var (
						oparg string = goStringForValue(expression.Args[0])
					)

					switch expression.Instruction {
					case arnoldc.ADD:
						if _, err = fmt.Fprintf(out, "%s = %s + %s\n", operand, operand, oparg); err != nil {
							return err
						}
					case arnoldc.SUBTRACT:
						if _, err = fmt.Fprintf(out, "%s = %s - %s\n", operand, operand, oparg); err != nil {
							return err
						}
					case arnoldc.MULTIPLY:
						if _, err = fmt.Fprintf(out, "%s = %s * %s\n", operand, operand, oparg); err != nil {
							return err
						}
					case arnoldc.DIVIDE:
						if _, err = fmt.Fprintf(out, "%s = %s / %s\n", operand, operand, oparg); err != nil {
							return err
						}
					case arnoldc.EQUAL_TO:
						if _, err = fmt.Fprintf(out, "if %s == %s { %s = 1 } else { %s = 0 }\n", operand, oparg, operand, operand); err != nil {
							return err
						}
					case arnoldc.GREATER_THAN:
						if _, err = fmt.Fprintf(out, "if %s > %s { %s = 1 } else { %s = 0 }\n", operand, oparg, operand, operand); err != nil {
							return err
						}
					case arnoldc.OR:
						if _, err = fmt.Fprintf(out, "if %s != 0 || %s != 0 { %s = 1 } else { %s = 0 }\n", operand, oparg, operand, operand); err != nil {
							return err
						}
					case arnoldc.AND:
						if _, err = fmt.Fprintf(out, "if %s !=0 &&  %s != 0 { %s = 1 } else { %s = 0 }\n", operand, oparg, operand, operand); err != nil {
							return err
						}
					default:
						return fmt.Errorf("illegal statement inside assignment: %q", expression.Instruction)
					}
				}

			case arnoldc.IF:
				if _, err = fmt.Fprintf(out, "if %s !=0 {\n", goStringForValue(block.Args[0])); err != nil {
					return err
				}
				if err = writeStatements(out, block.Statements[0].(arnoldc.Block).Statements); err != nil {
					return err
				}
				if len(block.Statements) > 1 {
					if _, err = fmt.Fprintf(out, "} else {\n"); err != nil {
						return err
					}
					if err = writeStatements(out, block.Statements[1].(arnoldc.Block).Statements); err != nil {
						return err
					}
				}
				if _, err = fmt.Fprintf(out, "}\n"); err != nil {
					return err
				}
				/*
					case arnoldc.WHILE:
						if _, err := i.whileBlock(block, vars); err != nil {
							return 0, err
						}
				*/
			default:
				return fmt.Errorf("runtime error; unknown block instruction %q", block.Instruction)
			}
		}
	}

	return nil
}

func goStringForValue(v arnoldc.Value) string {
	switch v.Type() {
	case arnoldc.IntegerType:
		return fmt.Sprintf("%d", v.Value().(int))
	case arnoldc.VariableType:
		return v.Value().(string)
	case arnoldc.StringType:
		return fmt.Sprintf("%q", v.Value())
	}
	panic("unknown value type")
}
