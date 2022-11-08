package compose

import (
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
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

	dc.SetShell(&shell.FakeShell{})

	dc.SetIsTTY(false)
	if strings.Contains(dc.String(), " -t ") {
		t.Error("unexpected -t flag when not TTY")
	}
	dc.SetIsTTY(true)
	if !strings.Contains(dc.String(), " -t ") {
		t.Error("missing -t flag when on TTY")
	}

	if !strings.HasPrefix(dc.String(), "docker compose") {
		t.Errorf("unexpected DockerCompose.String() prefix: %s", dc.String())
	}
}
