package cmd

import (
	"errors"
	"kool-dev/kool/cloud/k8s"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"strings"
	"testing"
)

func newFakeKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newFakeKoolService(),
		environment.NewFakeEnvStorage(),
		&fakeK8S{},
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

	var args = []string{"my-service"}

	if err = e.Execute(args); err == nil || !strings.Contains(err.Error(), "missing deploy domain") {
		t.Errorf("expected: missing deploy domain; got something else")
	}

	var domain string = "example.com"
	e.env.Set("KOOL_DEPLOY_DOMAIN", domain)

	mock := e.cloud.(*fakeK8S)
	mock.MockAuthenticateErr = errors.New("auth error")

	if err = e.Execute(args); !errors.Is(err, mock.MockAuthenticateErr) {
		t.Error("should return auth error")
	}

	mock.MockAuthenticateErr = nil
	mock.MockAuthenticateCloudService = "cloud-service"
	mock.MockKubectlErr = errors.New("kube error")

	if err = e.Execute(args); !errors.Is(err, mock.MockKubectlErr) {
		t.Error("should return kube error")
	}

	fakeKubectl := &builder.FakeCommand{}
	mock.MockKubectlErr = nil
	mock.MockKubectlKube = fakeKubectl

	fakeKubectl.MockInteractiveError = errors.New("interactive error")

	if err = e.Execute(args); !errors.Is(err, fakeKubectl.MockInteractiveError) {
		t.Error("should return interactive error")
	}

	fakeKubectl = &builder.FakeCommand{}
	mock.MockKubectlKube = fakeKubectl
	fakeKubectl.MockInteractiveError = nil
	e.term.(*shell.FakeTerminalChecker).MockIsTerminal = true

	if err = e.Execute(args); err != nil {
		t.Error("unexpected error")
	}

	str := strings.Join(fakeKubectl.ArgsAppend, " ")

	if !strings.Contains(str, "exec -i -t cloud-service -c default -- bash") {
		t.Error("bad kubectl command args")
	}
}
