package checker

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"testing"
)

type FakeCommand struct{}

func (c *FakeCommand) LookPath() (err error) {
	return
}

func (c *FakeCommand) Interactive(args ...string) (err error) {
	return
}

func (c *FakeCommand) Exec(args ...string) (outStr string, err error) {
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

func (c *NotRunningDockerCmd) Exec(args ...string) (outStr string, err error) {
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

	err := c.Check()

	if err == nil {
		t.Error("Expected an error, got none.")
		return
	}

	if !IsDockerNotFoundError(err) {
		t.Errorf("Expected the message '%s', got '%s'", ErrDockerNotFound.Error(), err.Error())
	}
}

func TestDockerComposeNotInstalled(t *testing.T) {
	var c Checker

	dockerCmd := &FakeCommand{}
	dockerComposeCmd := &NotInstalledDockerComposeCmd{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd}

	err := c.Check()

	if err == nil {
		t.Error("Expected an error, got none.")
		return
	}

	if !IsDockerComposeNotFoundError(err) {
		t.Errorf("Expected the message '%s', got '%s'", ErrDockerComposeNotFound.Error(), err.Error())
	}
}

func TestDockerNotRunning(t *testing.T) {
	var c Checker

	dockerCmd := &NotRunningDockerCmd{}
	dockerComposeCmd := &FakeCommand{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd}

	err := c.Check()

	if err == nil {
		t.Error("Expected an error, got none.")
		return
	}

	if !IsDockerNotRunningError(err) {
		t.Errorf("Expected the message '%s', got '%s'", ErrDockerNotRunning.Error(), err.Error())
	}
}

func TestCheckKoolDependencies(t *testing.T) {
	var c Checker

	dockerCmd := &FakeCommand{}
	dockerComposeCmd := &FakeCommand{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd}

	if err := c.Check(); err != nil {
		t.Errorf("Expected no errors, got %v.", err)
		return
	}
}
