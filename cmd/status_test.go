package cmd

import (
	"bytes"
	"errors"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type FakeDependenciesChecker struct{}

func (c *FakeDependenciesChecker) VerifyDependencies() (err error) {
	return
}

type FakeFailedDependenciesChecker struct{}

func (c *FakeFailedDependenciesChecker) VerifyDependencies() (err error) {
	err = errors.New("")
	return
}

type FakeNetworkHandler struct{}

func (c *FakeNetworkHandler) HandleGlobalNetwork(networkName string) (err error) {
	return
}

type FakeFailedNetworkHandler struct{}

func (c *FakeFailedNetworkHandler) HandleGlobalNetwork(networkName string) (err error) {
	err = errors.New("")
	return
}

type FakeRunner struct{}

func (c *FakeRunner) LookPath() (err error) {
	return
}

func (c *FakeRunner) Interactive(args ...string) (err error) {
	return
}

func (c *FakeRunner) Exec(args ...string) (outStr string, err error) {
	return
}

type FakeGetServicesRunner struct {
	FakeRunner
}

func (c *FakeGetServicesRunner) Exec(args ...string) (outStr string, err error) {
	outStr = `
app
cache
database
`
	return
}

type FakeFailedGetServicesRunner struct {
	FakeRunner
}

func (c *FakeFailedGetServicesRunner) Exec(args ...string) (outStr string, err error) {
	err = errors.New("")
	return
}

type FakeGetServiceIDRunner struct {
	FakeRunner
}

func (c *FakeGetServiceIDRunner) Exec(args ...string) (outStr string, err error) {
	outStr = "100"
	return
}

type FakeFailedGetServiceIDRunner struct {
	FakeRunner
}

func (c *FakeFailedGetServiceIDRunner) Exec(args ...string) (outStr string, err error) {
	err = errors.New("")
	return
}

type FakeGetServiceStatusPortRunner struct {
	FakeRunner
}

func (c *FakeGetServiceStatusPortRunner) Exec(args ...string) (outStr string, err error) {
	outStr = "Up About an hour|0.0.0.0:80->80/tcp, 9000/tcp"
	return
}

type FakeNotRunningGetServiceStatusPortRunner struct {
	FakeRunner
}

func (c *FakeNotRunningGetServiceStatusPortRunner) Exec(args ...string) (outStr string, err error) {
	outStr = "Exited an hour ago"
	return
}

func TestStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeDependenciesChecker{},
		&FakeNetworkHandler{},
		&FakeGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	output, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := `
+----------+---------+------------------------------+------------------+
| SERVICE  | RUNNING | PORTS                        | STATE            |
+----------+---------+------------------------------+------------------+
| app      | Running | 0.0.0.0:80->80/tcp, 9000/tcp | Up About an hour |
| cache    | Running | 0.0.0.0:80->80/tcp, 9000/tcp | Up About an hour |
| database | Running | 0.0.0.0:80->80/tcp, 9000/tcp | Up About an hour |
+----------+---------+------------------------------+------------------+
`
	expected = strings.Trim(expected, "\n")
	output = strings.Trim(output, "\n")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNotRunningStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeDependenciesChecker{},
		&FakeNetworkHandler{},
		&FakeGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeNotRunningGetServiceStatusPortRunner{},
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	output, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := `
+----------+-------------+-------+--------------------+
| SERVICE  | RUNNING     | PORTS | STATE              |
+----------+-------------+-------+--------------------+
| app      | Not running |       | Exited an hour ago |
| cache    | Not running |       | Exited an hour ago |
| database | Not running |       | Exited an hour ago |
+----------+-------------+-------+--------------------+
`
	expected = strings.Trim(expected, "\n")
	output = strings.Trim(output, "\n")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNoServicesStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeDependenciesChecker{},
		&FakeNetworkHandler{},
		&FakeRunner{},
		&FakeGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	output, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := color.New(color.Yellow).Sprint("No services found.")
	output = strings.Trim(output, "\n")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFailedGetServicesStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeDependenciesChecker{},
		&FakeNetworkHandler{},
		&FakeFailedGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	output, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := color.New(color.Yellow).Sprint("No services found.")
	output = strings.Trim(output, "\n")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFailedDependenciesStatusCommand(t *testing.T) {
	if os.Getenv("FLAG") == "1" {
		defaultStatusCmd := &DefaultStatusCmd{
			&FakeFailedDependenciesChecker{},
			&FakeNetworkHandler{},
			&FakeGetServicesRunner{},
			&FakeGetServiceIDRunner{},
			&FakeGetServiceStatusPortRunner{},
		}
		cmd := NewStatusCommand(defaultStatusCmd)
		_, _ = execStatusCommand(cmd)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFailedDependenciesStatusCommand")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()

	e, ok := err.(*exec.ExitError)

	if !ok {
		t.Error("Expected an exit error")
	}

	if e == nil {
		t.Error("Expected an error, got none")
	}
}

func TestFailedNetworkStatusCommand(t *testing.T) {
	if os.Getenv("FLAG") == "1" {
		defaultStatusCmd := &DefaultStatusCmd{
			&FakeDependenciesChecker{},
			&FakeFailedNetworkHandler{},
			&FakeGetServicesRunner{},
			&FakeGetServiceIDRunner{},
			&FakeGetServiceStatusPortRunner{},
		}
		cmd := NewStatusCommand(defaultStatusCmd)
		_, _ = execStatusCommand(cmd)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFailedNetworkStatusCommand")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()

	e, ok := err.(*exec.ExitError)

	if !ok {
		t.Error("Expected an exit error")
	}

	if e == nil {
		t.Error("Expected an error, got none")
	}
}

func TestFailedGetServiceIDStatusCommand(t *testing.T) {
	if os.Getenv("FLAG") == "1" {
		defaultStatusCmd := &DefaultStatusCmd{
			&FakeDependenciesChecker{},
			&FakeNetworkHandler{},
			&FakeGetServicesRunner{},
			&FakeFailedGetServiceIDRunner{},
			&FakeGetServiceStatusPortRunner{},
		}
		cmd := NewStatusCommand(defaultStatusCmd)
		_, _ = execStatusCommand(cmd)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFailedGetServiceIDStatusCommand")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()

	e, ok := err.(*exec.ExitError)

	if !ok {
		t.Error("Expected an exit error")
	}

	if e == nil {
		t.Error("Expected an error, got none")
	}
}

func execStatusCommand(cmd *cobra.Command) (output string, err error) {
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	if err = cmd.Execute(); err != nil {
		return
	}

	var out []byte
	if out, err = ioutil.ReadAll(b); err != nil {
		return
	}

	output = string(out)
	return
}
