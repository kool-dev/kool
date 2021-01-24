package compose

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"os"
	"strings"
	"testing"
)

func TestNewDockerCompose(t *testing.T) {
	dc := NewDockerCompose("cmd", "arg")

	if dc.isTTY {
		t.Error("unexpected default isTTY value for DockerCompose")
	}

	dc.SetIsTTY(true)

	if !dc.isTTY {
		t.Error("failed setting isTTY value for DockerCompose")
	}

	if _, ok := dc.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("unexpected default type for DockerCompose.env")
	}

	if !strings.HasSuffix(dc.String(), "cmd arg") {
		t.Errorf("unexpected DockerCompose.String() suffix: %s", dc.String())
	}

	dc.SetShell(&shell.FakeShell{})
	dc.SetLocalDockerCompose(&builder.FakeCommand{
		MockLookPathError: errors.New("some error"),
	})
	if !strings.HasPrefix(dc.String(), "docker run --rm -i") {
		t.Errorf("unexpected DockerCompose.String() prefix: %s", dc.String())
	}

	dc.SetIsTTY(false)
	if strings.Contains(dc.String(), " -t ") {
		t.Error("unexpected -t flag when not TTY")
	}
	dc.SetIsTTY(true)
	if !strings.Contains(dc.String(), " -t ") {
		t.Error("missing -t flag when on TTY")
	}

	dc.localDockerCompose = &builder.FakeCommand{
		MockLookPathError: nil,
	}
	if !strings.HasPrefix(dc.String(), "docker-compose") {
		t.Errorf("unexpected DockerCompose.String() prefix: %s", dc.String())
	}
}

func TestDockerComposeArgsParsing(t *testing.T) {
	dc := NewDockerCompose("cmd", "arg")
	dc.sh = &shell.FakeShell{}
	dc.localDockerCompose = &builder.FakeCommand{
		MockLookPathError: errors.New("some error"),
	}
	dc.env = environment.NewFakeEnvStorage()

	dc.env.Set("DOCKER_HOST", "")
	if !strings.Contains(dc.String(), "-e DOCKER_HOST") || !strings.Contains(dc.String(), "/var/run/docker.sock") {
		t.Error("failed parsing docker flags when DOCKER_HOST is not set")
	}
	dc.env.Set("DOCKER_HOST", "some-value")
	if !strings.Contains(dc.String(), "-e DOCKER_HOST") || !strings.Contains(dc.String(), "-e DOCKER_TLS_VERIFY") {
		t.Error("failed parsing docker flags when DOCKER_HOST is set to non-unix:// value")
	}

	cwd, _ := os.Getwd()
	if !strings.Contains(dc.String(), cwd) {
		t.Error("missing current working directory")
	}

	if strings.Contains(dc.String(), "-e HOME") {
		t.Error("unexpected passing down HOME empty variable")
	}
	dc.env.Set("HOME", "/my/home/path")
	if !strings.Contains(dc.String(), "-e HOME") || !strings.Contains(dc.String(), "/my/home/path") {
		t.Error("missing HOME variable and/or mount")
	}

	dc.env.Set("PATH", "some-path")
	dc.env.Set("SOME_VAR", "some-value")

	if strings.Contains(dc.String(), "-e PATH") {
		t.Error("unexpected passing down PATH variable")
	}
	if !strings.Contains(dc.String(), "-e SOME_VAR") {
		t.Error("missing passing down SOME_VAR variable")
	}

	if dc.Cmd() != "docker" {
		t.Error("unexpected Cmd() return")
	}
}
