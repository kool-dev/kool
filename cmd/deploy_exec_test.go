package cmd

import (
	"kool-dev/kool/cloud/k8s"
	"kool-dev/kool/environment"
	"strings"
	"testing"
)

func newFakeKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newFakeKoolService(),
		environment.NewFakeEnvStorage(),
		nil,
	}
}

func TestNewKoolDeployExec(t *testing.T) {
	e := NewKoolDeployExec()

	if _, ok := e.env.(*environment.DefaultEnvStorage); !ok {
		t.Errorf("unexpected type for env storage")
	}
	if _, ok := e.cloud.(*k8s.DefaultK8S); !ok {
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

	// var domain string = "example.com"
	// e.env.Set("KOOL_DEPLOY_DOMAIN", domain)
}
