package runtime

import (
	"testing"

	arnoldc "github.com/jmhobbs/go-arnoldc"
)

func TestScopeResolution(t *testing.T) {
	root := newScope(nil)
	child := newScope(root)
	secondChild := newScope(child)

	root.Declare("a", 5)
	root.Declare("b", 5)
	child.Declare("b", 7)
	child.Declare("c", 10)
	secondChild.Declare("d", 10)

	t.Run("root exists", func(t *testing.T) {
		v, err := root.Get(arnoldc.NewVariableValue("a"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if v != 5 {
			t.Errorf("incorrect value returned; expected 5, got %v", v)
		}
	})

	t.Run("root does not exist", func(t *testing.T) {
		_, err := root.Get(arnoldc.NewVariableValue("c"))
		if err == nil {
			t.Errorf("expected error, did not get one")
		}
	})

	t.Run("child exists", func(t *testing.T) {
		v, err := child.Get(arnoldc.NewVariableValue("c"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if v != 10 {
			t.Errorf("incorrect value returned; expected 10, got %v", v)
		}
	})

	t.Run("child shadows parent", func(t *testing.T) {
		v, err := child.Get(arnoldc.NewVariableValue("b"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if v != 7 {
			t.Errorf("incorrect value returned; expected 7, got %v", v)
		}
	})

	t.Run("child falls through to parent", func(t *testing.T) {
		v, err := child.Get(arnoldc.NewVariableValue("a"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if v != 5 {
			t.Errorf("incorrect value returned; expected 5, got %v", v)
		}
	})

	t.Run("second child falls through to parent", func(t *testing.T) {
		v, err := secondChild.Get(arnoldc.NewVariableValue("b"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if v != 7 {
			t.Errorf("incorrect value returned; expected 7, got %v", v)
		}
	})

	t.Run("second child falls through to root", func(t *testing.T) {
		v, err := secondChild.Get(arnoldc.NewVariableValue("a"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if v != 5 {
			t.Errorf("incorrect value returned; expected 5, got %v", v)
		}
	})
}
