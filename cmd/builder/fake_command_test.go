package builder

import (
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

	f.MockCmd = "cmd"

	if cmd := f.Cmd(); !f.CalledCmd || cmd != "cmd" {
		t.Errorf("failed to use mocked Cmd function on FakeCommand")
	}

	if args := f.Args(); !f.CalledArgs || len(args) != 2 || args[0] != "arg1" || args[1] != "arg2"{
		t.Errorf("failed to use mocked Args function on FakeCommand")
	}
}
