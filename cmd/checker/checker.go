package checker

import "kool-dev/kool/cmd/builder"

// Checker defines the check kool dependencies method
type Checker interface {
	VerifyDependencies() error
}

// DefaultChecker holds commands to be checked.
type DefaultChecker struct {
	DockerCmd        builder.Command
	dockerComposeCmd builder.Command
}

// NewChecker initializes checker
func NewChecker() *DefaultChecker {
	var dockerInfoCmd, dockerComposePsCmd *builder.DefaultCommand

	dockerInfoCmd = builder.NewCommand("docker", []string{"info"})
	dockerComposePsCmd = builder.NewCommand("docker-compose", []string{"ps"})

	return &DefaultChecker{dockerInfoCmd, dockerComposePsCmd}
}

// VerifyDependencies checks kool dependencies
func (c *DefaultChecker) VerifyDependencies() error {
	if err := c.DockerCmd.LookPath(); err != nil {
		return ErrDockerNotFound
	}

	if err := c.dockerComposeCmd.LookPath(); err != nil {
		return ErrDockerComposeNotFound
	}

	if _, err := c.DockerCmd.Exec(); err != nil {
		return ErrDockerNotRunning
	}

	return nil
}
