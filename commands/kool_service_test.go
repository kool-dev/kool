package commands

import (
	"kool-dev/kool/core/shell"
	"testing"
)

func TestKoolServiceProxies(t *testing.T) {
	k := &DefaultKoolService{
		&shell.FakeShell{},
	}

	if _, ok := k.Shell().(*shell.FakeShell); !ok {
		t.Error("unexpected Shell return")
	}
}

// func TestKoolServiceErrors(t *testing.T) {
// 	k := newFakeKoolService()

// 	err := k.Interactive(&builder.FakeCommand{MockInteractiveError: shell.ErrLookPath}, "extraArg")

// 	if err == nil || !strings.Contains(err.Error(), "failed to run") {
// 		t.Errorf("bad error returned: %v", err)
// 	}

// 	exitError := &exec.ExitError{ProcessState: &os.ProcessState{}}
// 	exitStatus := exitError.Sys().(syscall.WaitStatus).ExitStatus()

// 	command := &builder.FakeCommand{MockInteractiveError: exitError}

// 	err = k.Interactive(command)

// 	if ex, ok := err.(shell.ErrExitable); !ok {
// 		t.Error("should be ErrExitable")
// 	} else if ex.Code != exitStatus {
// 		t.Errorf("exit code should be %d but got %d", exitStatus, ex.Code)
// 	}
// }
