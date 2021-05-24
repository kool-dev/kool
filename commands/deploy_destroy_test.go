package commands

import (
	"errors"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/api"
	"strings"
	"testing"
)

func TestNewDeployDestroyCommand(t *testing.T) {
	destroy := NewKoolDeployDestroy()
	cmd := NewDeployDestroyCommand(destroy)
	if cmd.Use != "destroy" {
		t.Errorf("bad command use: %s", cmd.Use)
	}

	if _, ok := destroy.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("unexpected default env on destroy")
	}
}

type fakeDestroyCall struct {
	api.DefaultEndpoint

	err  error
	resp *api.DestroyResponse
}

func (d *fakeDestroyCall) Call() (*api.DestroyResponse, error) {
	return d.resp, d.err
}

func TestDeployDestroyExec(t *testing.T) {
	destroy := &KoolDeployDestroy{
		*newFakeKoolService(),
		environment.NewFakeEnvStorage(),
		&fakeDestroyCall{
			DefaultEndpoint: *api.NewDefaultEndpoint(""),
		},
	}

	destroy.env.Set("KOOL_API_TOKEN", "fake token")
	destroy.env.Set("KOOL_API_URL", "fake-url")

	args := []string{}

	if err := destroy.Execute(args); !strings.Contains(err.Error(), "missing deploy domain") {
		t.Errorf("unexpected error - expected missing deploy domain, got: %v", err)
	}

	destroy.env.Set("KOOL_DEPLOY_DOMAIN", "domain.com")

	destroy.apiDestroy.(*fakeDestroyCall).err = errors.New("failed call")

	if err := destroy.Execute(args); !strings.Contains(err.Error(), "failed call") {
		t.Errorf("unexpected error - expected failed call, got: %v", err)
	}

	destroy.apiDestroy.(*fakeDestroyCall).err = nil
	resp := new(api.DestroyResponse)
	resp.Environment.ID = 100
	destroy.apiDestroy.(*fakeDestroyCall).resp = resp

	if err := destroy.Execute(args); err != nil {
		t.Errorf("unexpected error, got: %v", err)
	}

	if !strings.Contains(destroy.shell.(*shell.FakeShell).SuccessOutput[0].(string), "ID: 100") {
		t.Errorf("did not get success message")
	}
}
