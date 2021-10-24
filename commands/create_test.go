package commands

import (
	"bytes"
	"errors"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/presets"
	"kool-dev/kool/core/shell"
	"strings"
	"testing"
)

func newFakeKoolCreate() *KoolCreate {
	return &KoolCreate{
		*(newDefaultKoolService().Fake()),
		&presets.FakeParser{},
		environment.NewFakeEnvStorage(),
	}
}

func TestNewKoolCreate(t *testing.T) {
	k := NewKoolCreate()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolCreate instance")
	}

	if _, ok := k.parser.(*presets.DefaultParser); !ok {
		t.Errorf("unexpected presets.Parser on default KoolCreate instance")
	}
}

func TestNewKoolCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockCreate = nil
	f.parser.(*presets.FakeParser).MockInstall = nil

	cmd := NewCreateCommand(f)
	cmd.SetArgs([]string{"laravel", "my-app"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing create command; error: %v", err)
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.parser.(*presets.FakeParser).CalledCreate {
		t.Error("did not call parser.Create")
	}

	if !f.parser.(*presets.FakeParser).CalledInstall {
		t.Error("did not call parser.Install")
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

func TestErrCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	f.parser.(*presets.FakeParser).MockExists = true
	createErr := errors.New("create error")
	f.parser.(*presets.FakeParser).MockCreate = createErr

	cmd := NewCreateCommand(f)

	cmd.SetArgs([]string{"laravel", "my-app"})

	assertExecGotError(t, cmd, "create error")
}

func TestErrInstallCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockInstall = errors.New("install error")

	cmd := NewCreateCommand(f)

	cmd.SetArgs([]string{"laravel", "my-app"})

	assertExecGotError(t, cmd, "install error")
}
