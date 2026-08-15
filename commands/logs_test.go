package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"testing"
)

func newFakeKoolLogs() *KoolLogs {
	return &KoolLogs{
		*(newDefaultKoolService().Fake()),
		&KoolLogsFlags{25, false},
		&builder.FakeCommand{MockCmd: "list", MockExecOut: "app"},
		&builder.FakeCommand{MockCmd: "logs"},
	}
}

func newFakeFailedKoolLogs() *KoolLogs {
	return &KoolLogs{
		*(newDefaultKoolService().Fake()),
		&KoolLogsFlags{25, false},
		&builder.FakeCommand{MockCmd: "list", MockExecOut: "app"},
		&builder.FakeCommand{MockCmd: "logs", MockInteractiveError: errors.New("error logs")},
	}
}

func TestNewKoolLogs(t *testing.T) {
	k := NewKoolLogs()

	if _, ok := k.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolLogs instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolLogs instance")
	} else {
		if k.Flags.Tail != 25 {
			t.Errorf("bad default value for Tail flag on default KoolLogs instance")
		}

		if k.Flags.Follow {
			t.Errorf("bad default value for Follow flag on default KoolLogs instance")
		}
	}

	if _, ok := k.logs.(*builder.DefaultCommand); !ok {
		t.Error("unexpected builder.Command on default KoolLogs instance")
	}

	if k.logs.String() != "docker compose logs" {
		t.Error("unexpected logs .String() on default KoolLogs instance logs")
	}
}

func TestNewLogsCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.logs.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolLogs.logs Command")
	}

	argsAppend := f.logs.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 2 || argsAppend[0] != "--tail" || argsAppend[1] != "25" {
		t.Errorf("bad arguments to KoolLogs.logs Command with default flags")
	}
}

func TestNewLogsTailCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	cmd.SetArgs([]string{"--tail=10"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.logs.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolLogs.logs Command")
	}

	argsAppend := f.logs.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 2 || argsAppend[0] != "--tail" || argsAppend[1] != "10" {
		t.Errorf("bad arguments to KoolLogs.logs Command when passing --tail flag")
	}
}

func TestNewLogsTailAllCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	cmd.SetArgs([]string{"--tail=0"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.logs.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolLogs.logs Command")
	}

	argsAppend := f.logs.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 2 || argsAppend[0] != "--tail" || argsAppend[1] != "all" {
		t.Errorf("bad arguments to KoolLogs.logs Command when passing 0 to --tail flag")
	}
}

func TestNewLogsFollowCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	cmd.SetArgs([]string{"--follow"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.logs.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolLogs.logs Command")
	}

	argsAppend := f.logs.(*builder.FakeCommand).ArgsAppend
	if len(argsAppend) != 3 || argsAppend[2] != "--follow" {
		t.Errorf("bad arguments to KoolLogs.logs Command when passing --follow flag")
	}
}

func TestNewLogsServiceCommand(t *testing.T) {
	f := newFakeKoolLogs()
	cmd := NewLogsCommand(f)

	cmd.SetArgs([]string{"app"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	args, ok := f.shell.(*shell.FakeShell).ArgsInteractive["logs"]
	if !ok || len(args) != 1 || args[0] != "app" {
		t.Errorf("bad arguments to KoolLogs.logs Command when executing it")
	}
}

func TestFailingNewLogsCommand(t *testing.T) {
	f := newFakeFailedKoolLogs()
	cmd := NewLogsCommand(f)

	assertExecGotError(t, cmd, "error logs")
}

func TestNoContainersNewLogsCommand(t *testing.T) {
	f := newFakeKoolLogs()
	f.list.(*builder.FakeCommand).MockExecOut = ""

	cmd := NewLogsCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledWarning {
		t.Error("did not call Warning")
	}

	expectedWarning := "There are no containers"

	if gotWarning := fmt.Sprint(f.shell.(*shell.FakeShell).WarningOutput...); gotWarning != expectedWarning {
		t.Errorf("expecting warning '%s', got '%s'", expectedWarning, gotWarning)
	}

	if val, ok := f.shell.(*shell.FakeShell).CalledInteractive["logs"]; val && ok {
		t.Error("should not call docker compose logs if there are no containers")
	}
}

func TestFailingNoContainersNewLogsCommand(t *testing.T) {
	f := newFakeKoolLogs()
	f.list.(*builder.FakeCommand).MockExecError = errors.New("error list")

	cmd := NewLogsCommand(f)

	assertExecGotError(t, cmd, "error list")
}

func TestParseLogLine(t *testing.T) {
	tests := []struct {
		line       string
		service    string
		message    string
	}{
		{"web | GET / 200", "web", "GET / 200"},
		{"web  |  GET / 200", "web", "GET / 200"},
		{"app | starting worker process", "app", "starting worker process"},
		{"plain text without delimiter", "", "plain text without delimiter"},
		{"", "", ""},
	}

	for _, tt := range tests {
		entry := parseLogLine(tt.line)
		if entry.Service != tt.service {
			t.Errorf("parseLogLine(%q): service = %q, want %q", tt.line, entry.Service, tt.service)
		}
		if entry.Message != tt.message {
			t.Errorf("parseLogLine(%q): message = %q, want %q", tt.line, entry.Message, tt.message)
		}
	}
}

func newFakeKoolLogsJSON() *KoolLogs {
	f := &KoolLogs{
		*(newDefaultKoolService().Fake()),
		&KoolLogsFlags{25, false},
		&builder.FakeCommand{MockCmd: "list", MockExecOut: "app"},
		&builder.FakeCommand{MockCmd: "logs"},
	}
	f.shell.(*shell.FakeShell).MockIsJSONOutput = true
	f.shell.(*shell.FakeShell).MockErrStream = io.Discard
	f.shell.(*shell.FakeShell).MockOutStream = io.Discard
	return f
}

func TestLogsJSONOutput(t *testing.T) {
	f := newFakeKoolLogsJSON()
	f.logs.(*builder.FakeCommand).MockExecOut = "web | GET / 200\ndb | connected to redis"

	cmd := NewLogsCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	fakeShell := f.shell.(*shell.FakeShell)
	if len(fakeShell.OutLines) != 2 {
		t.Fatalf("expected 2 JSON lines, got %d: %v", len(fakeShell.OutLines), fakeShell.OutLines)
	}

	var entry logEntryJSON
	if err := json.Unmarshal([]byte(fakeShell.OutLines[0]), &entry); err != nil {
		t.Fatalf("failed to parse first JSON line: %v", err)
	}
	if entry.Service != "web" {
		t.Errorf("expected service 'web', got '%s'", entry.Service)
	}
	if entry.Message != "GET / 200" {
		t.Errorf("expected message 'GET / 200', got '%s'", entry.Message)
	}

	if err := json.Unmarshal([]byte(fakeShell.OutLines[1]), &entry); err != nil {
		t.Fatalf("failed to parse second JSON line: %v", err)
	}
	if entry.Service != "db" {
		t.Errorf("expected service 'db', got '%s'", entry.Service)
	}
	if entry.Message != "connected to redis" {
		t.Errorf("expected message 'connected to redis', got '%s'", entry.Message)
	}
}

func TestLogsJSONNoContainers(t *testing.T) {
	f := newFakeKoolLogsJSON()
	f.list.(*builder.FakeCommand).MockExecOut = ""

	cmd := NewLogsCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing logs command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledWarning {
		t.Error("expected Warning to be called when no containers")
	}
}
