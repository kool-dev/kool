package parser

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"testing"
)

func TestFakeKooYml(t *testing.T) {
	f := &FakeKoolYaml{}

	f.MockParseError = map[string]error{
		"path": errors.New("parse error"),
	}

	if err := f.Parse("path"); err == nil {
		t.Error("expecting error on Parse, got none")
	} else if err.Error() != "parse error" {
		t.Errorf("expecting error 'parse error' on Parse, got %v", err)
	}

	if val, ok := f.CalledParse["path"]; !val || !ok {
		t.Error("failed calling Parse")
	}

	f.MockHasScript = map[string]bool{"script": true}

	if has := f.HasScript("script"); !has {
		t.Error("expecting to get true on HasScript, got false")
	}

	if val, ok := f.CalledHasScript["script"]; !val || !ok {
		t.Error("failed calling HasScript")
	}

	f.MockCommands = map[string][]builder.Command{
		"script": []builder.Command{
			&builder.FakeCommand{MockCmd: "command"},
		},
	}
	f.MockParseCommandsError = map[string]error{
		"script": errors.New("parse commands error"),
	}

	parsedCmds, parsedError := f.ParseCommands("script")

	if parsedError == nil {
		t.Error("expecting error 'parse commands error', got none")
	} else if parsedError.Error() != "parse commands error" {
		t.Errorf("expecting error 'parse commands error', got %v", parsedError)
	}

	if len(parsedCmds) != 1 || parsedCmds[0].Cmd() != "command" {
		t.Error("failed getting mocked commands")
	}

	if val, ok := f.CalledParseCommands["script"]; !val || !ok {
		t.Error("failed calling ParseCommands")
	}

	f.SetScript("script", []string{"command"})

	if val, ok := f.CalledSetScript["script"]; !val || !ok {
		t.Error("failed calling SetScript")
	}

	if cmds, ok := f.ScriptCommands["script"]; !ok || len(cmds) != 1 {
		t.Error("failed setting script")
	}

	f.MockString = "string"
	f.MockStringError = errors.New("string error")

	stringContent, stringError := f.String()

	if stringError == nil {
		t.Error("expecting error 'string error', got none")
	} else if stringError.Error() != "string error" {
		t.Errorf("expecting error 'string error', got %v", parsedError)
	}

	if stringContent != "string" {
		t.Errorf("expecting to get 'string' on String, got %s", stringContent)
	}
}
