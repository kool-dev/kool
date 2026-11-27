package cloud

import (
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"reflect"
	"testing"
)

func TestDockerLoginError(t *testing.T) {
	exitErr := errors.New("exit status 1")

	err := dockerLoginError("registry.example.com", "Error response from daemon: unauthorized\n", exitErr)
	if err.Error() != "docker login to registry.example.com failed: Error response from daemon: unauthorized (exit status 1)" {
		t.Errorf("unexpected error message: %v", err)
	}
	if !errors.Is(err, exitErr) {
		t.Error("expected the docker login error to wrap the original error")
	}

	// the shell already merges the output into the error; do not repeat it
	shellErr := errors.New("Error response from daemon: unauthorized (exit status 1)")
	err = dockerLoginError("registry.example.com", "Error response from daemon: unauthorized", shellErr)
	if err.Error() != "docker login to registry.example.com failed: Error response from daemon: unauthorized (exit status 1)" {
		t.Errorf("unexpected error message: %v", err)
	}

	err = dockerLoginError("registry.example.com", "", exitErr)
	if err.Error() != "docker login to registry.example.com failed: exit status 1" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDockerLogoutCommand(t *testing.T) {
	logout := dockerLogoutCommand("123456789012.dkr.ecr.us-east-1.amazonaws.com")

	// a bare docker logout would log the user out of Docker Hub
	if logout.Cmd() != "docker" || !reflect.DeepEqual(logout.Args(), []string{"logout", "123456789012.dkr.ecr.us-east-1.amazonaws.com"}) {
		t.Errorf("expected docker logout of the Kool.dev Cloud registry only; got %s %v", logout.Cmd(), logout.Args())
	}
}

func TestDockerLogout(t *testing.T) {
	sh := &shell.FakeShell{}
	dockerLogout(sh, &builder.FakeCommand{MockCmd: "docker"})

	if !sh.CalledExec["docker"] {
		t.Error("expected docker logout to run")
	}
	if sh.CalledWarning {
		t.Error("unexpected warning for a successful logout")
	}

	// the image is pushed by then; a failed logout must not fail the deploy
	sh = &shell.FakeShell{}
	dockerLogout(sh, &builder.FakeCommand{MockCmd: "docker", MockExecError: errors.New("logout failed")})

	if !sh.CalledWarning {
		t.Error("expected a warning for a failed logout")
	}
}
