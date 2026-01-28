package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/parser"
	"kool-dev/kool/core/shell"
	"strings"
	"testing"
)

func newFakeKoolScripts(mockScripts []string, mockParseErr error) *KoolScripts {
	var details []parser.ScriptDetail
	for _, script := range mockScripts {
		details = append(details, parser.ScriptDetail{Name: script, Comments: []string{}, Commands: []string{}})
	}
	return &KoolScripts{
		*(newDefaultKoolService().Fake()),
		&KoolScriptsFlags{},
		&parser.FakeParser{
			MockScripts:                    mockScripts,
			MockScriptDetails:              details,
			MockParseAvailableScriptsError: mockParseErr,
		},
		environment.NewFakeEnvStorage(),
	}
}

func TestScriptsCommandListsScripts(t *testing.T) {
	f := newFakeKoolScripts([]string{"setup", "lint"}, nil)
	cmd := NewScriptsCommand(f)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing scripts command; error: %v", err)
	}

	if !f.parser.(*parser.FakeParser).CalledParseAvailableScripts {
		t.Errorf("did not call ParseAvailableScripts")
	}

	fakeShell := f.shell.(*shell.FakeShell)

	if !containsLine(fakeShell.OutLines, "setup") || !containsLine(fakeShell.OutLines, "lint") {
		t.Errorf("missing scripts on output: %v", fakeShell.OutLines)
	}
}

func TestScriptsCommandFiltersScripts(t *testing.T) {
	f := newFakeKoolScripts([]string{"setup", "lint"}, nil)
	cmd := NewScriptsCommand(f)

	cmd.SetArgs([]string{"se"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing scripts command; error: %v", err)
	}

	fakeShell := f.shell.(*shell.FakeShell)

	if containsLine(fakeShell.OutLines, "lint") {
		t.Errorf("unexpected script on output: %v", fakeShell.OutLines)
	}

	if !containsLine(fakeShell.OutLines, "setup") {
		t.Errorf("missing filtered script on output: %v", fakeShell.OutLines)
	}
}

func TestScriptsCommandNoScripts(t *testing.T) {
	f := newFakeKoolScripts([]string{}, nil)
	cmd := NewScriptsCommand(f)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing scripts command; error: %v", err)
	}

	fakeShell := f.shell.(*shell.FakeShell)

	if !fakeShell.CalledWarning {
		t.Errorf("did not warn about missing scripts")
	}

	if !strings.Contains(fmt.Sprint(fakeShell.WarningOutput...), "No scripts found") {
		t.Errorf("unexpected warning output: %v", fakeShell.WarningOutput)
	}
}

func TestScriptsCommandParseError(t *testing.T) {
	f := newFakeKoolScripts([]string{"setup"}, errors.New("parse error"))
	cmd := NewScriptsCommand(f)
	cmd.SetArgs([]string{})

	assertExecGotError(t, cmd, "parse error")
}

func TestScriptsCommandJsonOutput(t *testing.T) {
	parserMock := &parser.FakeParser{
		MockScriptDetails: []parser.ScriptDetail{
			{
				Name:     "setup",
				Comments: []string{"Sets up dependencies"},
				Commands: []string{"kool run composer install"},
			},
			{
				Name:     "lint",
				Comments: []string{},
				Commands: []string{"kool run go:linux fmt ./..."},
			},
		},
	}

	f := newFakeKoolScripts([]string{}, nil)
	f.parser = parserMock
	cmd := NewScriptsCommand(f)
	cmd.SetArgs([]string{"--json"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing scripts command; error: %v", err)
	}

	fakeShell := f.shell.(*shell.FakeShell)

	if len(fakeShell.OutLines) == 0 {
		t.Errorf("expected JSON output")
		return
	}

	var output []parser.ScriptDetail
	if err := json.Unmarshal([]byte(fakeShell.OutLines[0]), &output); err != nil {
		t.Fatalf("failed to parse json output: %v", err)
	}

	if len(output) != 2 {
		t.Fatalf("expected 2 script entries, got %d", len(output))
	}

	if output[0].Name != "lint" || output[1].Name != "setup" {
		t.Errorf("unexpected scripts order or names: %v", output)
	}
}

func containsLine(lines []string, match string) bool {
	for _, line := range lines {
		if strings.Contains(line, match) {
			return true
		}
	}

	return false
}
