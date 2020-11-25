package checker

import (
	"errors"
	"io/ioutil"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/shell"
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

	s := shell.NewShell()
	s.SetOutStream(ioutil.Discard)

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
	dockerComposeCmd := &builder.FakeCommand{MockLookPathError: errors.New("not installed")}

	s := shell.NewShell()
	s.SetOutStream(ioutil.Discard)

	c = &DefaultChecker{dockerCmd, dockerComposeCmd, s}

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

	dockerCmd := &builder.FakeCommand{MockError: errors.New("not running")}
	dockerComposeCmd := &builder.FakeCommand{}

	s := shell.NewShell()
	s.SetOutStream(ioutil.Discard)

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
	dockerComposeCmd.MockCmd = "docker-compose"

	s := &shell.FakeShell{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd, s}

	if err := c.Check(); err != nil {
		t.Errorf("Expected no errors, got %v.", err)
		return
	}

	if val, ok := c.(*DefaultChecker).shell.(*shell.FakeShell).CalledLookPath["docker"]; !val || !ok {
		t.Error("did not call LookPath for dockerCmd")
	}

	if val, ok := c.(*DefaultChecker).shell.(*shell.FakeShell).CalledLookPath["docker-compose"]; !val || !ok {
		t.Error("did not call LookPath for dockerComposeCmd")
	}

	if val, ok := c.(*DefaultChecker).shell.(*shell.FakeShell).CalledExec["docker"]; !val || !ok {
		t.Error("did not call Exec for dockerCmd")
	}
}
