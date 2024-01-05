package k8s

import (
	"errors"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/api"
	"os"
	"strings"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	k := NewDefaultK8S()
	k.deployExec.Endpoint.(*api.DefaultEndpoint).Fake()

	expectedErr := errors.New("call error")
	k.deployExec.Endpoint.(*api.DefaultEndpoint).MockErr(expectedErr)

	if _, err := k.Authenticate("foo", "bar"); !errors.Is(err, expectedErr) {
		t.Error("unexpected error return from Authenticate")
	}

	k.deployExec.Endpoint.(*api.DefaultEndpoint).MockErr(nil)
	k.deployExec.Endpoint.(*api.DefaultEndpoint).MockResp(&api.DeployExecResponse{
		Server:    "server",
		Namespace: "ns",
		Path:      "path",
		Token:     "",
		CA:        "ca",
	})

	if _, err := k.Authenticate("foo", "bar"); !strings.Contains(err.Error(), "failed to generate access credentials") {
		t.Errorf("unexpected error from DeployExec call: %v", err)
	}

	k.deployExec.Endpoint.(*api.DefaultEndpoint).MockResp(&api.DeployExecResponse{
		Server:    "server",
		Namespace: "ns",
		Path:      "path",
		Token:     "token",
		CA:        "ca",
	})

	authTempPath = t.TempDir()

	if cloudService, err := k.Authenticate("foo", "bar"); err != nil {
		t.Errorf("unexpected error from Authenticate call: %v", err)
	} else if cloudService != "path" {
		t.Errorf("unexpected cloudService return: %s", cloudService)
	}
}

func TestTempCAPath(t *testing.T) {
	k := NewDefaultK8S()

	authTempPath = "fake-path"

	if !strings.Contains(k.getTempCAPath(), authTempPath) {
		t.Error("missing authTempPath from temp CA path")
	}
}

func TestCleanup(t *testing.T) {
	k := NewDefaultK8S()

	authTempPath = t.TempDir()
	if err := os.WriteFile(k.getTempCAPath(), []byte("ca"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fakeShell := &shell.FakeShell{}
	k.Cleanup(fakeShell)

	if fakeShell.CalledWarning {
		t.Error("should not have warned on removing the file")
	}

	authTempPath = t.TempDir() + "test"
	fakeShell = &shell.FakeShell{}
	k.Cleanup(fakeShell)

	if !fakeShell.CalledWarning || len(fakeShell.WarningOutput) != 2 {
		t.Error("should have warned on removing the file once")
	}
}

func TestKubectl(t *testing.T) {
	authTempPath = t.TempDir()

	k := NewDefaultK8S()
	k.deployExec.Endpoint.(*api.DefaultEndpoint).Fake()

	k.deployExec.Endpoint.(*api.DefaultEndpoint).MockResp(&api.DeployExecResponse{
		Server:    "server",
		Namespace: "ns",
		Path:      "path",
		Token:     "token",
		CA:        "ca",
	})

	fakeShell := &shell.FakeShell{}

	if _, err := k.Kubectl(fakeShell); !strings.Contains(err.Error(), "but did not auth") {
		t.Error("should get error before authenticating")
	}

	_, _ = k.Authenticate("foo", "bar")

	if cmd, _ := k.Kubectl(fakeShell); cmd.Cmd() != "kubectl" {
		t.Error("should use kubectl")
	}

	fakeShell.MockLookPath = errors.New("err")

	if cmd, _ := k.Kubectl(fakeShell); cmd.Cmd() != "kool" {
		t.Error("should use kool")
	}
}
