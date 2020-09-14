package checker

import (
	"kool-dev/kool/cmd/builder"
	"log"
)

// Checker defines the check kool dependencies method
type Checker interface {
	CheckKoolDependencies() (string, error)
}

// DefaultChecker holds commands to be checked.
type DefaultChecker struct {
	DockerCmd builder.Executor
	dockerComposeCmd builder.Executor
}

// NewChecker initializes checker
func NewChecker() Checker {
	var (
		err error
		dockerInfoCmd *builder.Command
		dockerComposePsCmd *builder.Command
	)

	dockerInfoCmd, err = builder.ParseCommand("docker info")

	if err != nil {
		log.Fatal(err)
	}
	dockerComposePsCmd, err = builder.ParseCommand("docker-compose ps")

	if err != nil {
		log.Fatal(err)
	}

	return &DefaultChecker{dockerInfoCmd, dockerComposePsCmd}
}

// CheckKoolDependencies checks kool dependencies
func (c *DefaultChecker) CheckKoolDependencies() (message string, err error) {
	if err = c.DockerCmd.LookPath(); err != nil {
		message = "Docker doesn't seem to be installed, install it first and retry."
		return
	}

	if err = c.dockerComposeCmd.LookPath(); err != nil {
		message = "Docker-compose doesn't seem to be installed, install it first and retry."
		return
	}

	if _, err = c.DockerCmd.Exec(); err != nil {
		message = "Docker daemon doesn't seem to be running, run it first and retry."
		return
	}

	return
}
