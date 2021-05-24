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

func TestStartAllCommand(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
	}

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

func TestStartServicesCommand(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
	}

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
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{MockError: errors.New("dependencies")},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
	}

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
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{MockError: errors.New("network")},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
	}

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
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start", MockInteractiveError: errors.New("start")},
	}

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
