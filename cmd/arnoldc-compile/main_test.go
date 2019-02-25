package main

import (
	"testing"

	arnoldc "github.com/jmhobbs/go-arnoldc"
)

func TestGoStringForValue(t *testing.T) {

	cases := []struct {
		value  arnoldc.Value
		expect string
	}{
		{
			arnoldc.NewVariableValue("aVariable"),
			"aVariable",
		},
		{
			arnoldc.NewIntegerValue(5),
			"5",
		},
		{
			arnoldc.NewStringValue("a string here"),
			`"a string here"`,
		},
	}

	for _, test := range cases {
		v := goStringForValue(test.value)
		if v != test.expect {
			t.Errorf("bad string for value; got %v expected %v", v, test.expect)
		}
	}
}
