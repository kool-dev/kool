package checker

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/compose"
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
		compose.NewDockerCompose("ps"),
		s,
	}
}

// Check checks kool dependencies
func (c *DefaultChecker) Check() error {
	if err := c.shell.LookPath(c.dockerCmd); err != nil {
		return ErrDockerNotFound
	}

	if err := c.shell.LookPath(c.dockerComposeCmd); err != nil {
		return ErrDockerComposeNotFound
	}

	if _, err := c.shell.Exec(c.dockerCmd); err != nil {
		return ErrDockerNotRunning
	}

	return nil
}
