package cmd

import (
	"errors"
	"kool-dev/kool/api"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"strings"
	"testing"
)

type fakeExecCall struct {
	api.DefaultEndpoint

	err  error
	resp *api.ExecResponse
}

func (d *fakeExecCall) Call() (*api.ExecResponse, error) {
	return d.resp, d.err
}

func newFakeKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newFakeKoolService(),
		&builder.FakeCommand{}, // kubectl
		&builder.FakeCommand{}, // kool
		environment.NewFakeEnvStorage(),
		&fakeExecCall{
			DefaultEndpoint: *api.NewDefaultEndpoint(""),
		},
	}
}

func TestNewKoolDeployExec(t *testing.T) {
	e := NewKoolDeployExec()

	if _, ok := e.kubectl.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected type for kubectl command")
	}
	if _, ok := e.kool.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected type for kool command")
	}
	if _, ok := e.env.(*environment.DefaultEnvStorage); !ok {
		t.Errorf("unexpected type for env storage")
	}
	if _, ok := e.apiExec.(*api.DefaultExecCall); !ok {
		t.Errorf("unexpected type for apiExec endpoint")
	}
}

func TestKoolDeployExec(t *testing.T) {
	e := newFakeKoolDeployExec()
	err := e.Execute([]string{})

	if err == nil || !strings.Contains(err.Error(), "required at least one argument") {
		t.Errorf("expected: missing required parameter; got something else")
	}

	var service string = "my-service"
	err = e.Execute([]string{service})

	if err == nil || !strings.Contains(err.Error(), "missing deploy domain") {
		t.Errorf("expected: missing deploy domain; got something else")
	}

	var domain string = "example.com"
	e.env.Set("KOOL_DEPLOY_DOMAIN", domain)

	e.apiExec.(*fakeExecCall).err = errors.New("call error")

	if err = e.Execute([]string{service}); !errors.Is(err, e.apiExec.(*fakeExecCall).err) {
		t.Errorf("unexpected error from DeployExec call: %v", err)
	}

	e.apiExec.(*fakeExecCall).err = nil
	e.apiExec.(*fakeExecCall).resp = &api.ExecResponse{
		Server:    "server",
		Namespace: "ns",
		Path:      "path",
		Token:     "",
		CA:        "ca",
	}

	if err = e.Execute([]string{service}); !strings.Contains(err.Error(), "failed to generate access credentials") {
		t.Errorf("unexpected error from DeployExec call: %v", err)
	}

	e.apiExec.(*fakeExecCall).resp.Token = "token"
	authTempPath = t.TempDir()

	if err = e.Execute([]string{service, "foo", "bar"}); err != nil {
		t.Errorf("unexpected error from DeployExec call: %v", err)
	}

	args := e.kubectl.(*builder.FakeCommand).ArgsAppend
	if args[1] != "server" || args[3] != "token" || args[5] != "ns" || !strings.Contains(args[7], ".kool-cluster-CA") {
		t.Errorf("unexpected arguments to kubectl: %v", args)
	}

	if len(e.kool.(*builder.FakeCommand).ArgsAppend) > 0 {
		t.Errorf("should not have used kool")
	}

	e.kubectl.(*builder.FakeCommand).MockLookPathError = errors.New("not found")
	e.kubectl.(*builder.FakeCommand).MockCmd = "kub-foo"

	if err = e.Execute([]string{service, "foo", "bar"}); err != nil {
		t.Errorf("unexpected error from DeployExec call: %v", err)
	}

	args = e.kool.(*builder.FakeCommand).ArgsAppend

	if len(args) == 0 {
		t.Errorf("should have used kool")
	}

	if args[5] != "kub-foo" {
		t.Errorf("unexpected kubectl Cmd on kool: %v", args)
	}
}
