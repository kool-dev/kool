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

type FakeStatusCmd struct{}

// CheckDependencies check kool dependencies
func (s *FakeStatusCmd) CheckDependencies() (err error) {
	return
}

// GetServices get docker services
func (s *FakeStatusCmd) GetServices() (services []string, err error) {
	services = []string{"app", "database", "cache"}
	return
}

// GetServiceID get docker service ID
func (s *FakeStatusCmd) GetServiceID(service string) (serviceID string, err error) {
	serviceID = "100"
	return
}

// GetStatusPort get docker service port and status
func (s *FakeStatusCmd) GetStatusPort(serviceID string) (status string, port string) {
	status = "Up About an hour"
	port = "0.0.0.0:80->80/tcp, 9000/tcp"
	return
}

type NotRunningStatusCmd struct {
	FakeStatusCmd
}

// GetStatusPort get docker service port and status
func (s *NotRunningStatusCmd) GetStatusPort(serviceID string) (status string, port string) {
	return
}

type NoServicesStatusCmd struct {
	FakeStatusCmd
}

// GetServices get docker services
func (s *NoServicesStatusCmd) GetServices() (services []string, err error) {
	return
}

type FailedDependenciesStatusCmd struct {
	FakeStatusCmd
}

// CheckDependencies check kool dependencies
func (s *FailedDependenciesStatusCmd) CheckDependencies() (err error) {
	err = errors.New("failed dependencies")
	return
}

type FailedGetServicesStatusCmd struct {
	FakeStatusCmd
}

// GetServices get docker services
func (s *FailedGetServicesStatusCmd) GetServices() (services []string, err error) {
	err = errors.New("failed get services")
	return
}

type FailedGetServiceIDStatusCmd struct {
	FakeStatusCmd
}

// GetServiceID get docker service ID
func (s *FailedGetServiceIDStatusCmd) GetServiceID(service string) (serviceID string, err error) {
	err = errors.New("failed get service id")
	return
}

func TestStatusCommand(t *testing.T) {
	cmd := NewStatusCommand(&FakeStatusCmd{})
	output, err := execCommand(cmd)

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
	cmd := NewStatusCommand(&NotRunningStatusCmd{})
	output, err := execCommand(cmd)

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
	output = strings.Trim(output, "\n")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNoServicesStatusCommand(t *testing.T) {
	cmd := NewStatusCommand(&NoServicesStatusCmd{})
	output, err := execCommand(cmd)

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
	cmd := NewStatusCommand(&FailedGetServicesStatusCmd{})
	output, err := execCommand(cmd)

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
		cmd := NewStatusCommand(&FailedDependenciesStatusCmd{})
		_, _ = execCommand(cmd)
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

func TestFailedGetServiceIDStatusCommand(t *testing.T) {
	if os.Getenv("FLAG") == "1" {
		cmd := NewStatusCommand(&FailedGetServiceIDStatusCmd{})
		_, _ = execCommand(cmd)
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

func execCommand(cmd *cobra.Command) (output string, err error) {
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
