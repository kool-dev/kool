package cmd

import (
	"kool-dev/kool/api"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"strings"
	"testing"
)

func newFakeKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newFakeKoolService(),
		&builder.FakeCommand{}, // kubectl
		&builder.FakeCommand{}, // kool
		environment.NewFakeEnvStorage(),
		nil, // api.NewDefaultExecCall()
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

	// to continue we need to mock api.Endpoint
	// var domain string = "example.com"
	// e.env.Set("KOOL_DEPLOY_DOMAIN", domain)
	// err = e.Execute([]string{service})
}
