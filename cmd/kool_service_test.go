package cmd

import (
	"errors"
	"kool-dev/kool/cmd/shell"
	"testing"
)

func newFakeKoolService() *DefaultKoolService {
	return &DefaultKoolService{
		&shell.FakeExiter{},
		&shell.FakeOutputWriter{},
	}
}

func TestKoolServiceProxies(t *testing.T) {
	code := 100
	k := newFakeKoolService()

	k.Exit(code)

	if !k.exiter.(*shell.FakeExiter).Exited() {
		t.Error("Exit was not proxied by DefaultKoolService")
	}

	if k.exiter.(*shell.FakeExiter).Code() != code {
		t.Errorf("Exit did not proxy the proper code by DefaultKoolService; expected %d got %d", code, k.exiter.(*shell.FakeExiter).Code())
	}

	err := errors.New("fake error")
	k.Error(err)

	if !k.out.(*shell.FakeOutputWriter).CalledError {
		t.Error("Error was not proxied by DefaultKoolService")
	}

	if k.out.(*shell.FakeOutputWriter).Err != err {
		t.Errorf("Error did not proxy the proper error on DefaultKoolService; expected %v got %v", err, k.out.(*shell.FakeOutputWriter).Err)
	}

	out := []interface{}{"out"}
	k.Warning(out...)

	if !k.out.(*shell.FakeOutputWriter).CalledWarning {
		t.Error("Warning was not proxied by DefaultKoolService")
	}

	if len(k.out.(*shell.FakeOutputWriter).WarningOutput) != len(out) {
		t.Errorf("Warning did not proxy the proper output on DefaultKoolService; expected %v got %v", out, k.out.(*shell.FakeOutputWriter).WarningOutput)
	}

	out = []interface{}{"success"}
	k.Success(out...)

	if !k.out.(*shell.FakeOutputWriter).CalledSuccess {
		t.Error("Success was not proxied by DefaultKoolService")
	}

	if len(k.out.(*shell.FakeOutputWriter).SuccessOutput) != len(out) {
		t.Errorf("Success did not proxy the proper output on DefaultKoolService; expected %v got %v", out, k.out.(*shell.FakeOutputWriter).SuccessOutput)
	}

	k.SetWriter(nil)

	if !k.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Error("SetWriter was not proxied by DefaultKoolService")
	}
}
