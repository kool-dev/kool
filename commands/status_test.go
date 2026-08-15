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
	"strings"
	"testing"
)

type FakeRaceShell struct {
	shell.FakeShell
}

type statusRecordingShell struct {
	shell.FakeShell
	args [][]string
}

func (f *statusRecordingShell) Exec(command builder.Command, extraArgs ...string) (string, error) {
	if len(extraArgs) > 0 {
		f.args = append(f.args, append([]string(nil), extraArgs...))
	}
	if fake, ok := command.(*builder.FakeCommand); ok {
		return fake.MockExecOut, fake.MockExecError
	}
	return "", nil
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
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&shell.FakeTableWriter{},
	}

	fs.shell.(*shell.FakeShell).MockErrStream = io.Discard
	fs.shell.(*shell.FakeShell).MockOutStream = io.Discard
	fs.env.Set("KOOL_NAME", "example")

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

func TestStatusFiltersWorkspaceServices(t *testing.T) {
	f := newFakeKoolStatus()
	f.shell = &FakeRaceShell{FakeShell: shell.FakeShell{MockErrStream: io.Discard, MockOutStream: io.Discard}}
	f.env.Set("KOOL_WORKSPACE", "true")
	f.env.Set("KOOL_WORKSPACES_ENABLED", "true")
	f.env.Set("KOOL_WORKSPACE_SERVICES", "app,node")
	f.env.Set("KOOL_WORKSPACE_SOURCE_PROJECT", "example")
	f.env.Set("KOOL_WORKSPACE_PROJECT", "example-workspace-task-a")
	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "app\ndatabase\nnode"
	f.getProjectServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"

	cmd := NewStatusCommand(f)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := f.table.(*shell.FakeTableWriter).TableOut
	for _, expected := range []string{"example | database", "example-workspace-task-a | app", "example-workspace-task-a | node"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected status to contain %q, got %q", expected, output)
		}
	}
	if strings.Contains(output, "example-workspace-task-a | database") {
		t.Errorf("did not expect shared database in workspace project, got %q", output)
	}
}

func TestStatusShowsMainAndAllActiveWorkspaces(t *testing.T) {
	f := newFakeKoolStatus()
	f.shell = &FakeRaceShell{FakeShell: shell.FakeShell{MockErrStream: io.Discard, MockOutStream: io.Discard}}
	f.env.Set("KOOL_NAME", "example")
	f.env.Set("KOOL_WORKSPACES_ENABLED", "true")
	f.env.Set("KOOL_WORKSPACE_SERVICES", "app")
	f.getServicesCmd.(*builder.FakeCommand).MockExecOut = "app\ndatabase"
	f.getProjectsCmd.(*builder.FakeCommand).MockExecOut = "example-workspace-task-b\nexample-workspace-task-a\nexample-workspace-task-b"
	f.getProjectServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"

	if err := NewStatusCommand(f).Execute(); err != nil {
		t.Fatal(err)
	}

	output := f.table.(*shell.FakeTableWriter).TableOut
	for _, expected := range []string{
		"example | app",
		"example | database",
		"example-workspace-task-a | app",
		"example-workspace-task-b | app",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected status to contain %q, got %q", expected, output)
		}
	}
	if strings.Contains(output, "example-workspace-task-a | database") || strings.Contains(output, "example-workspace-task-b | database") {
		t.Errorf("did not expect shared database in workspace projects, got %q", output)
	}
}

func TestWorkspaceServiceInfoExcludesOneOffContainers(t *testing.T) {
	f := newFakeKoolStatus()
	recorder := &statusRecordingShell{}
	f.shell = recorder
	f.getProjectServiceIDCmd.(*builder.FakeCommand).MockExecOut = "regular-id\noneoff-id\n"
	f.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Up 1 minute|80/tcp"

	running, _, _, err := f.getServiceInfo("example-workspace-task", "app")
	if err != nil {
		t.Fatal(err)
	}
	if !running {
		t.Error("expected canonical service container to be running")
	}
	var recorded []string
	for _, args := range recorder.args {
		recorded = append(recorded, args...)
	}
	if !strings.Contains(strings.Join(recorded, " "), "label=com.docker.compose.oneoff=False") {
		t.Fatalf("expected one-off container filter, got %v", recorder.args)
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
	f.env.Set("KOOL_NAME", "example")
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
