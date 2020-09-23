package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"kool-dev/kool/cmd/shell"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

type FakeStartDependenciesChecker struct{}

func (c *FakeStartDependenciesChecker) Check() (err error) {
	return
}

type FakeStartFailedDependenciesChecker struct{}

func (c *FakeStartFailedDependenciesChecker) Check() (err error) {
	err = errors.New("dependencies")
	return
}

type FakeStartNetworkHandler struct{}

func (c *FakeStartNetworkHandler) HandleGlobalNetwork(networkName string) (err error) {
	return
}

type FakeStartFailedNetworkHandler struct{}

func (c *FakeStartFailedNetworkHandler) HandleGlobalNetwork(networkName string) (err error) {
	err = errors.New("network")
	return
}

type FakeStartRunner struct{}

var startedServices []string

func (c *FakeStartRunner) LookPath() (err error) {
	return
}

func (c *FakeStartRunner) Interactive(args ...string) (err error) {
	startedServices = []string{}
	if len(args) > 0 {
		startedServices = args
	}
	return
}

func (c *FakeStartRunner) Exec(args ...string) (outStr string, err error) {
	return
}

type FakeFailedStartRunner struct {
	FakeStartRunner
}

func (c *FakeFailedStartRunner) Interactive(args ...string) (err error) {
	err = errors.New("")
	return
}

var startExitCode int

type FakeStartExiter struct{}

func (e *FakeStartExiter) Exit(code int) {
	startExitCode = code
}

func TestStartAllCommand(t *testing.T) {
	defaultStartCmd := &KoolStart{
		&FakeStartDependenciesChecker{},
		&FakeStartNetworkHandler{},
		&FakeStartRunner{},
		&FakeStartExiter{},
		&shell.FakeOutputWriter{},
	}

	cmd := NewStartCommand(defaultStartCmd)
	startExitCode = 0

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if startExitCode != 0 {
		t.Errorf("Expected no exit error, got '%v'", startExitCode)
	}

	if len(startedServices) > 0 {
		t.Errorf("Expected no arguments, got '%v'", startedServices)
	}
}

func TestStartServicesCommand(t *testing.T) {
	defaultStartCmd := &KoolStart{
		&FakeStartDependenciesChecker{},
		&FakeStartNetworkHandler{},
		&FakeStartRunner{},
		&FakeStartExiter{},
		&shell.FakeOutputWriter{},
	}

	cmd := NewStartCommand(defaultStartCmd)
	startExitCode = 0
	expected := []string{"app", "database"}
	cmd.SetArgs(expected)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if startExitCode != 0 {
		t.Errorf("Expected no exit error, got '%v'", startExitCode)
	}

	if !startedServicesAreEqual(startedServices, expected) {
		t.Errorf("Expect to start '%v', got '%v'", expected, startedServices)
	}
}

func TestFailedDependenciesStartCommand(t *testing.T) {
	defaultStartCmd := &KoolStart{
		&FakeStartFailedDependenciesChecker{},
		&FakeStartNetworkHandler{},
		&FakeStartRunner{},
		&FakeStartExiter{},
		&shell.FakeOutputWriter{},
	}

	cmd := NewStartCommand(defaultStartCmd)
	startExitCode = 0

	_, err := execStartCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if startExitCode != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", statusExitCode)
	}
}

func TestFailedNetworkStartCommand(t *testing.T) {
	defaultStartCmd := &KoolStart{
		&FakeStartDependenciesChecker{},
		&FakeStartFailedNetworkHandler{},
		&FakeStartRunner{},
		&FakeStartExiter{},
		&shell.FakeOutputWriter{},
	}

	cmd := NewStartCommand(defaultStartCmd)
	startExitCode = 0

	_, err := execStartCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if startExitCode != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", statusExitCode)
	}
}

func TestStartWithError(t *testing.T) {
	defaultStartCmd := &KoolStart{
		&FakeStartDependenciesChecker{},
		&FakeStartNetworkHandler{},
		&FakeFailedStartRunner{},
		&FakeStartExiter{},
		&shell.FakeOutputWriter{},
	}

	cmd := NewStartCommand(defaultStartCmd)
	startExitCode = 0

	_, err := execStartCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if startExitCode != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", statusExitCode)
	}
}

func execStartCommand(cmd *cobra.Command) (output string, err error) {
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	if err = cmd.Execute(); err != nil {
		return
	}

	var out []byte
	if out, err = ioutil.ReadAll(b); err != nil {
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
