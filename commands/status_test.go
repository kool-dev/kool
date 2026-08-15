package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/network"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/checker"
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
		*(newDefaultKoolService().Fake()),
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

	if _, ok := k.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolStatus instance")
	}

	if _, ok := k.check.(*checker.DefaultChecker); !ok {
		t.Errorf("unexpected checker.Checker on default KoolStatus instance")
	}

	if _, ok := k.net.(*network.DefaultHandler); !ok {
		t.Errorf("unexpected network.Handler on default KoolStatus instance")
	}

	if _, ok := k.getServicesCmd.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Command on default KoolStatus instance")
	}

	if _, ok := k.getServiceIDCmd.(*builder.DefaultCommand); !ok {
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

	expected := ""

	output := fmt.Sprint(f.shell.(*shell.FakeShell).WarningOutput...)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}

	f.getServicesCmd.(*builder.FakeCommand).MockExecError = nil

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected = "No services found."

	output = fmt.Sprint(f.shell.(*shell.FakeShell).WarningOutput...)

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
		*(newDefaultKoolService().Fake()),
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

func newFakeKoolStatusJSON() *KoolStatus {
	f := &KoolStatus{
		*(newDefaultKoolService().Fake()),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&shell.FakeTableWriter{},
	}

	f.shell.(*shell.FakeShell).MockIsJSONOutput = true
	f.shell.(*shell.FakeShell).MockErrStream = io.Discard
	f.shell.(*shell.FakeShell).MockOutStream = io.Discard

	return f
}

func TestStatusJSONOutput(t *testing.T) {
	f := newFakeKoolStatusJSON()

	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"
	f.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Up About an hour|0.0.0.0:80->80/tcp, 9000/tcp"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	fakeShell := f.shell.(*shell.FakeShell)
	if len(fakeShell.OutLines) == 0 {
		t.Fatal("expected JSON output")
	}

	var output statusOutputJSON
	if err := json.Unmarshal([]byte(fakeShell.OutLines[0]), &output); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nraw: %s", err, fakeShell.OutLines[0])
	}

	if output.Count != 1 {
		t.Errorf("expected count 1, got %d", output.Count)
	}

	if len(output.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(output.Services))
	}

	svc := output.Services[0]
	if svc.Service != "app" {
		t.Errorf("expected service 'app', got '%s'", svc.Service)
	}
	if !svc.Running {
		t.Error("expected running to be true")
	}
	if svc.State != "Up About an hour" {
		t.Errorf("expected state 'Up About an hour', got '%s'", svc.State)
	}
	if svc.Ports != "0.0.0.0:80->80/tcp, 9000/tcp" {
		t.Errorf("expected ports '0.0.0.0:80->80/tcp, 9000/tcp', got '%s'", svc.Ports)
	}
}

func TestStatusJSONOutputNotRunning(t *testing.T) {
	f := newFakeKoolStatusJSON()

	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"
	f.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Exited an hour ago"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	fakeShell := f.shell.(*shell.FakeShell)
	var output statusOutputJSON
	if err := json.Unmarshal([]byte(fakeShell.OutLines[0]), &output); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if output.Services[0].Running {
		t.Error("expected running to be false for exited container")
	}
}

func TestStatusJSONOutputMultipleServices(t *testing.T) {
	f := newFakeKoolStatusJSON()

	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "cache\napp"
	f.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"
	f.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Up|0.0.0.0:80->80/tcp"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	fakeShell := f.shell.(*shell.FakeShell)
	var output statusOutputJSON
	if err := json.Unmarshal([]byte(fakeShell.OutLines[0]), &output); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if output.Count != 2 {
		t.Errorf("expected count 2, got %d", output.Count)
	}
}
