package cmd

import (
	"bytes"
	"errors"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/shell"
	"os"
	"testing"
)

func newFakeKoolDocker() *KoolDocker {
	return &KoolDocker{
		*newFakeKoolService(),
		&KoolDockerFlags{false, []string{}, []string{}, []string{}},
		&builder.FakeCommand{},
	}
}

func newFailedFakeKoolDocker() *KoolDocker {
	return &KoolDocker{
		*newFakeKoolService(),
		&KoolDockerFlags{false, []string{}, []string{}, []string{}},
		&builder.FakeFailedCommand{MockError: errors.New("error docker")},
	}
}

func TestNewKoolDocker(t *testing.T) {
	k := NewKoolDocker()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Errorf("unexpected shell.OutputWriter on default KoolDocker instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolDocker instance")
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

	if !f.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Errorf("did not call SetWriter")
	}

	if !f.dockerRun.(*builder.FakeCommand).CalledAppendArgs {
		t.Errorf("did not call AppendArgs on KoolDocker.dockerRun Command")
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 3 || argsAppend[0] != "-t" || argsAppend[1] != "--volume" || argsAppend[2] != workDir+":/app:delegated" {
		t.Errorf("bad arguments to KoolDocker.logs Command with default flags")
	}

	if !f.dockerRun.(*builder.FakeCommand).CalledInteractive {
		t.Errorf("did not call Interactive on KoolDocker.dockerRun Command")
	}

	interactiveArgs := f.dockerRun.(*builder.FakeCommand).ArgsInteractive

	if len(interactiveArgs) != 1 || interactiveArgs[0] != "image" {
		t.Errorf("bad arguments to Interactive on KoolDocker.dockerRun Command")
	}
}

func TestNoArgsNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	cmd.SetOut(bytes.NewBufferString(""))

	if err := cmd.Execute(); err == nil {
		t.Error("expecting no arguments error executing docker command")
	}
}

func TestDisableTTYFlagNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"--disable-tty", "image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 2 || argsAppend[0] == "-t" {
		t.Errorf("bad arguments to KoolDocker.logs Command with disable-tty flag")
	}
}

func TestDisableTTYEnvNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	os.Setenv("KOOL_TTY_DISABLE", "1")
	defer os.Unsetenv("KOOL_TTY_DISABLE")

	cmd.SetArgs([]string{"image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 2 || argsAppend[0] == "-t" {
		t.Errorf("bad arguments to KoolDocker.logs Command with 'KOOL_TTY_DISABLE' variable")
	}
}

func TestAsUserEnvKoolImageNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	os.Setenv("KOOL_ASUSER", "kooldev_user_test")
	defer os.Unsetenv("KOOL_ASUSER")

	cmd.SetArgs([]string{"--disable-tty", "kooldev/image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[0] != "--env" || argsAppend[1] != "ASUSER=kooldev_user_test" {
		t.Errorf("bad arguments to KoolDocker.logs Command with 'KOOL_ASUSER' variable")
	}
}

func TestAsUserEnvFireworkImageNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	os.Setenv("KOOL_ASUSER", "kooldev_user_test")
	defer os.Unsetenv("KOOL_ASUSER")

	cmd.SetArgs([]string{"--disable-tty", "fireworkweb/image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[0] != "--env" || argsAppend[1] != "ASUSER=kooldev_user_test" {
		t.Errorf("bad arguments to KoolDocker.logs Command with 'KOOL_ASUSER' variable")
	}
}

func TestEnvFlagNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"--disable-tty", "--env=VAR_TEST=1", "image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[0] != "--env" || argsAppend[1] != "VAR_TEST=1" {
		t.Errorf("bad arguments to KoolDocker.logs Command with EnvVariables flag")
	}
}

func TestVolumesFlagNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"--disable-tty", "--volume=volume_test", "image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[2] != "--volume" || argsAppend[3] != "volume_test" {
		t.Errorf("bad arguments to KoolDocker.logs Command with Volumes flag")
	}
}

func TestPublishFlagNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"--disable-tty", "--publish=publish_test", "image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	argsAppend := f.dockerRun.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 4 || argsAppend[2] != "--publish" || argsAppend[3] != "publish_test" {
		t.Errorf("bad arguments to KoolDocker.logs Command with Volumes flag")
	}
}

func TestImageCommandsNewDockerCommand(t *testing.T) {
	f := newFakeKoolDocker()
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"image", "command1", "command2"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	interactiveArgs := f.dockerRun.(*builder.FakeCommand).ArgsInteractive

	if len(interactiveArgs) != 3 || interactiveArgs[0] != "image" || interactiveArgs[1] != "command1" || interactiveArgs[2] != "command2" {
		t.Errorf("bad arguments to Interactive on KoolDocker.dockerRun Command")
	}
}

func TestFailingNewDockerCommand(t *testing.T) {
	f := newFailedFakeKoolDocker()
	cmd := NewDockerCommand(f)

	cmd.SetArgs([]string{"image"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing docker command; error: %v", err)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("expecting command to exit due to an error.")
	}

	if err := f.out.(*shell.FakeOutputWriter).Err; err.Error() != "error docker" {
		t.Errorf("expecting error 'error docker', got '%s'", err.Error())
	}
}
