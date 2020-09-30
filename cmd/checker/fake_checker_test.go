package checker

import (
	"errors"
	"testing"
)

func TestFakeChecker(t *testing.T) {
	f := &FakeChecker{}

	_ = f.Check()

	if !f.CalledCheck {
		t.Error("failed to use mocked Check function on FakeChecker")
	}
}

func TestFailedFakeChecker(t *testing.T) {
	f := &FakeChecker{MockError: errors.New("fake error")}

	err := f.Check()

	if err == nil {
		t.Error("failed to use mocked failed Check function on FakeChecker")
	} else if err.Error() != "fake error" {
		t.Error("failed to use mocked failed Check function on FakeChecker")
	}
}
