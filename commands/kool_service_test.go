package commands

import (
	"errors"
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"strings"
	"testing"
)

func newFakeKoolService() *DefaultKoolService {
	return &DefaultKoolService{
		&shell.FakeTerminalChecker{MockIsTerminal: true},
		&shell.FakeShell{},
	}
}

func TestKoolServiceProxies(t *testing.T) {
	k := newFakeKoolService()

	err := errors.New("fake error")
	k.Error(err)

	if !k.shell.(*shell.FakeShell).CalledError {
		t.Error("Error was not proxied by DefaultKoolService")
	}

	if k.shell.(*shell.FakeShell).Err != err {
		t.Errorf("Error did not proxy the proper error on DefaultKoolService; expected %v got %v", err, k.shell.(*shell.FakeShell).Err)
	}

	out := []interface{}{"out"}
	k.Warning(out...)

	if !k.shell.(*shell.FakeShell).CalledWarning {
		t.Error("Warning was not proxied by DefaultKoolService")
	}

	if len(k.shell.(*shell.FakeShell).WarningOutput) != len(out) {
		t.Errorf("Warning did not proxy the proper output on DefaultKoolService; expected %v got %v", out, k.shell.(*shell.FakeShell).WarningOutput)
	}

	out = []interface{}{"success"}
	k.Success(out...)

	if !k.shell.(*shell.FakeShell).CalledSuccess {
		t.Error("Success was not proxied by DefaultKoolService")
	}

	if len(k.shell.(*shell.FakeShell).SuccessOutput) != len(out) {
		t.Errorf("Success did not proxy the proper output on DefaultKoolService; expected %v got %v", out, k.shell.(*shell.FakeShell).SuccessOutput)
	}

	out = []interface{}{"success"}
	k.Println(out...)

	if !k.shell.(*shell.FakeShell).CalledPrintln {
		t.Error("Println was not proxied by DefaultKoolService")
	}

	expected := strings.TrimSpace(fmt.Sprintln(out...))
	if len(k.shell.(*shell.FakeShell).OutLines[0]) != len(expected) {
		t.Errorf("Println did not proxy the proper output on DefaultKoolService; expected %v got %v", expected, k.shell.(*shell.FakeShell).OutLines[0])
	}

	k.Printf("testing %s", "format")

	if !k.shell.(*shell.FakeShell).CalledPrintf {
		t.Error("Printf was not proxied by DefaultKoolService")
	}

	expectedFOutput := "testing format"
	if fOutput := k.shell.(*shell.FakeShell).FOutput; fOutput != expectedFOutput {
		t.Errorf("Printf did not proxy the proper output on DefaultKoolService; expected '%s', got %s", expectedFOutput, fOutput)
	}

	k.InStream()

	if !k.shell.(*shell.FakeShell).CalledInStream {
		t.Errorf("failed to assert calling method InStream on FakeKoolService")
	}

	k.OutStream()

	if !k.shell.(*shell.FakeShell).CalledOutStream {
		t.Errorf("failed to assert calling method OutStream on FakeKoolService")
	}

	k.ErrStream()

	if !k.shell.(*shell.FakeShell).CalledErrStream {
		t.Errorf("failed to assert calling method ErrStream on FakeKoolService")
	}

	k.SetInStream(nil)

	if !k.shell.(*shell.FakeShell).CalledSetInStream {
		t.Errorf("failed to assert calling method SetInStream on FakeKoolService")
	}

	k.SetOutStream(nil)

	if !k.shell.(*shell.FakeShell).CalledSetOutStream {
		t.Errorf("failed to assert calling method SetOutStream on FakeKoolService")
	}

	k.SetErrStream(nil)

	if !k.shell.(*shell.FakeShell).CalledSetErrStream {
		t.Errorf("failed to assert calling method SetErrStream on FakeKoolService")
	}

	_, _ = k.Exec(&builder.FakeCommand{MockCmd: "cmd"}, "extraArg")

	if val, ok := k.shell.(*shell.FakeShell).CalledExec["cmd"]; !val || !ok {
		t.Errorf("failed to assert calling method Exec on FakeKoolService")
	}

	_ = k.Interactive(&builder.FakeCommand{MockCmd: "cmd"}, "extraArg")

	if val, ok := k.shell.(*shell.FakeShell).CalledInteractive["cmd"]; !val || !ok {
		t.Errorf("failed to assert calling method Interactive on FakeKoolService")
	}

	_ = k.LookPath(&builder.FakeCommand{MockCmd: "cmd"})

	if val, ok := k.shell.(*shell.FakeShell).CalledLookPath["cmd"]; !val || !ok {
		t.Errorf("failed to assert calling method LookPath on FakeKoolService")
	}
}
