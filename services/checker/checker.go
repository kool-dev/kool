package checker

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"strings"
)

// Checker defines the check kool dependencies method
type Checker interface {
	Check() error
}

// DefaultChecker holds commands to be checked.
type DefaultChecker struct {
	dockerCmd        builder.Command
	dockerComposeCmd builder.Command
	shell            shell.Shell
}

// NewChecker initializes checker
func NewChecker(s shell.Shell) *DefaultChecker {
	return &DefaultChecker{
		builder.NewCommand("docker", "info"),
		builder.NewCommand("docker", "compose", "ps"),
		s,
	}
}

// Check checks kool dependencies
func (c *DefaultChecker) Check() error {
	if err := c.shell.LookPath(c.dockerCmd); err != nil {
		return ErrDockerNotFound
	}

	if _, err := c.shell.Exec(c.dockerComposeCmd); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "is not a docker command") {
			return ErrDockerComposeNotFound
		}
		// anything else, raise the original error
		return err
	}

	if _, err := c.shell.Exec(c.dockerCmd); err != nil {
		return ErrDockerNotRunning
	}

	return nil
}
