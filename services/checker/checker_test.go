package checker

import (
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"testing"
)

func TestDefaultChecker(t *testing.T) {
	var c Checker = NewChecker(&shell.FakeShell{})

	if _, assert := c.(*DefaultChecker); !assert {
		t.Errorf("NewChecker() did not return a *DefaultChecker")
	}
}

func TestDockerNotInstalled(t *testing.T) {
	var c Checker

	dockerCmd := &builder.FakeCommand{MockLookPathError: errors.New("not installed")}
	dockerComposeCmd := &builder.FakeCommand{}

	s := &shell.FakeShell{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd, s}

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

	dockerCmd := &builder.FakeCommand{}
	dockerComposeCmd := &builder.FakeCommand{MockExecError: errors.New("is not a docker command")}

	s := &shell.FakeShell{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd, s}

	err := c.Check()

	if err == nil {
		t.Error("Expected an error, got none.")
		return
	}

	if !IsDockerComposeNotFoundError(err) {
		t.Errorf("Expected the message '%s', got '%s'", ErrDockerComposeNotFound.Error(), err.Error())
	}

	dockerComposeCmd.MockExecError = errors.New("some other error")

	if err := c.Check(); err == nil || err.Error() != "some other error" {
		t.Errorf("Expected the error message 'some other error', got '%v'", err)
	}
}

func TestDockerNotRunning(t *testing.T) {
	var c Checker

	dockerCmd := &builder.FakeCommand{MockExecError: errors.New("not running")}
	dockerComposeCmd := &builder.FakeCommand{}

	s := &shell.FakeShell{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd, s}

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

	dockerCmd := &builder.FakeCommand{}
	dockerCmd.MockCmd = "docker"

	dockerComposeCmd := &builder.FakeCommand{}
	dockerComposeCmd.MockCmd = "command"

	s := &shell.FakeShell{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd, s}

	if err := c.Check(); err != nil {
		t.Errorf("Expected no errors, got %v.", err)
		return
	}

	if val, ok := c.(*DefaultChecker).shell.(*shell.FakeShell).CalledLookPath["docker"]; !val || !ok {
		t.Error("did not call LookPath for dockerCmd")
	}

	if val, ok := c.(*DefaultChecker).shell.(*shell.FakeShell).CalledExec["command"]; !val || !ok {
		t.Error("did not call LookPath for dockerComposeCmd")
	}

	if val, ok := c.(*DefaultChecker).shell.(*shell.FakeShell).CalledExec["docker"]; !val || !ok {
		t.Error("did not call Exec for dockerCmd")
	}
}
