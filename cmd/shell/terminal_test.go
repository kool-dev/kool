// +build !windows

package shell

import (
	"io"
	"os/exec"
	"testing"

	"github.com/creack/pty"
)

func TestPipeIsTerminal(t *testing.T) {
	c1 := exec.Command("echo", "testing")
	c2 := exec.Command("cat")

	r, w := io.Pipe()

	c1.Stdout = w
	c2.Stdin = r

	_ = c1.Start()
	_ = c2.Start()
	go func() {
		defer w.Close()

		_ = c1.Wait()
	}()
	_ = c2.Wait()

	terminalChecker := NewTerminalChecker()

	if terminalChecker.IsTerminal(c1.Stdout) {
		t.Error("unexpected tty terminal on piped command")
	}
}

func TestPtyIsTerminal(t *testing.T) {
	c := exec.Command("echo", "testing")
	f, err := pty.Start(c)

	if err != nil {
		t.Fatal(err)
	}

	terminalChecker := NewTerminalChecker()

	if !terminalChecker.IsTerminal(f) {
		t.Error("expecting tty terminal on command")
	}
}
