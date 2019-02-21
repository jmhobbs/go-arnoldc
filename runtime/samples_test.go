package runtime

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	arnoldc "github.com/jmhobbs/go-arnoldc"
)

func TestSamples(t *testing.T) {
	samples := []struct {
		file     string
		expected string
	}{
		{
			"hello_world.arnoldc",
			"hello world\n",
		},
		{
			"hello_variable.arnoldc",
			"10\n",
		},
		{
			"hello_math.arnoldc",
			"28\n",
		},
		{
			"hello_logic.arnoldc",
			"1\n0\n",
		},
		{
			"hello_conditionals.arnoldc",
			"true is true\nfalse is not true\na is greater than b\n",
		},
		{
			"count_to_ten.arnoldc",
			"1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n",
		},
	}

	for _, sample := range samples {
		t.Run(sample.file, func(t *testing.T) {
			var stdout bytes.Buffer

			f, err := os.Open(filepath.Join(".", "samples", sample.file))
			if err != nil {
				t.Errorf("error opening sample: %v", err)
				return
			}
			defer f.Close()

			a := arnoldc.New(f)
			program, err := a.Parse()
			if err != nil {
				t.Errorf("error parsing program: %v", err)
				return
			}

			err = New(&stdout, &stdout).Run(program)
			if err != nil {
				t.Errorf("error running program: %v", err)
				return
			}

			output := stdout.String()
			if sample.expected != output {
				t.Errorf("program output does not match expectations\nexpect: %q\noutput: %q", sample.expected, output)
			}
		})
	}
}
