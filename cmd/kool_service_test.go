package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/shell"
	"strings"
	"testing"
)

func newFakeKoolService() *DefaultKoolService {
	return &DefaultKoolService{
		&shell.FakeExiter{},
		&shell.FakeOutputWriter{},
		&shell.FakeInputReader{},
		&shell.FakeTerminalChecker{MockIsTerminal: true},
		&shell.FakeShell{},
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

	out = []interface{}{"success"}
	k.Println(out...)

	if !k.out.(*shell.FakeOutputWriter).CalledPrintln {
		t.Error("Println was not proxied by DefaultKoolService")
	}

	expected := strings.TrimSpace(fmt.Sprintln(out...))
	if len(k.out.(*shell.FakeOutputWriter).OutLines[0]) != len(expected) {
		t.Errorf("Println did not proxy the proper output on DefaultKoolService; expected %v got %v", expected, k.out.(*shell.FakeOutputWriter).OutLines[0])
	}

	k.Printf("testing %s", "format")

	if !k.out.(*shell.FakeOutputWriter).CalledPrintf {
		t.Error("Printf was not proxied by DefaultKoolService")
	}

	expectedFOutput := "testing format"
	if fOutput := k.out.(*shell.FakeOutputWriter).FOutput; fOutput != expectedFOutput {
		t.Errorf("Printf did not proxy the proper output on DefaultKoolService; expected '%s', got %s", expectedFOutput, fOutput)
	}

	k.SetWriter(nil)

	if !k.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Error("SetWriter was not proxied by DefaultKoolService")
	}

	k.GetWriter()

	if !k.out.(*shell.FakeOutputWriter).CalledGetWriter {
		t.Error("GetWriter was not proxied by DefaultKoolService")
	}

	k.SetReader(nil)

	if !k.in.(*shell.FakeInputReader).CalledSetReader {
		t.Error("SetReader was not proxied by DefaultKoolService")
	}

	k.GetReader()

	if !k.in.(*shell.FakeInputReader).CalledGetReader {
		t.Error("GetReader was not proxied by DefaultKoolService")
	}
}
