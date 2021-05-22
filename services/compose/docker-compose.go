package compose

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"os"
	"strings"
)

// DockerComposeImage holds the Docker image:tag to use for Docker Compose
const DockerComposeImage = "docker/compose:1.28.0"

// TtyAware interface holds functions for becoming aware of TTY
type TtyAware interface {
	SetIsTTY(bool)
}

// DockerCompose holds data and logic to wrap docker-compose command
// within a container for flexibility
type DockerCompose struct {
	builder.Command
	localDockerCompose builder.Command
	env                environment.EnvStorage
	sh                 shell.Shell
	isTTY              bool
}

// NewDockerCompose creates a new instance of DockerCompose
func NewDockerCompose(cmd string, args ...string) *DockerCompose {
	return &DockerCompose{
		Command:            builder.NewCommand(cmd, args...),
		env:                environment.NewEnvStorage(),
		sh:                 shell.NewShell(),
		localDockerCompose: builder.NewCommand("docker-compose"),
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

// SetLocalDockerCompose sets the builder.Command to be used for checking
// docker-compose on PATH
func (c *DockerCompose) SetLocalDockerCompose(cmd builder.Command) *DockerCompose {
	c.localDockerCompose = cmd

	return c
}

// Args returns the command arguments
func (c *DockerCompose) Args() (args []string) {
	if c.sh.LookPath(c.localDockerCompose) == nil {
		return append([]string{c.Command.Cmd()}, c.Command.Args()...)
	}

	args = append(args, "run", "--rm", "-i")

	if c.isTTY {
		args = append(args, "-t")
	}

	dockerHost := c.env.Get("DOCKER_HOST")

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
		c.env.Set("DOCKER_HOST", dockerHost)
	}

	if strings.HasPrefix(dockerHost, "unix://") {
		path := strings.TrimPrefix(dockerHost, "unix://")
		args = append(args, "-v", fmt.Sprintf("%s:%s", path, path), "-e", "DOCKER_HOST")
	} else {
		args = append(args, "-e", "DOCKER_HOST", "-e", "DOCKER_TLS_VERIFY", "-e", "DOCKER_CERT_PATH")
	}

	cwd, _ := os.Getwd()
	if cwd != "/" {
		args = append(args, "-v", fmt.Sprintf("%s:%s", cwd, cwd))
	}
	args = append(args, "-w", cwd)
	if home := c.env.Get("HOME"); home != "" {
		args = append(args, "-v", fmt.Sprintf("%s:%s", home, home), "-e HOME")
	}

	for _, env := range c.env.All() {
		key := strings.SplitN(env, "=", 2)[0]

		if key == "PATH" {
			continue
		}

		args = append(args, "-e", key)
	}

	args = append(args, DockerComposeImage, "-p", c.env.Get("KOOL_NAME"), c.Command.Cmd())
	return append(args, c.Command.Args()...)
}

// Cmd returns the command executable
func (c *DockerCompose) Cmd() string {
	if c.sh.LookPath(c.localDockerCompose) == nil {
		return "docker-compose"
	}

	return "docker"
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
