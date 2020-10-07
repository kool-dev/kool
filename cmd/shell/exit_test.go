package shell

import "testing"

func TestNewExiter(t *testing.T) {
	e := NewExiter()

	if _, ok := e.(*DefaultExiter); !ok {
		t.Errorf("NewExiter() did not return a *DefaultExiter")
	}
}

func TestExitCode1NewExiter(t *testing.T) {
	var exitCode int

	oldExit := exitFn

	exitFn = func(code int) {
		exitCode = code
	}

	defer func() { exitFn = oldExit }()

	e := NewExiter()
	e.Exit(1)

	if exitCode != 1 {
		t.Errorf("expecting exit code 1, got '%v'", exitCode)
	}
}

func TestExitCode2NewExiter(t *testing.T) {
	var exitCode int

	oldExit := exitFn

	exitFn = func(code int) {
		exitCode = code
	}

	defer func() { exitFn = oldExit }()

	e := NewExiter()
	e.Exit(2)

	if exitCode != 2 {
		t.Errorf("expecting exit code 2, got '%v'", exitCode)
	}
}
