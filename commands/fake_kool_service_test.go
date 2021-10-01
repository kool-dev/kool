package commands

import (
	"testing"
)

func TestFakeKoolService(t *testing.T) {
	f := &FakeKoolService{}

	_ = f.Execute([]string{"arg1", "arg2"})

	if !f.CalledExecute || len(f.ArgsExecute) != 2 || f.ArgsExecute[0] != "arg1" || f.ArgsExecute[1] != "arg2" {
		t.Errorf("failed to assert calling method Execute on FakeKoolService")
	}

	if f.Shell() != nil {
		t.Errorf("unexpected non-nil default shell; got: %v - %v", f.Shell(), f.Shell() == nil)
	}
}
