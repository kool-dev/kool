package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"

	"github.com/spf13/cobra"
)

func TestStatusJSONEmptyAndSorted(t *testing.T) {
	for _, services := range []string{"", "zebra\napp"} {
		f := newFakeKoolStatusJSON()
		f.getServicesCmd.(*builder.FakeCommand).MockExecOut = services
		if err := NewStatusCommand(f).Execute(); err != nil {
			t.Fatal(err)
		}
		var output statusOutputJSON
		if err := json.Unmarshal([]byte(f.shell.(*shell.FakeShell).OutLines[0]), &output); err != nil {
			t.Fatal(err)
		}
		if output.Services == nil || (services != "" && output.Services[0].Service != "app") {
			t.Fatalf("unexpected services: %#v", output.Services)
		}
	}
}

func TestOutputFormatValidation(t *testing.T) {
	for _, format := range []string{"JSON", "text", "jsno"} {
		env := environment.NewFakeEnvStorage()
		root := NewRootCmd(env)
		root.AddCommand(&cobra.Command{Use: "test", Run: func(*cobra.Command, []string) {}})
		root.SetArgs([]string{"--output", format, "test"})
		err := root.Execute()
		if format == "jsno" {
			if err == nil {
				t.Fatal("expected invalid format error")
			}
		} else if err != nil || env.Get("KOOL_OUTPUT") != strings.ToLower(format) {
			t.Fatalf("format %s: %v", format, err)
		}
	}
}

func TestLogsJSONLongLinesAndSeparateDiagnostics(t *testing.T) {
	original := execLogsCmd
	t.Cleanup(func() { execLogsCmd = original })
	for _, follow := range []bool{false, true} {
		f := newFakeKoolLogsJSON()
		f.Flags.Follow = follow
		var output, diagnostics bytes.Buffer
		f.shell.(*shell.FakeShell).MockOutStream = &output
		f.shell.(*shell.FakeShell).MockErrStream = &diagnostics
		execLogsCmd = func(string, ...string) *exec.Cmd {
			return exec.Command("sh", "-c", "printf 'warning\\n' >&2; printf 'app | '; head -c 131072 /dev/zero | tr '\\000' x; printf '\\napp | done'")
		}
		if err := NewLogsCommand(f).Execute(); err != nil {
			t.Fatal(err)
		}
		decoder := json.NewDecoder(&output)
		var entry logEntryJSON
		if err := decoder.Decode(&entry); err != nil || len(entry.Message) != 131072 {
			t.Fatalf("long log lost: %d bytes, %v", len(entry.Message), err)
		}
		if err := decoder.Decode(&entry); err != nil || entry.Message != "done" {
			t.Fatalf("last log lost: %v", err)
		}
		if diagnostics.String() != "warning\n" {
			t.Fatalf("stderr lost: %q", diagnostics.String())
		}
	}
}

func TestRunJSONPreservesParseError(t *testing.T) {
	parseErr := errors.New("invalid YAML document")
	f := newFakeKoolRun(nil, map[string]error{"build": parseErr})
	f.shell.(*shell.FakeShell).MockIsJSONOutput = true
	var diagnostics bytes.Buffer
	f.shell.(*shell.FakeShell).MockErrStream = &diagnostics
	cmd := NewRunCommand(f)
	cmd.SetArgs([]string{"build"})
	err := cmd.Execute()
	if !errors.Is(err, parseErr) {
		t.Fatalf("lost original error: %v", err)
	}
	var output struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(diagnostics.Bytes(), &output); err != nil || output.Error != parseErr.Error() {
		t.Fatalf("incorrect diagnostic: %s (%v)", diagnostics.String(), err)
	}
}
