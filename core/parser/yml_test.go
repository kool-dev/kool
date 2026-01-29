package parser

import (
	"errors"
	"kool-dev/kool/core/builder"
	"os"
	"path"
	"testing"

	"gopkg.in/yaml.v3"
)

const KoolYmlOK = `scripts:
  single-line: single line script
  multi-line:
    - line 1
    - line 2
`

func TestParseKoolYaml(t *testing.T) {
	var (
		err         error
		tmpPath     string
		parsed      *KoolYaml
		cmds        []builder.Command
		koolContent string
	)

	tmpPath = path.Join(t.TempDir(), "kool.yml")

	err = os.WriteFile(tmpPath, []byte(KoolYmlOK), os.ModePerm)

	if err != nil {
		t.Fatal("failed creating temporary file for test", err)
	}

	parsed, err = ParseKoolYaml(tmpPath)

	if err != nil {
		t.Errorf("failed parsing proper kool.yml file; error: %s", err)
		return
	}

	if len(parsed.Scripts) != 2 {
		t.Errorf("expected to parse 2 scripts; got %d", len(parsed.Scripts))
		return
	}

	if !parsed.HasScript("single-line") || !parsed.HasScript("multi-line") {
		t.Errorf("expected to have single-line and multi-line script")
		return
	}

	if hasSimilars, similars := parsed.GetSimilars("single-lne"); !hasSimilars || len(similars) != 1 {
		t.Errorf("unexpected return on GetSimilars %v - %v", hasSimilars, similars)
	}

	if cmds, err = parsed.ParseCommands("single-line"); err != nil {
		t.Errorf("failed to parse proper single-line; error: %s", err)
		return
	}

	if len(cmds) != 1 {
		t.Errorf("expected single-line to parse 1 command; got %d", len(cmds))
		return
	}

	if cmds, err = parsed.ParseCommands("multi-line"); err != nil {
		t.Errorf("failed to parse proper multi-line; error: %s", err)
		return
	}

	if len(cmds) != 2 {
		t.Errorf("expected multi-line to parse 1 command; got %d", len(cmds))
		return
	}

	parsed.SetScript("new-script", []string{"new-command 1"})

	if len(parsed.Scripts) != 3 {
		t.Errorf("expected to get 3 scripts after setting a new one; got %d", len(parsed.Scripts))
		return
	}

	if !parsed.HasScript("new-script") {
		t.Errorf("expected to have new-script script")
		return
	}

	if cmds, err = parsed.ParseCommands("new-script"); err != nil {
		t.Errorf("failed to parse proper new-script; error: %s", err)
		return
	}

	if len(cmds) != 1 {
		t.Errorf("expected new-script to parse 1 command; got %d", len(cmds))
		return
	}

	parsed.SetScript("new-script", []string{"new-command 1", "new-command 2"})

	if len(parsed.Scripts) != 3 {
		t.Errorf("expected to get 3 scripts after setting a existing one; got %d", len(parsed.Scripts))
		return
	}

	if !parsed.HasScript("new-script") {
		t.Errorf("expected to have new-script script")
		return
	}

	if cmds, err = parsed.ParseCommands("new-script"); err != nil {
		t.Errorf("failed to parse proper new-script; error: %s", err)
		return
	}

	if len(cmds) != 2 {
		t.Errorf("expected new-script to parse 2 commands; got %d", len(cmds))
		return
	}

	if koolContent, err = parsed.String(); err != nil {
		t.Errorf("failed to get kool.yml content; error: %s", err)
		return
	}

	parsedOutput := new(KoolYaml)
	if err = yaml.Unmarshal([]byte(koolContent), parsedOutput); err != nil {
		t.Errorf("failed to parse generated kool.yml content; error: %s", err)
		return
	}

	if len(parsedOutput.Scripts) != 3 {
		t.Errorf("expected to parse 3 scripts from generated content; got %d", len(parsedOutput.Scripts))
		return
	}

	if !parsedOutput.HasScript("single-line") || !parsedOutput.HasScript("multi-line") || !parsedOutput.HasScript("new-script") {
		t.Errorf("expected to have single-line, multi-line and new-script scripts")
		return
	}

	if cmds, err = parsedOutput.ParseCommands("new-script"); err != nil {
		t.Errorf("failed to parse new-script from generated content; error: %s", err)
		return
	}

	if len(cmds) != 2 {
		t.Errorf("expected new-script to parse 2 commands; got %d", len(cmds))
		return
	}
}

