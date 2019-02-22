package arnoldc

import "testing"

func TestInstructionToString(t *testing.T) {
	mo := instructionToString(MAIN_OPEN)
	if mo != "IT'S SHOWTIME" {
		t.Errorf("expected %q, got %q", "IT'S SHOWTIME", mo)
	}
}
