package cmd

import (
	"kool-dev/kool/cmd/shell"
	"testing"
)

func newFakeKoolService() *DefaultKoolService {
	return &DefaultKoolService{
		&shell.FakeExiter{},
		&shell.FakeOutputWriter{},
	}
}

func TestExitProxy(t *testing.T) {
	code := 100
	k := newFakeKoolService()

	k.Exit(code)

	if !k.exiter.(*shell.FakeExiter).Exited() {
		t.Error("Exit was not proxied by DefaultKoolService")
	}

	if k.exiter.(*shell.FakeExiter).Code() != code {
		t.Errorf("Exit did not proxy the proper code by DefaultKoolService; expected %d got %d", code, k.exiter.(*shell.FakeExiter).Code())
	}
}
