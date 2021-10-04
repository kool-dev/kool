package commands

import (
	"bytes"
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"os"
	"testing"
)

func newFakeKoolDocker() *KoolDocker {
	return &KoolDocker{
		*(newDefaultKoolService().Fake()),
		&KoolDockerFlags{false, []string{}, []string{}, []string{}, []string{}},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "docker"},
	}
}

func newFailedFakeKoolDocker() *KoolDocker {
	return &KoolDocker{
		*(newDefaultKoolService().Fake()),
		&KoolDockerFlags{false, []string{}, []string{}, []string{}, []string{}},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "docker", MockInteractiveError: errors.New("error docker")},
	}
}

func TestNewKoolDocker(t *testing.T) {
	k := NewKoolDocker()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolDocker instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolDocker instance")
	} else {
		if k.Flags.DisableTty {
			t.Errorf("bad default value for DisableTty flag on default KoolDocker instance")
		}

		if len(k.Flags.EnvVariables) > 0 {
			t.Errorf("bad default value for EnvVariables flag on default KoolDocker instance")
		}

		if len(k.Flags.Volumes) > 0 {
			t.Errorf("bad default value for Volumes flag on default KoolDocker instance")
		}

		if len(k.Flags.Publish) > 0 {
			t.Errorf("bad default value for Publish flag on default KoolDocker instance")
		}
	}

	if _, ok := k.dockerRun.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Command on default KoolDocker instance")
	}
}

func TestNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()

	cmd := NewDockerCommand(f)
	workDir, _ := os.Getwd()

	cmd.SetArgs([]string{"image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	if !f.dockerRun.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolDocker.dockerRun Command")
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 3 || argsAppend[0] != "-t" || argsAppend[1] != "--volume" || argsAppend[2] != workDir+":/app:delegated" {
		t.Errorf("bad arguments to KoolDocker.dockerRun Command with default flags: %v", argsAppend)
	}

	if val, ok := f.shell.(*shell.FakeShell).CalledInteractive["docker"]; !ok || !val {
		t.Errorf("did not call Interactive on KoolDocker.dockerRun Command")
	}

	interactiveArgs, ok := f.shell.(*shell.FakeShell).ArgsInteractive["docker"]

	if !ok || len(interactiveArgs) != 1 || interactiveArgs[0] != "image" {
		t.Errorf("bad arguments to Interactive on KoolDocker.dockerRun Command")
	}
}

func TestNoArgsNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false

	cmd := NewDockerCommand(f)

	cmd.SetOut(bytes.NewBufferString(""))

	if err := cmd.Execute(); err == nil {
		t.Error("expecting no arguments error executing docker command")
	}
}

func TestAsUserEnvKoolImageNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false
	cmd := NewDockerCommand(f)

	f.envStorage.(*environment.FakeEnvStorage).Envs["KOOL_ASUSER"] = "kooldev_user_test"

	cmd.SetArgs([]string{"kooldev/image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[0] != "--env" || argsAppend[1] != "ASUSER=kooldev_user_test" {
		t.Errorf("bad arguments to KoolDocker.dockerRun Command with 'KOOL_ASUSER' variable")
	}
}

func TestAsUserEnvFireworkImageNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false
	cmd := NewDockerCommand(f)

	f.envStorage.(*environment.FakeEnvStorage).Envs["KOOL_ASUSER"] = "kooldev_user_test"

	cmd.SetArgs([]string{"fireworkweb/image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[0] != "--env" || argsAppend[1] != "ASUSER=kooldev_user_test" {
		t.Errorf("bad arguments to KoolDocker.dockerRun Command with 'KOOL_ASUSER' variable")
	}
}

func TestEnvFlagNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"--env=VAR_TEST=1", "image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[0] != "--env" || argsAppend[1] != "VAR_TEST=1" {
		t.Errorf("bad arguments to KoolDocker.dockerRun Command with EnvVariables flag")
	}
}

func TestVolumesFlagNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"--volume=volume_test", "image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[2] != "--volume" || argsAppend[3] != "volume_test" {
		t.Errorf("bad arguments to KoolDocker.dockerRun Command with Volumes flag")
	}
}

func TestPublishFlagNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"--publish=publish_test", "image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[2] != "--publish" || argsAppend[3] != "publish_test" {
		t.Errorf("bad arguments to KoolDocker.dockerRun Command with Volumes flag")
	}
}

func TestNetworkFlagNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"--network=kool_global", "image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[2] != "--network" || argsAppend[3] != "kool_global" {
		t.Errorf("bad arguments to KoolDocker.dockerRun Command with Network flag")
	}
}

func TestImageCommandsNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"image", "command1", "command2"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	interactiveArgs, ok := f.shell.(*shell.FakeShell).ArgsInteractive["docker"]

	if !ok || len(interactiveArgs) != 3 || interactiveArgs[0] != "image" || interactiveArgs[1] != "command1" || interactiveArgs[2] != "command2" {
		t.Errorf("bad arguments to Interactive on KoolDocker.dockerRun Command")
	}
}

func TestFailingNewDockerCommand(t *testing.T) {
	f := newFailedFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"image"})

	assertExecGotError(t, cmd, "error docker")
}

func TestNonTerminalNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	f.shell.(*shell.FakeShell).MockIsTerminal = false

	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 2 || argsAppend[0] == "-t" {
		t.Errorf("bad arguments to KoolDocker.dockerRun Command on non terminal environment")
	}
}
