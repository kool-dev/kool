package shell

import "testing"

func TestFakeExiter(t *testing.T) {
	f := &FakeExiter{}

	f.Exit(1)

	code := f.Code()
	if !f.Exited() || code != 1 {
		t.Error("failed to use mocked Exit function on FakeExiter")
	}
}

func TestCode2FakeExiter(t *testing.T) {
	f := &FakeExiter{}

	f.Exit(2)

	code := f.Code()
	if !f.Exited() || code != 2 {
		t.Error("failed to use mocked Exit function on FakeExiter")
	}
}
