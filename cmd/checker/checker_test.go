package checker

import (
	"kool-dev/kool/cmd/builder"
	"errors"
	"testing"
)

type FakeCommand struct {}

func (c *FakeCommand) AppendArgs(args ...string) {}

func (c *FakeCommand) LookPath() (err error) {
	return
}

func (c *FakeCommand) Interactive() (err error) {
	return
}

func (c *FakeCommand) Exec() (outStr string, err error) {
	return
}

func (c *FakeCommand) String() (strCommand string) {
	return
}

type NotInstalledDockerCmd struct {
	FakeCommand
}

func (c *NotInstalledDockerCmd) LookPath() (err error) {
	err = errors.New("not installed")
	return
}

type NotRunningDockerCmd struct {
	FakeCommand
}

func (c *NotRunningDockerCmd) Exec() (outStr string, err error) {
	err = errors.New("not running")
	outStr = "error"
	return
}

type NotInstalledDockerComposeCmd struct {
	FakeCommand
}

func (c *NotInstalledDockerComposeCmd) LookPath() (err error) {
	err = errors.New("not installed")
	return
}

func TestDefaultChecker(t *testing.T) {
	var c Checker = NewChecker()

	if _, assert := c.(*DefaultChecker); !assert {
		t.Errorf("NewChecker() did not return a *DefaultChecker")
	}
}

func TestDockerNotInstalled(t *testing.T) {
	var c Checker

	dockerCmd := &NotInstalledDockerCmd{}
	dockerComposeCmd, _ := builder.ParseCommand("docker-compose ps")

	c = &DefaultChecker{dockerCmd, dockerComposeCmd}

	message, err := c.CheckKoolDependencies()

	if err == nil {
		t.Error("Expected an error, got none.")
		return
	}

	if message != "Docker doesn't seem to be installed, install it first and retry." {
		t.Errorf("Expected the message 'Docker doesn't seem to be installed, install it first and retry.', got %s", message)
	}
}

func TestDockerComposeNotInstalled(t *testing.T) {
	var c Checker

	dockerCmd := &FakeCommand{}
	dockerComposeCmd := &NotInstalledDockerComposeCmd{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd}

	message, err := c.CheckKoolDependencies()

	if err == nil {
		t.Error("Expected an error, got none.")
		return
	}

	if message != "Docker-compose doesn't seem to be installed, install it first and retry." {
		t.Errorf("Expected the message 'Docker-compose doesn't seem to be installed, install it first and retry.', got %s", message)
	}
}

func TestDockerNotRunning(t *testing.T) {
	var c Checker

	dockerCmd := &NotRunningDockerCmd{}
	dockerComposeCmd := &FakeCommand{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd}

	message, err := c.CheckKoolDependencies()

	if err == nil {
		t.Error("Expected an error, got none.")
		return
	}

	if message != "Docker daemon doesn't seem to be running, run it first and retry." {
		t.Errorf("Expected the message 'Docker daemon doesn't seem to be running, run it first and retry.', got %s", message)
	}
}

func TestCheckKoolDependencies(t *testing.T) {
	var c Checker

	dockerCmd := &FakeCommand{}
	dockerComposeCmd := &FakeCommand{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd}

	message, err := c.CheckKoolDependencies()

	if err != nil {
		t.Errorf("Expected no errors, got %v.", err)
		return
	}

	if message != "" {
		t.Errorf("Expected no message, got %s", message)
	}
}
