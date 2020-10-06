package shell

import (
	"os"
	"os/exec"
	"testing"
)

func TestNewExiter(t *testing.T) {
	e := NewExiter()

	if _, ok := e.(*DefaultExiter); !ok {
		t.Errorf("NewExiter() did not return a *DefaultExiter")
	}
}

func TestExitCode1NewExiter(t *testing.T) {
	if os.Getenv("TESTING_FLAG") == "1" {
		e := NewExiter()
		e.Exit(1)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExitCode1NewExiter")
	cmd.Env = append(os.Environ(), "TESTING_FLAG=1")

	err := cmd.Run()

	e, ok := err.(*exec.ExitError)

	if !ok {
		t.Errorf("exiter did not exit")
	}

	if e.Error() != "exit status 1" {
		t.Errorf("expecting exit status 1, got '%s'", e.Error())
	}
}

func TestExitCode2NewExiter(t *testing.T) {
	if os.Getenv("TESTING_FLAG") == "1" {
		e := NewExiter()
		e.Exit(2)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExitCode2NewExiter")
	cmd.Env = append(os.Environ(), "TESTING_FLAG=1")

	err := cmd.Run()

	e, ok := err.(*exec.ExitError)

	if !ok {
		t.Errorf("exiter did not exit")
	}

	if e.Error() != "exit status 2" {
		t.Errorf("expecting exit status 2, got '%s'", e.Error())
	}
}
