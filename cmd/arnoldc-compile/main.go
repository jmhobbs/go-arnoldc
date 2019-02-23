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
		fmt.Fprintf(os.Stderr, "error opening program: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	a := arnoldc.New(f)
	program, err := a.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing program: %v", err)
		os.Exit(1)
	}

	var out *os.File

	if goOut != "" {
		out, err = os.OpenFile(goOut, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening go source output file: %v", err)
			os.Exit(1)
		}
		defer out.Close()
	} else {
		out, err = ioutil.TempFile("", fmt.Sprintf("*_%s.arnoldc.go", filepath.Base(sourcefile)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening tempfile: %v", err)
			os.Exit(1)
		}
		defer func() {
			out.Close()
			os.Remove(out.Name())
		}()
	}

	err = writePreamble(out)
	if err != nil {
		panic(err)
	}

	err = writeMain(out, program.Main)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "build", "-o", binOut, out.Name())
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error compiling with go: %v", err)
		fmt.Fprintf(os.Stderr, "command output: %q", stdoutStderr)
		os.Exit(1)
	}
}

func writePreamble(out io.Writer) error {
	_, err := out.Write([]byte(`package main

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

	for _, statement := range main.Statements {
		if arnoldc.ExpressionType == statement.Type() {
			expression := statement.(arnoldc.Expression)
			switch expression.Instruction {
			case arnoldc.DECLARE:
				_, err = fmt.Fprintf(out, "var %s int = %d\n", expression.Args[0].Value().(string), expression.Args[1].Value().(int))
				if err != nil {
					return err
				}
			case arnoldc.PRINT:
				switch expression.Args[0].Type() {
				case arnoldc.VariableType:
					_, err = fmt.Fprintf(out, "fmt.Println(%s)\n", expression.Args[0].Value())
				case arnoldc.IntegerType:
					_, err = fmt.Fprintf(out, "fmt.Println(%d)\n", expression.Args[0].Value())
				case arnoldc.StringType:
					_, err = fmt.Fprintf(out, "fmt.Println(%q)\n", expression.Args[0].Value())
				}
				if err != nil {
					return err
				}
			}
		}
	}

	if _, err = out.Write([]byte("}")); err != nil {
		return err
	}

	return nil
}
