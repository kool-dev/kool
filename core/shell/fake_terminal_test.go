package shell

import (
	"os"
	"testing"
)

func TestFakeTerminalChecker(t *testing.T) {
	f := &FakeTerminalChecker{}
	f.MockIsTerminal = true

	isTerminal := f.IsTerminal(os.Stdin, os.Stdout)

	if !f.CalledIsTerminal || !isTerminal {
		t.Error("failed to use mocked IsTerminal function on FakeTerminalChecker")
	}
}
