package k8s

import (
	"errors"
	"kool-dev/kool/api"
	"kool-dev/kool/cmd/shell"
	"os"
	"strings"
	"testing"
)

// fake api.ExecCall
type fakeExecCall struct {
	api.DefaultEndpoint

	err  error
	resp *api.ExecResponse
}

func (d *fakeExecCall) Call() (*api.ExecResponse, error) {
	return d.resp, d.err
}

func newFakeExecCall() *fakeExecCall {
	return &fakeExecCall{
		DefaultEndpoint: *api.NewDefaultEndpoint(""),
	}
}

// fake shell.OutputWritter
type fakeOutputWritter struct {
	warned []interface{}
}

func (*fakeOutputWritter) Println(args ...interface{}) {
}

func (*fakeOutputWritter) Printf(s string, args ...interface{}) {
}

func (f *fakeOutputWritter) Warning(args ...interface{}) {
	f.warned = append(f.warned, args...)
}

func (*fakeOutputWritter) Success(args ...interface{}) {
}

func TestNewDefaultK8S(t *testing.T) {
	k := NewDefaultK8S()
	if _, ok := k.apiExec.(*api.DefaultExecCall); !ok {
		t.Error("invalid type on apiExec")
	}
}

func TestAuthenticate(t *testing.T) {
	k := &DefaultK8S{
		apiExec: newFakeExecCall(),
	}

	expectedErr := errors.New("call error")
	k.apiExec.(*fakeExecCall).err = expectedErr

	if _, err := k.Authenticate("foo", "bar"); !errors.Is(err, expectedErr) {
		t.Error("unexpected error return from Authenticate")
	}

	k.apiExec.(*fakeExecCall).err = nil
	k.apiExec.(*fakeExecCall).resp = &api.ExecResponse{
		Server:    "server",
		Namespace: "ns",
		Path:      "path",
		Token:     "",
		CA:        "ca",
	}

	if _, err := k.Authenticate("foo", "bar"); !strings.Contains(err.Error(), "failed to generate access credentials") {
		t.Errorf("unexpected error from DeployExec call: %v", err)
	}

	k.apiExec.(*fakeExecCall).resp.Token = "token"
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

	fakeOut := &fakeOutputWritter{}

	k.Cleanup(fakeOut)

	if len(fakeOut.warned) != 0 {
		t.Error("should not have warned on removing the file")
	}

	authTempPath = t.TempDir() + "test"
	k.Cleanup(fakeOut)

	if len(fakeOut.warned) != 2 {
		t.Error("should have warned on removing the file once")
	}
}

func TestKubectl(t *testing.T) {
	authTempPath = t.TempDir()

	k := &DefaultK8S{
		apiExec: newFakeExecCall(),
	}

	k.apiExec.(*fakeExecCall).resp = &api.ExecResponse{
		Server:    "server",
		Namespace: "ns",
		Path:      "path",
		Token:     "token",
		CA:        "ca",
	}

	fakeShell := &shell.FakeShell{}

	if _, err := k.Kubectl(fakeShell); !strings.Contains(err.Error(), "but did not auth") {
		t.Error("should get error before authenticating")
	}

	k.Authenticate("foo", "bar")

	if cmd, _ := k.Kubectl(fakeShell); cmd.Cmd() != "kubectl" {
		t.Error("should use kubectl")
	}

	fakeShell.MockLookPath = errors.New("err")

	if cmd, _ := k.Kubectl(fakeShell); cmd.Cmd() != "kool" {
		t.Error("should use kool")
	}
}
