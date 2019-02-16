package main

import (
	"fmt"
	"os"

	arnoldc "github.com/jmhobbs/go-arnoldc"
)

func main() {
	f, err := os.Open(os.Args[1])
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

	err = program.Run(os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
