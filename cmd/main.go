package main

import (
	"flag"
	"fmt"
	"os"

	arnoldc "github.com/jmhobbs/go-arnoldc"
)

func main() {
	var debug = flag.Bool("debug", false, "Run with verbose debugging.")
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening program: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	a := arnoldc.New(f)
	a.Debug = *debug
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
