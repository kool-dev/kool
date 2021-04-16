package cmd

import (
	"kool-dev/kool/cloud/k8s"
	"kool-dev/kool/environment"
	"strings"
	"testing"
)

func fakeKoolDeployLogs() *KoolDeployLogs {
	return &KoolDeployLogs{
		*newFakeKoolService(),
		&KoolLogsFlags{},
		environment.NewFakeEnvStorage(),
		nil,
	}
}

func TestNewKoolDeployLogs(t *testing.T) {
	l := NewKoolDeployLogs()

	if l.Flags.Follow {
		t.Error("unexpected default Follow behaviour")
	}
	if l.Flags.Tail != 25 {
		t.Error("unexpected default Tail behaviour")
	}
	if _, ok := l.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("bad default type for env storage")
	}
	if _, ok := l.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("bad default type for env storage")
	}
	if _, ok := l.cloud.(*k8s.DefaultK8S); !ok {
		t.Error("bad default type for k8s cloud")
	}
}

func TestNewDeployLogsCommand(t *testing.T) {
	cmd := NewDeployLogsCommand(fakeKoolDeployLogs())

	if cmd.Flags().Lookup("tail") == nil {
		t.Error("missing flag: tailt")
	}

	if cmd.Flags().Lookup("follow") == nil {
		t.Error("missing flag: tailt")
	}
}

func TestKoolDeployLogsExecute(t *testing.T) {
	l := fakeKoolDeployLogs()
	args := []string{"foo"}

	l.env.Set("KOOL_API_URL", "api-url")

	if err := l.Execute(args); !strings.Contains(err.Error(), "missing deploy domain") {
		t.Error("should get error on missing domain")
	}

	// l.env.Set("KOOL_DEPLOY_DOMAIN", "deploy.domain")
}
