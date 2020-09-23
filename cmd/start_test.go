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

type FakeStartDependenciesChecker struct {
	called bool
}

func (c *FakeStartDependenciesChecker) Check() (err error) {
	c.called = true
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

func TestStartAllCommand(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&FakeStartDependenciesChecker{},
		&FakeStartNetworkHandler{},
		&FakeStartRunner{},
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

	if len(startedServices) > 0 {
		t.Errorf("Expected no arguments, got '%v'", startedServices)
	}
}

func TestStartServicesCommand(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&FakeStartDependenciesChecker{},
		&FakeStartNetworkHandler{},
		&FakeStartRunner{},
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

	if !startedServicesAreEqual(startedServices, expected) {
		t.Errorf("Expect to start '%v', got '%v'", expected, startedServices)
	}
}

func TestFailedDependenciesStartCommand(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&FakeStartFailedDependenciesChecker{},
		&FakeStartNetworkHandler{},
		&FakeStartRunner{},
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
		&FakeStartDependenciesChecker{},
		&FakeStartFailedNetworkHandler{},
		&FakeStartRunner{},
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
		&FakeStartDependenciesChecker{},
		&FakeStartNetworkHandler{},
		&FakeFailedStartRunner{},
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
