package checker

import "kool-dev/kool/cmd/builder"

// Checker defines the check kool dependencies method
type Checker interface {
	Check() error
}

// DefaultChecker holds commands to be checked.
type DefaultChecker struct {
	dockerCmd        builder.Runner
	dockerComposeCmd builder.Runner
}

// NewChecker initializes checker
func NewChecker() *DefaultChecker {
	var dockerInfoCmd, dockerComposePsCmd *builder.DefaultCommand

	dockerInfoCmd = builder.NewCommand("docker", "info")
	dockerComposePsCmd = builder.NewCommand("docker-compose", "ps")

	return &DefaultChecker{dockerInfoCmd, dockerComposePsCmd}
}

// Check checks kool dependencies
func (c *DefaultChecker) Check() error {
	if err := c.dockerCmd.LookPath(); err != nil {
		return ErrDockerNotFound
	}

	if err := c.dockerComposeCmd.LookPath(); err != nil {
		return ErrDockerComposeNotFound
	}

	if _, err := c.dockerCmd.Exec(); err != nil {
		return ErrDockerNotRunning
	}

	return nil
}