func TestParseKoolYamlScriptDetails(t *testing.T) {
	const KoolYmlWithComments = `scripts:
  # build the app
  # and setup
  setup: kool run go build
  lint:
    - kool run go vet
    - kool run go fmt
`

	var (
		err    error
		parsed *KoolYaml
		tmp    string
	)

	tmp = path.Join(t.TempDir(), "kool.yml")
	if err = os.WriteFile(tmp, []byte(KoolYmlWithComments), os.ModePerm); err != nil {
		t.Fatal("failed creating temporary file for test", err)
	}

	if parsed, err = ParseKoolYamlWithDetails(tmp); err != nil {
		t.Fatalf("failed parsing kool.yml with comments; error: %s", err)
	}

	setup, ok := parsed.ScriptDetails["setup"]
	if !ok {
		t.Fatal("expected to find setup script details")
	}

	if len(setup.Comments) != 2 {
		t.Fatalf("expected 2 comments for setup, got %v", setup.Comments)
	}

	if setup.Comments[0] != "build the app" || setup.Comments[1] != "and setup" {
		t.Fatalf("unexpected setup comments: %v", setup.Comments)
	}

	if setup.Commands == nil || len(setup.Commands) != 1 || setup.Commands[0] != "kool run go build" {
		t.Fatalf("unexpected setup commands: %v", setup.Commands)
	}

	lint, ok := parsed.ScriptDetails["lint"]
	if !ok {
		t.Fatal("expected to find lint script details")
	}

	if len(lint.Comments) != 0 {
		t.Fatalf("expected no comments for lint, got %v", lint.Comments)
	}

	if len(lint.Commands) != 2 {
		t.Fatalf("expected 2 commands for lint, got %v", lint.Commands)
	}
}

func TestParseKoolYamlStruct(t *testing.T) {
	var (
		err     error
		tmpPath string
		parsed  *KoolYaml
	)

	tmpPath = path.Join(t.TempDir(), "kool.yml")

	err = os.WriteFile(tmpPath, []byte(KoolYmlOK), os.ModePerm)

	if err != nil {
		t.Fatal("failed creating temporary file for test", err)
	}

	parsed = new(KoolYaml)

	err = parsed.Parse(tmpPath)

	if err != nil {
		t.Errorf("failed parsing proper kool.yml file; error: %s", err)
		return
	}

	if len(parsed.Scripts) != 2 {
		t.Errorf("expected to parse 2 scripts; got %d", len(parsed.Scripts))
		return
	}
}

func TestErrorParseKoolYamlStruct(t *testing.T) {
	var (
		err     error
		tmpPath string
		parsed  *KoolYaml
	)

	tmpPath = path.Join(t.TempDir(), "kool.yml")

	invalidKoolYml := "	invalid"

	err = os.WriteFile(tmpPath, []byte(invalidKoolYml), os.ModePerm)

	if err != nil {
		t.Fatal("failed creating temporary file for test", err)
	}

	parsed = new(KoolYaml)

	if err = parsed.Parse(tmpPath); err == nil {
		t.Error("expecting error on Parse, got none")
		return
	}
}

func TestSetScriptEmptyCommandsKoolYmlParser(t *testing.T) {
	parsed := new(KoolYaml)
	var emptyCommands []string

	parsed.SetScript("key", []string{"command", "another-command"})

	parsed.SetScript("key", emptyCommands)

	if commands := parsed.Scripts["key"]; len(commands.([]interface{})) == 0 {
		t.Error("calling SetScript with no command should no override existing commands")
	}
}

func TestErrorStringKoolYmlParser(t *testing.T) {
	originalYamlMarshalFn := yamlMarshalFn

	defer func() {
		yamlMarshalFn = originalYamlMarshalFn
	}()

	yamlMarshalFn = func(in interface{}) ([]byte, error) {
		return nil, errors.New("marshal error")
	}

	parsed := new(KoolYaml)

	_, err := parsed.String()

	if err == nil {
		t.Error("expecting an error on String, got none")
	} else if err.Error() != "marshal error" {
		t.Errorf("expecting error 'marshal error' on String, got '%v'", err)
	}
}
