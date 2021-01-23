package compose

import (
	"fmt"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"os"
	"strings"
)

const DOCKER_COMPOSE_IMAGE = "docker/compose:1.28.0"

type DockerCompose struct {
	builder.Command
	env   environment.EnvStorage
	isTTY bool
}

func NewDockerCompose(cmd string, args ...string) *DockerCompose {
	return &DockerCompose{
		Command: builder.NewCommand(cmd, args...),
		env:     environment.NewEnvStorage(),
	}
}

// SetIsTTY sets whether we are under TTY
func (c *DockerCompose) SetIsTTY(tty bool) *DockerCompose {
	c.isTTY = tty

	return c
}

// Args returns the command arguments
func (c *DockerCompose) Args() (args []string) {
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

	args = append(args, DOCKER_COMPOSE_IMAGE, "-p", c.env.Get("KOOL_NAME"), c.Command.Cmd())
	return append(args, c.Command.Args()...)
}

// Cmd returns the command executable
func (c *DockerCompose) Cmd() string {
	return "docker"
}

func (c *DockerCompose) String() string {
	return strings.Trim(fmt.Sprintf("%s %s", c.Cmd(), strings.Join(c.Args(), " ")), " ")
}
