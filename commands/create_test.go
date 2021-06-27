package commands

import (
	"bytes"
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/presets"
	"kool-dev/kool/core/shell"
	"strings"
	"testing"
)

func newFakeKoolCreate() *KoolCreate {
	return &KoolCreate{
		*newFakeKoolService(),
		&presets.FakeParser{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
		*newFakeKoolPreset(),
	}
}

func TestNewKoolCreate(t *testing.T) {
	k := NewKoolCreate()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolCreate instance")
	}

	if _, ok := k.DefaultKoolService.term.(*shell.DefaultTerminalChecker); !ok {
		t.Errorf("unexpected shell.TerminalChecker on default KoolCreate instance")
	}

	if _, ok := k.createCommand.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Command on default KoolCreate instance")
	}

	if _, ok := k.parser.(*presets.DefaultParser); !ok {
		t.Errorf("unexpected presets.Parser on default KoolCreate instance")
	}
}

func TestNewKoolCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	f.parser.(*presets.FakeParser).MockExists = true
	f.KoolPreset.presetsParser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": {
			Commands: map[string][]string{
				"create": {"kool docker create command"},
			},
		},
	}
	f.KoolPreset.presetsParser.(*presets.FakeParser).MockConfig = f.parser.(*presets.FakeParser).MockConfig
	f.createCommand.(*builder.FakeCommand).MockCmd = "create"

	cmd := NewCreateCommand(f)
	cmd.SetArgs([]string{"laravel", "my-app"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing create command; error: %v", err)
	}

	if !f.parser.(*presets.FakeParser).CalledLoadPresets {
		t.Error("did not call parser.LoadPresets")
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if val, ok := f.parser.(*presets.FakeParser).CalledGetConfig["laravel"]; !ok || !val {
		t.Error("did not call parser.GetConfig for preset 'laravel'")
	}

	if !f.createCommand.(*builder.FakeCommand).CalledParseCommand {
		t.Error("did not call Parse on KoolCreate.createCommand Command")
	}

	if val, ok := f.shell.(*shell.FakeShell).CalledInteractive["create"]; !val || !ok {
		t.Error("did not call Interactive on KoolCreate.createCommand Command")
	}
}

func TestInvalidPresetCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()
	cmd := NewCreateCommand(f)

	cmd.SetArgs([]string{"invalid", "my-app"})

	if err := cmd.Execute(); err == nil {
		t.Error("should have got an error")
	} else if !strings.Contains(err.Error(), "unknown preset") {
		t.Errorf("unexpected error: %s", err)
	}

	if !f.parser.(*presets.FakeParser).CalledLoadPresets {
		t.Error("did not call parser.LoadPresets")
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}
}

func TestNoArgsNewCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	cmd := NewCreateCommand(f)
	cmd.SetOut(bytes.NewBufferString(""))

	if err := cmd.Execute(); err == nil {
		t.Error("expecting no arguments error executing create command")
	}
}

func TestErrorConfigCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	f.parser.(*presets.FakeParser).MockExists = true
	getConfigError := errors.New("get config error")
	f.parser.(*presets.FakeParser).MockGetConfigError = map[string]error{
		"laravel": getConfigError,
	}

	cmd := NewCreateCommand(f)

	cmd.SetArgs([]string{"laravel", "my-app"})

	assertExecGotError(t, cmd, "error parsing preset config")
}

func TestNoCreateCommandsCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockConfig = map[string]*presets.PresetConfig{
		"laravel": {
			Commands: make(map[string][]string),
		},
	}

	cmd := NewCreateCommand(f)

	cmd.SetArgs([]string{"laravel", "my-app"})

	assertExecGotError(t, cmd, "no create commands were found for preset laravel")
}
