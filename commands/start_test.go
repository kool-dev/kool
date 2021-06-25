package commands

import (
	"bytes"
	"errors"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/network"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/checker"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newFakeKoolStart() *KoolStart {
	return &KoolStart{
		*newFakeKoolService(),
		&KoolStartFlags{},
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
		&KoolRebuild{
			*newFakeKoolService(),
			&builder.FakeCommand{MockCmd: "pull"},
			&builder.FakeCommand{MockCmd: "build"},
		},
	}
}

func TestStartAllCommand(t *testing.T) {
	koolStart := newFakeKoolStart()

	cmd := NewStartCommand(koolStart)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if koolStart.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not expect for KoolStart service to call exit")
	}

	if koolStart.exiter.(*shell.FakeExiter).Code() != 0 {
		t.Errorf("did not expect KoolStart service to have exit code different than 0; got '%d", koolStart.exiter.(*shell.FakeExiter).Code())
	}

	interactiveArgs, ok := koolStart.shell.(*shell.FakeShell).ArgsInteractive["start"]

	if ok && len(interactiveArgs) > 0 {
		t.Errorf("Expected no arguments, got '%v'", interactiveArgs)
	}
}

func TestStartForegroundFlag(t *testing.T) {
	koolStart := newFakeKoolStart()

	if err := koolStart.Execute(nil); err != nil {
		t.Fatal(err)
	}

	args := koolStart.start.(*builder.FakeCommand).ArgsAppend
	if len(args) == 0 || args[0] != "-d" {
		t.Error("did not set -d on start")
	}

	koolStart = newFakeKoolStart()
	koolStart.Flags.Foreground = true

	if err := koolStart.Execute(nil); err != nil {
		t.Fatal(err)
	}

	args = koolStart.start.(*builder.FakeCommand).ArgsAppend
	if len(args) != 0 {
		t.Error("shoul not have appended args")
	}
}

func TestStartRebuildFlag(t *testing.T) {
	koolStart := newFakeKoolStart()

	if err := koolStart.Execute(nil); err != nil {
		t.Fatal(err)
	}

	rebuilder := koolStart.rebuilder.(*KoolRebuild)
	if rebuilder.pull.(*builder.FakeCommand).CalledCmd || rebuilder.build.(*builder.FakeCommand).CalledCmd {
		t.Error("should not have executed pull or build")
	}

	koolStart = newFakeKoolStart()
	rebuilder = koolStart.rebuilder.(*KoolRebuild)

	koolStart.Flags.Rebuild = true

	rebuilder.shell.(*shell.FakeShell).MockOutStream = io.Discard

	if err := koolStart.Execute(nil); err != nil {
		t.Fatal(err)
	}

	if !rebuilder.pull.(*builder.FakeCommand).CalledCmd || !rebuilder.build.(*builder.FakeCommand).CalledCmd {
		t.Error("should have executed pull and build")
	}
}

func TestStartServicesCommand(t *testing.T) {
	koolStart := newFakeKoolStart()

	cmd := NewStartCommand(koolStart)
	expected := []string{"app", "database"}
	cmd.SetArgs(expected)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if koolStart.exiter.(*shell.FakeExiter).Code() != 0 {
		t.Errorf("did not expect KoolStart to exit with error, got %d", koolStart.exiter.(*shell.FakeExiter).Code())
	}
	var startedServices []string
	if interactiveArgs, ok := koolStart.shell.(*shell.FakeShell).ArgsInteractive["start"]; ok {
		startedServices = interactiveArgs
	}

	if !startedServicesAreEqual(startedServices, expected) {
		t.Errorf("Expect to start '%v', got '%v'", expected, startedServices)
	}
}

func TestFailedDependenciesStartCommand(t *testing.T) {
	koolStart := newFakeKoolStart()
	koolStart.check.(*checker.FakeChecker).MockError = errors.New("dependencies")

	cmd := NewStartCommand(koolStart)

	_, err := execStartCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if koolStart.exiter.(*shell.FakeExiter).Code() != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", koolStart.exiter.(*shell.FakeExiter).Code())
	}
}

func TestFailedNetworkStartCommand(t *testing.T) {
	koolStart := newFakeKoolStart()
	koolStart.net.(*network.FakeHandler).MockError = errors.New("network")

	cmd := NewStartCommand(koolStart)

	_, err := execStartCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if koolStart.exiter.(*shell.FakeExiter).Code() != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", koolStart.exiter.(*shell.FakeExiter).Code())
	}
}

func TestStartWithError(t *testing.T) {
	koolStart := newFakeKoolStart()
	koolStart.start.(*builder.FakeCommand).MockInteractiveError = errors.New("start")

	cmd := NewStartCommand(koolStart)

	_, err := execStartCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if koolStart.exiter.(*shell.FakeExiter).Code() != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", koolStart.exiter.(*shell.FakeExiter).Code())
	}
}

func execStartCommand(cmd *cobra.Command) (output string, err error) {
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	if err = cmd.Execute(); err != nil {
		return
	}

	var out []byte
	if out, err = io.ReadAll(b); err != nil {
		return
	}

	output = strings.Trim(string(out), "\n")
	return
}

func startedServicesAreEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
