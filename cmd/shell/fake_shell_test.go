package shell

import (
	"errors"
	"io/ioutil"
	"kool-dev/kool/cmd/builder"
	"os"
	"testing"
)

func TestFakeShell(t *testing.T) {
	f := &FakeShell{}

	f.MockInStream = os.Stdin

	in := f.InStream()

	if !f.CalledInStream || in != os.Stdin {
		t.Error("failed to use mocked InStream function on FakeShell")
	}

	f.SetInStream(nil)

	if !f.CalledSetInStream {
		t.Error("failed to use mocked SetInStream function on FakeShell")
	}

	f.MockOutStream = ioutil.Discard

	out := f.OutStream()

	if !f.CalledOutStream || out != ioutil.Discard {
		t.Error("failed to use mocked OutStream function on FakeShell")
	}

	f.SetOutStream(nil)

	if !f.CalledSetOutStream {
		t.Error("failed to use mocked SetOutStream function on FakeShell")
	}

	f.MockErrStream = ioutil.Discard

	err := f.ErrStream()

	if !f.CalledErrStream || err != ioutil.Discard {
		t.Error("failed to use mocked ErrStream function on FakeShell")
	}

	f.SetErrStream(nil)

	if !f.CalledSetErrStream {
		t.Error("failed to use mocked SetErrStream function on FakeShell")
	}

	command := &builder.FakeCommand{
		MockCmd:              "cmd",
		MockExecError:        errors.New("error exec"),
		MockExecOut:          "exec output",
		MockInteractiveError: errors.New("error interactive"),
		MockLookPathError:    errors.New("error lookPath"),
	}

	execOut, execError := f.Exec(command)

	if val, ok := f.CalledExec["cmd"]; !val || !ok || execOut != command.MockExecOut || execError != command.MockExecError {
		t.Error("failed to use mocked Exec function on FakeShell")
	}

	interactiveError := f.Interactive(command)

	if val, ok := f.CalledInteractive["cmd"]; !val || !ok || interactiveError != command.MockInteractiveError {
		t.Error("failed to use mocked Interactive function on FakeShell")
	}

	lookPathError := f.LookPath(command)

	if val, ok := f.CalledLookPath["cmd"]; !val || !ok || lookPathError != command.MockLookPathError {
		t.Error("failed to use mocked LookPath function on FakeShell")
	}

	f.Println()

	if !f.CalledPrintln {
		t.Errorf("failed to assert calling method Println on FakeShell")
	}

	f.Printf("")

	if !f.CalledPrintf {
		t.Errorf("failed to assert calling method Printf on FakeShell")
	}

	f.Error(nil)

	if !f.CalledError {
		t.Errorf("failed to assert calling method Error on FakeShell")
	}

	f.Warning()

	if !f.CalledWarning {
		t.Errorf("failed to assert calling method Warning on FakeShell")
	}

	f.Success()

	if !f.CalledSuccess {
		t.Errorf("failed to assert calling method Success on FakeShell")
	}
}
