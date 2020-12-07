package parser

import (
	"kool-dev/kool/cmd/builder"
	"os"
	"path"
	"testing"
)

func TestDefaultParser(t *testing.T) {
	var p Parser = NewParser()

	if _, assert := p.(*DefaultParser); !assert {
		t.Errorf("NewParser() did not return a *DefaultParser")
	}
}

func TestParserAddLookupPath(t *testing.T) {
	var (
		p      Parser = NewParser()
		err    error
		tmpDir = t.TempDir()
	)

	err = p.AddLookupPath(tmpDir)

	if err == nil || ErrKoolYmlNotFound.Error() != err.Error() {
		t.Errorf("expected ErrKoolYmlNotFound; got %s", err)
	}

	workDir, _ := os.Getwd()
	if err = p.AddLookupPath(path.Join(workDir, "testing_files")); err != nil {
		t.Errorf("unexpected error; error: %s", err)
	}

	_ = p.AddLookupPath(path.Join(workDir, "testing_files"))
	_ = p.AddLookupPath(path.Join(workDir, "testing_files"))

	if commands, _ := p.Parse("testing"); len(commands) != 1 {
		t.Errorf("expecting to get only one command, got '%v'", len(commands))
	}
}

func TestParserParse(t *testing.T) {
	var (
		p        Parser = NewParser()
		commands []builder.Command
		err      error
	)

	if _, err = p.Parse("testing"); err == nil {
		t.Error("expecting 'kool.yml not found' error, got none")
	}

	if err != nil && err.Error() != "kool.yml not found" {
		t.Errorf("expecting error 'kool.yml not found', got '%s'", err.Error())
	}

	workDir, _ := os.Getwd()
	_ = p.AddLookupPath(path.Join(workDir, "testing_files"))

	if commands, err = p.Parse("testing"); err != nil {
		t.Errorf("unexpected error; error: %s", err)
	}

	if len(commands) != 1 || commands[0].String() != "echo testing" {
		t.Error("failed to parse testing kool.yml")
	}

	if commands, err = p.Parse("invalid"); err != nil {
		t.Errorf("unexpected error; error: %s", err)
	}

	if len(commands) > 0 {
		t.Error("should not find scripts")
	}
}

func TestParserParseAvailableScripts(t *testing.T) {
	var (
		p       Parser = NewParser()
		scripts []string
		err     error
	)

	if _, err = p.ParseAvailableScripts(""); err == nil {
		t.Error("expecting 'kool.yml not found' error, got none")
	}

	if err != nil && err.Error() != "kool.yml not found" {
		t.Errorf("expecting error 'kool.yml not found', got '%s'", err.Error())
	}

	workDir, _ := os.Getwd()
	_ = p.AddLookupPath(path.Join(workDir, "testing_files"))

	if scripts, err = p.ParseAvailableScripts(""); err != nil {
		t.Errorf("unexpected error; error: %s", err)
	}

	if len(scripts) != 2 || scripts[0] != "testing" || scripts[1] != "variables" {
		t.Error("failed to get all scripts from kool.yml")
	}
}

func TestParserParseAvailableScriptsFilter(t *testing.T) {
	var (
		p       Parser = NewParser()
		scripts []string
		err     error
	)

	workDir, _ := os.Getwd()
	_ = p.AddLookupPath(path.Join(workDir, "testing_files"))

	if scripts, err = p.ParseAvailableScripts("tes"); err != nil {
		t.Errorf("unexpected error; error: %s", err)
	}

	if len(scripts) != 1 || scripts[0] != "testing" {
		t.Error("failed to get filtered scripts from kool.yml")
	}

	if scripts, err = p.ParseAvailableScripts("invalidFilter"); err != nil {
		t.Errorf("unexpected error; error: %s", err)
	}

	if len(scripts) != 0 {
		t.Error("failed to get filtered scripts from kool.yml")
	}
}

func TestParserLookUpVariables(t *testing.T) {
	var variables []string
	p := NewParser()

	if variables = p.LookUpVariables("variables"); len(variables) > 0 {
		t.Error("should not find variables without looking for kool.yml")
	}

	workDir, _ := os.Getwd()
	_ = p.AddLookupPath(path.Join(workDir, "testing_files"))

	if variables = p.LookUpVariables("invalid"); len(variables) > 0 {
		t.Error("should not find variables in an invalid script")
	}

	if variables = p.LookUpVariables("testing"); len(variables) > 0 {
		t.Error("should not find variables in a script without variables")
	}

	if variables = p.LookUpVariables("variables"); len(variables) != 1 || variables[0] != "testing_var" {
		t.Error("failed to get variables in a script")
	}
}
