package compose

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"strings"
)

// TtyAware interface holds functions for becoming aware of TTY
type TtyAware interface {
	SetIsTTY(bool)
}

// DockerCompose holds data and logic to wrap docker-compose command
// within a container for flexibility
type DockerCompose struct {
	builder.Command
	env   environment.EnvStorage
	sh    shell.Shell
	isTTY bool
}

// NewDockerCompose creates a new instance of DockerCompose
func NewDockerCompose(cmd string, args ...string) *DockerCompose {
	command := builder.NewCommand("docker compose")
	command.AppendArgs(cmd)
	command.AppendArgs(args...)

	return &DockerCompose{
		Command: command,
		env:     environment.NewEnvStorage(),
		sh:      shell.NewShell(),
	}
}

// SetIsTTY sets whether we are under TTY
func (c *DockerCompose) SetIsTTY(tty bool) {
	c.isTTY = tty
}

// SetShell sets the shell.Shell to be used
func (c *DockerCompose) SetShell(sh shell.Shell) *DockerCompose {
	c.sh = sh

	return c
}

// Args returns the command arguments
func (c *DockerCompose) Args() (args []string) {
	args = c.Command.Args()

	if c.isTTY {
		args = append([]string{"-t"}, args...)
	}

	return
}

// Cmd returns the command executable
func (c *DockerCompose) Cmd() string {
	return "docker compose"
}

func (c *DockerCompose) String() string {
	return strings.Trim(fmt.Sprintf("%s %s", c.Cmd(), strings.Join(c.Args(), " ")), " ")
}

// Copy clones the pointer to avoid unintended modifications
func (c *DockerCompose) Copy() (copied builder.Command) {
	copied = NewDockerCompose(c.Command.Cmd(), c.Command.Args()...)
	copied.(*DockerCompose).SetIsTTY(c.isTTY)
	return
}
