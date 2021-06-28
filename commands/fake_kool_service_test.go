package commands

import (
	"errors"
	"kool-dev/kool/core/builder"
	"testing"
)

func TestFakeKoolService(t *testing.T) {
	f := &FakeKoolService{}

	_ = f.Execute([]string{"arg1", "arg2"})

	if !f.CalledExecute || len(f.ArgsExecute) != 2 || f.ArgsExecute[0] != "arg1" || f.ArgsExecute[1] != "arg2" {
		t.Errorf("failed to assert calling method Execute on FakeKoolService")
	}

	f.Println()

	if !f.CalledPrintln {
		t.Errorf("failed to assert calling method Println on FakeKoolService")
	}

	f.Printf("")

	if !f.CalledPrintf {
		t.Errorf("failed to assert calling method Printf on FakeKoolService")
	}

	f.Error(nil)

	if !f.CalledError {
		t.Errorf("failed to assert calling method Error on FakeKoolService")
	}

	f.Warning()

	if !f.CalledWarning {
		t.Errorf("failed to assert calling method Warning on FakeKoolService")
	}

	f.Success()

	if !f.CalledSuccess {
		t.Errorf("failed to assert calling method Success on FakeKoolService")
	}

	f.MockExecError = errors.New("error")

	err := f.Execute(nil)

	if err == nil || err.Error() != f.MockExecError.Error() {
		t.Errorf("failed to assert returning Execute mocked error on FakeKoolService")
	}

	f.IsTerminal()

	if !f.CalledIsTerminal {
		t.Errorf("failed to assert calling method IsTerminal on FakeKoolService")
	}

	f.InStream()

	if !f.CalledInStream {
		t.Errorf("failed to assert calling method InStream on FakeKoolService")
	}

	f.OutStream()

	if !f.CalledOutStream {
		t.Errorf("failed to assert calling method OutStream on FakeKoolService")
	}

	f.ErrStream()

	if !f.CalledErrStream {
		t.Errorf("failed to assert calling method ErrStream on FakeKoolService")
	}

	f.SetInStream(nil)

	if !f.CalledSetInStream {
		t.Errorf("failed to assert calling method SetInStream on FakeKoolService")
	}

	f.SetOutStream(nil)

	if !f.CalledSetOutStream {
		t.Errorf("failed to assert calling method SetOutStream on FakeKoolService")
	}

	f.SetErrStream(nil)

	if !f.CalledSetErrStream {
		t.Errorf("failed to assert calling method SetErrStream on FakeKoolService")
	}

	_, _ = f.Exec(&builder.FakeCommand{}, "extraArg")

	if !f.CalledExec {
		t.Errorf("failed to assert calling method Exec on FakeKoolService")
	}

	_ = f.Interactive(&builder.FakeCommand{}, "extraArg")

	if !f.CalledInteractive {
		t.Errorf("failed to assert calling method Interactive on FakeKoolService")
	}

	_ = f.LookPath(&builder.FakeCommand{})

	if !f.CalledLookPath {
		t.Errorf("failed to assert calling method LookPath on FakeKoolService")
	}
}
