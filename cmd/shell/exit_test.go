package shell

import "testing"

func TestNewExiter(t *testing.T) {
	exiter := NewExiter()

	if _, assert := exiter.(*DefaultExiter); !assert {
		t.Errorf("NewExiter() did not return a *DefaultExiter")
	}
}

func TestExitCodeZero(t *testing.T) {
	if code := execExitCode(0); code != 0 {
		t.Errorf("Exit() failed; expected exit code '0', got %v", code)
	}
}

func TestExitCodeOne(t *testing.T) {
	if code := execExitCode(1); code != 1 {
		t.Errorf("Exit() failed; expected exit code '1', got %v", code)
	}
}

func TestExitCodeTwo(t *testing.T) {
	if code := execExitCode(2); code != 2 {
		t.Errorf("Exit() failed; expected exit code '2', got %v", code)
	}
}

func execExitCode(originalCode int) int {
	var (
		mockExitCode  = -1
		defaultExiter = &DefaultExiter{func(code int) {
			mockExitCode = code
		}}
	)

	defaultExiter.Exit(originalCode)

	return mockExitCode
}
