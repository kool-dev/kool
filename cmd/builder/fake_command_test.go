package builder

import "testing"

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

	f.Interactive("arg1", "arg2")

	if !f.CalledInteractive || f.ArgsInteractive == nil || f.ArgsInteractive[0] != "arg1" || f.ArgsInteractive[1] != "arg2" {
		t.Errorf("failed to use mocked Interactive function on FakeCommand")
	}

	f.Exec("arg1", "arg2")

	if !f.CalledExec || f.ArgsExec == nil || f.ArgsExec[0] != "arg1" || f.ArgsExec[1] != "arg2" {
		t.Errorf("failed to use mocked Exec function on FakeCommand")
	}
}
