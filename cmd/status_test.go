package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"kool-dev/kool/cmd/shell"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type FakeStatusDependenciesChecker struct{}

func (c *FakeStatusDependenciesChecker) Check() (err error) {
	return
}

type FakeStatusFailedDependenciesChecker struct{}

func (c *FakeStatusFailedDependenciesChecker) Check() (err error) {
	err = errors.New("")
	return
}

type FakeStatusNetworkHandler struct{}

func (c *FakeStatusNetworkHandler) HandleGlobalNetwork(networkName string) (err error) {
	return
}

type FakeStatusFailedNetworkHandler struct{}

func (c *FakeStatusFailedNetworkHandler) HandleGlobalNetwork(networkName string) (err error) {
	err = errors.New("")
	return
}

type FakeStatusRunner struct{}

func (c *FakeStatusRunner) LookPath() (err error) {
	return
}

func (c *FakeStatusRunner) Interactive(args ...string) (err error) {
	return
}

func (c *FakeStatusRunner) Exec(args ...string) (outStr string, err error) {
	return
}

type FakeGetServicesRunner struct {
	FakeStatusRunner
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
	FakeStatusRunner
}

func (c *FakeFailedGetServicesRunner) Exec(args ...string) (outStr string, err error) {
	err = errors.New("")
	return
}

type FakeGetServiceIDRunner struct {
	FakeStatusRunner
}

func (c *FakeGetServiceIDRunner) Exec(args ...string) (outStr string, err error) {
	outStr = "100"
	return
}

type FakeFailedGetServiceIDRunner struct {
	FakeStatusRunner
}

func (c *FakeFailedGetServiceIDRunner) Exec(args ...string) (outStr string, err error) {
	err = errors.New("")
	return
}

type FakeGetServiceStatusPortRunner struct {
	FakeStatusRunner
}

func (c *FakeGetServiceStatusPortRunner) Exec(args ...string) (outStr string, err error) {
	outStr = "Up About an hour|0.0.0.0:80->80/tcp, 9000/tcp"
	return
}

type FakeNotRunningGetServiceStatusPortRunner struct {
	FakeStatusRunner
}

func (c *FakeNotRunningGetServiceStatusPortRunner) Exec(args ...string) (outStr string, err error) {
	outStr = "Exited an hour ago"
	return
}

var statusExitCode int

type FakeStatusExiter struct{}

func (e *FakeStatusExiter) Exit(code int) {
	statusExitCode = code
}

func TestStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeStatusDependenciesChecker{},
		&FakeStatusNetworkHandler{},
		&FakeGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
		&FakeStatusExiter{},
		&shell.FakeOutputWriter{},
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

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNotRunningStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeStatusDependenciesChecker{},
		&FakeStatusNetworkHandler{},
		&FakeGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeNotRunningGetServiceStatusPortRunner{},
		&FakeStatusExiter{},
		&shell.FakeOutputWriter{},
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

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNoStatusPortStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeStatusDependenciesChecker{},
		&FakeStatusNetworkHandler{},
		&FakeGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeStatusRunner{},
		&FakeStatusExiter{},
		&shell.FakeOutputWriter{},
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	output, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := `
+----------+-------------+-------+-------+
| SERVICE  | RUNNING     | PORTS | STATE |
+----------+-------------+-------+-------+
| app      | Not running |       |       |
| cache    | Not running |       |       |
| database | Not running |       |       |
+----------+-------------+-------+-------+
`
	expected = strings.Trim(expected, "\n")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNoServicesStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeStatusDependenciesChecker{},
		&FakeStatusNetworkHandler{},
		&FakeStatusRunner{},
		&FakeGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
		&FakeStatusExiter{},
		shell.NewOutputWriter(),
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	output, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := color.New(color.Yellow).Sprint("No services found.")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFailedGetServicesStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeStatusDependenciesChecker{},
		&FakeStatusNetworkHandler{},
		&FakeFailedGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
		&FakeStatusExiter{},
		shell.NewOutputWriter(),
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	output, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := color.New(color.Yellow).Sprint("No services found.")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFailedDependenciesStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeStatusFailedDependenciesChecker{},
		&FakeStatusNetworkHandler{},
		&FakeGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
		&FakeStatusExiter{},
		&shell.FakeOutputWriter{},
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	statusExitCode = 0

	_, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if statusExitCode != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", statusExitCode)
	}
}

func TestFailedNetworkStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeStatusDependenciesChecker{},
		&FakeStatusFailedNetworkHandler{},
		&FakeGetServicesRunner{},
		&FakeGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
		&FakeStatusExiter{},
		&shell.FakeOutputWriter{},
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	statusExitCode = 0

	_, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if statusExitCode != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", statusExitCode)
	}
}

func TestFailedGetServiceIDStatusCommand(t *testing.T) {
	defaultStatusCmd := &DefaultStatusCmd{
		&FakeStatusDependenciesChecker{},
		&FakeStatusNetworkHandler{},
		&FakeGetServicesRunner{},
		&FakeFailedGetServiceIDRunner{},
		&FakeGetServiceStatusPortRunner{},
		&FakeStatusExiter{},
		&shell.FakeOutputWriter{},
	}
	cmd := NewStatusCommand(defaultStatusCmd)
	statusExitCode = 0

	_, err := execStatusCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	if statusExitCode != 1 {
		t.Errorf("Expected an exit code 1, got '%v'", statusExitCode)
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

	output = strings.Trim(string(out), "\n")
	return
}
