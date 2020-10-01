package checker

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"testing"
)

func TestDefaultChecker(t *testing.T) {
	var c Checker = NewChecker()

	if _, assert := c.(*DefaultChecker); !assert {
		t.Errorf("NewChecker() did not return a *DefaultChecker")
	}
}

func TestDockerNotInstalled(t *testing.T) {
	var c Checker

	dockerCmd := &builder.FakeCommand{MockLookPathError: errors.New("not installed")}
	dockerComposeCmd := &builder.FakeCommand{}

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

	dockerCmd := &builder.FakeCommand{}
	dockerComposeCmd := &builder.FakeCommand{MockLookPathError: errors.New("not installed")}

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

	dockerCmd := &builder.FakeCommand{MockError: errors.New("not running")}
	dockerComposeCmd := &builder.FakeCommand{}

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

	dockerCmd := &builder.FakeCommand{}
	dockerComposeCmd := &builder.FakeCommand{}

	c = &DefaultChecker{dockerCmd, dockerComposeCmd}

	if err := c.Check(); err != nil {
		t.Errorf("Expected no errors, got %v.", err)
		return
	}

	if !c.(*DefaultChecker).dockerCmd.(*builder.FakeCommand).CalledLookPath {
		t.Error("did not call LookPath for dockerCmd")
	}

	if !c.(*DefaultChecker).dockerComposeCmd.(*builder.FakeCommand).CalledLookPath {
		t.Error("did not call LookPath for dockerComposeCmd")
	}

	if !c.(*DefaultChecker).dockerCmd.(*builder.FakeCommand).CalledExec {
		t.Error("did not call Exec for dockerCmd")
	}
}
