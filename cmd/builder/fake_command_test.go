package builder

import (
	"errors"
	"testing"
)

func TestFakeCommand(t *testing.T) {
	f := &FakeCommand{}

	f.AppendArgs("arg1", "arg2")

	if !f.CalledAppendArgs || f.ArgsAppend == nil || f.ArgsAppend[0] != "arg1" || f.ArgsAppend[1] != "arg2" {
		t.Errorf("failed to use mocked AppendArgs function on FakeCommand")
	}

	_ = f.String()

	if !f.CalledString {
		t.Errorf("failed to use mocked String function on FakeCommand")
	}

	_ = f.LookPath()

	if !f.CalledLookPath {
		t.Errorf("failed to use mocked LookPath function on FakeCommand")
	}

	_ = f.Interactive("arg1", "arg2")

	if !f.CalledInteractive || f.ArgsInteractive == nil || f.ArgsInteractive[0] != "arg1" || f.ArgsInteractive[1] != "arg2" {
		t.Errorf("failed to use mocked Interactive function on FakeCommand")
	}

	_, _ = f.Exec("arg1", "arg2")

	if !f.CalledExec || f.ArgsExec == nil || f.ArgsExec[0] != "arg1" || f.ArgsExec[1] != "arg2" {
		t.Errorf("failed to use mocked Exec function on FakeCommand")
	}
}

func TestFakeFailedCommand(t *testing.T) {
	mockErr := errors.New("error")
	f := &FakeCommand{MockError: mockErr, MockLookPathError: mockErr}

	if err := f.LookPath(); err == nil {
		t.Errorf("failed to mock error calling LookPath function on FakeCommand")
	}

	if !f.CalledLookPath {
		t.Errorf("failed to use mocked LookPath function on FakeCommand")
	}

	if err := f.Interactive("arg1", "arg2"); err == nil {
		t.Errorf("failed to mock error calling Interactive function on FakeCommand")
	}

	if !f.CalledInteractive || f.ArgsInteractive == nil || f.ArgsInteractive[0] != "arg1" || f.ArgsInteractive[1] != "arg2" {
		t.Errorf("failed to use mocked Interactive function on FakeCommand")
	}

	if _, err := f.Exec("arg1", "arg2"); err == nil {
		t.Errorf("failed to mock error calling Exec function on FakeCommand")
	}

	if !f.CalledExec || f.ArgsExec == nil || f.ArgsExec[0] != "arg1" || f.ArgsExec[1] != "arg2" {
		t.Errorf("failed to use mocked Exec function on FakeCommand")
	}
}
