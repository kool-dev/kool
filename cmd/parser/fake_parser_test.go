package parser

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"testing"
)

func TestFakeParser(t *testing.T) {
	f := &FakeParser{
		MockParsedCommands: []builder.Command{&builder.FakeCommand{}},
	}

	_ = f.AddLookupPath("path")

	if !f.CalledAddLookupPath || len(f.TargetFiles) != 1 {
		t.Error("failed to use mocked AddLookupPath function on FakeParser")
	}

	_ = f.AddLookupPath("path")

	if len(f.TargetFiles) != 2 {
		t.Error("failed to use mocked AddLookupPath function more then once on FakeParser")
	}

	commands, _ := f.Parse("script")

	if !f.CalledParse || len(commands) != 1 {
		t.Error("failed to use mocked Parse function more then once on FakeParser")
	}
}

func TestFakeFailedParser(t *testing.T) {
	f := &FakeParser{
		MockParseError: errors.New("error"),
	}

	_, err := f.Parse("script")

	if !f.CalledParse || err == nil {
		t.Error("failed to use mocked failing Parse function more then once on FakeParser")
	}
}
