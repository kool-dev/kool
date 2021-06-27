package commands

import (
	"errors"
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/network"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/checker"
	"kool-dev/kool/services/compose"
	"strings"
	"testing"
)

type FakeRaceShell struct {
	shell.FakeShell
}

func (f *FakeRaceShell) Exec(command builder.Command, extraArgs ...string) (string, error) {
	output := command.(*builder.FakeCommand).MockExecOut
	return output, nil
}

func newFakeKoolStatus() *KoolStatus {
	fs := &KoolStatus{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&shell.FakeTableWriter{},
	}

	fs.shell.(*shell.FakeShell).MockErrStream = io.Discard
	fs.shell.(*shell.FakeShell).MockOutStream = io.Discard

	return fs
}

func TestNewKoolStatus(t *testing.T) {
	k := NewKoolStatus()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolStatus instance")
	}

	if _, ok := k.DefaultKoolService.term.(*shell.DefaultTerminalChecker); !ok {
		t.Errorf("unexpected shell.TerminalChecker on default KoolStatus instance")
	}

	if _, ok := k.check.(*checker.DefaultChecker); !ok {
		t.Errorf("unexpected checker.Checker on default KoolStatus instance")
	}

	if _, ok := k.net.(*network.DefaultHandler); !ok {
		t.Errorf("unexpected network.Handler on default KoolStatus instance")
	}

	if _, ok := k.getServicesCmd.(*compose.DockerCompose); !ok {
		t.Errorf("unexpected builder.Command on default KoolStatus instance")
	}

	if _, ok := k.getServiceIDCmd.(*compose.DockerCompose); !ok {
		t.Errorf("unexpected builder.Command on default KoolStatus instance")
	}

	if _, ok := k.getServiceStatusPortCmd.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Command on default KoolStatus instance")
	}

	if _, ok := k.table.(*shell.DefaultTableWriter); !ok {
		t.Errorf("unexpected shell.TableWriter on default KoolStatus instance")
	}
}

func TestStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"
	f.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Up About an hour|0.0.0.0:80->80/tcp, 9000/tcp"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := `Service | Running | Ports | State
app | Running | 0.0.0.0:80->80/tcp, 9000/tcp | Up About an hour`

	output := strings.TrimSpace(f.table.(*shell.FakeTableWriter).TableOut)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNotRunningStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"
	f.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Exited an hour ago"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := `Service | Running | Ports | State
app | Not running |  | Exited an hour ago`

	output := strings.TrimSpace(f.table.(*shell.FakeTableWriter).TableOut)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNoStatusPortStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := `Service | Running | Ports | State
app | Not running |  |`

	output := strings.TrimSpace(f.table.(*shell.FakeTableWriter).TableOut)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNoServicesStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()
	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := "No services found."

	output := fmt.Sprint(f.shell.(*shell.FakeShell).WarningOutput...)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFailedGetServicesStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesCmd.(*builder.FakeCommand).MockExecError = errors.New("exec err")

	cmd := NewStatusCommand(f)

	assertExecGotError(t, cmd, "exec err")

	expected := "No services found."

	output := fmt.Sprint(f.shell.(*shell.FakeShell).WarningOutput...)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFailedDependenciesStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()
	f.check.(*checker.FakeChecker).MockError = errors.New("exec error")

	cmd := NewStatusCommand(f)

	assertExecGotError(t, cmd, "exec error")
}

func TestFailedNetworkStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()
	f.net.(*network.FakeHandler).MockError = errors.New("exec network error")

	cmd := NewStatusCommand(f)

	assertExecGotError(t, cmd, "exec network error")
}

func TestFailedGetServiceIDStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDCmd.(*builder.FakeCommand).MockExecError = errors.New("get service error")

	cmd := NewStatusCommand(f)

	assertExecGotError(t, cmd, "get service error")
}

func TestServicesOrderStatusCommand(t *testing.T) {
	f := &KoolStatus{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&shell.FakeTableWriter{},
	}

	f.shell = &FakeRaceShell{
		FakeShell: shell.FakeShell{
			MockErrStream: io.Discard,
			MockOutStream: io.Discard,
		},
	}
	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = `cache
app`
	f.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "output"
	f.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "output"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := `Service | Running | Ports | State
app | Not running |  | output
cache | Not running |  | output`

	output := strings.TrimSpace(f.table.(*shell.FakeTableWriter).TableOut)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}
