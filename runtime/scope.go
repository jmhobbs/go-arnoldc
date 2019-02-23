package runtime

import (
	"fmt"

	arnoldc "github.com/jmhobbs/go-arnoldc"
)

type scope struct {
	parent    *scope
	variables map[string]int
}

func newScope(parent *scope) *scope {
	return &scope{parent, make(map[string]int)}
}

func (s *scope) Declare(name string, initial_value int) {
	s.variables[name] = initial_value
}

func (s *scope) Set(name string, value int) error {
	// TODO: This check should run after parsing to catch early
	if _, ok := s.variables[name]; !ok {
		if s.parent != nil {
			return s.parent.Set(name, value)
		}
		return fmt.Errorf("undefined variable %q", name)
	}
	s.variables[name] = value

	return nil
}

// Resolve a value to it's underlying integer, following variable references and up the scope chain.
func (s *scope) Get(v arnoldc.Value) (int, error) {
	switch v.Type() {
	case arnoldc.VariableType:
		var varName string = v.Value().(string)
		value, ok := s.variables[varName]
		if !ok {
			if s.parent != nil {
				return s.parent.Get(v)
			}
			return 0, fmt.Errorf("undefined variable %q", varName)
		}
		return value, nil
	case arnoldc.IntegerType:
		return v.Value().(int), nil
	default:
		return 0, fmt.Errorf("invalid value for number type")
	}
}

// Resolve a value to it's underlying type, following variable references and up the scope chain.
func (s *scope) resolveValue(v arnoldc.Value) (interface{}, error) {
	if v.Type() == arnoldc.VariableType {
		var varName string = v.Value().(string)
		value, ok := s.variables[varName]
		if !ok {
			if s.parent != nil {
				return s.parent.resolveValue(v)
			}
			return nil, fmt.Errorf("undefined variable %q", varName)
		}
		return value, nil
	}
	return v.Value(), nil
}
