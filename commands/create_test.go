package commands

import (
	"bytes"
	"errors"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/presets"
	"kool-dev/kool/core/shell"
	"os"
	"path/filepath"
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

	if _, ok := k.shell.(*shell.DefaultShell); !ok {
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
	cmd.SetArgs([]string{"laravel", "/tmp"})

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

	cwd, _ := os.Getwd()

	cmd.SetArgs([]string{"laravel", t.TempDir()})

	assertExecGotError(t, cmd, "install error")

	// return to original folder
	_ = os.Chdir(cwd)
}

func TestCreateCommandUsesRelativeCreateDirectory(t *testing.T) {
	f := newFakeKoolCreate()
	f.parser.(*presets.FakeParser).MockExists = true

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	parent := t.TempDir()
	dest := filepath.Join(parent, "my-app")
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatal(err)
	}

	cmd := NewCreateCommand(f)
	cmd.SetArgs([]string{"laravel", dest})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing create command; error: %v", err)
	}

	if got := f.env.Get("CREATE_DIRECTORY"); got != "my-app" {
		t.Errorf("CREATE_DIRECTORY should be the folder name for docker mounts, got %q", got)
	}
}
