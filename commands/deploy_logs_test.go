package commands

import (
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/k8s"
	"strings"
	"testing"
)

type fakeK8S struct {
	// Authenticate
	CalledAuthenticate             bool
	CalledAuthenticateParamDomain  string
	CalledAuthenticateParamService string
	MockAuthenticateCloudService   string
	MockAuthenticateErr            error

	// Kubectl
	CalledKubectl            bool
	CalledKubectlParamLooker shell.PathChecker
	MockKubectlKube          builder.Command
	MockKubectlErr           error

	// Cleanup
	CalledCleanup         bool
	CalledCleanupParamOut shell.OutputWritter
}

func (f *fakeK8S) Authenticate(domain, service string) (cloudService string, err error) {
	f.CalledAuthenticate = true
	f.CalledAuthenticateParamDomain = domain
	f.CalledAuthenticateParamService = service

	cloudService = f.MockAuthenticateCloudService
	err = f.MockAuthenticateErr
	return
}

func (f *fakeK8S) Kubectl(looker shell.PathChecker) (kube builder.Command, err error) {
	f.CalledKubectl = true
	f.CalledKubectlParamLooker = looker
	kube = f.MockKubectlKube
	err = f.MockKubectlErr
	return
}

func (f *fakeK8S) Cleanup(out shell.OutputWritter) {
	f.CalledCleanup = true
	f.CalledCleanupParamOut = out
}

func fakeKoolDeployLogs() *KoolDeployLogs {
	return &KoolDeployLogs{
		*newFakeKoolService(),
		&KoolDeployLogsFlags{},
		environment.NewFakeEnvStorage(),
		&fakeK8S{},
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

	l.env.Set("KOOL_DEPLOY_DOMAIN", "deploy.domain")

	l.cloud.(*fakeK8S).MockAuthenticateErr = errors.New("authenticate error")

	if err := l.Execute(args); !errors.Is(err, l.cloud.(*fakeK8S).MockAuthenticateErr) {
		t.Error("should get error on authenticate")
	}

	l.cloud.(*fakeK8S).MockAuthenticateErr = nil
	l.cloud.(*fakeK8S).MockAuthenticateCloudService = "app"
	l.cloud.(*fakeK8S).MockKubectlErr = errors.New("kubectl error")

	if err := l.Execute(args); !errors.Is(err, l.cloud.(*fakeK8S).MockKubectlErr) {
		t.Error("should get error on kubectl")
	}

	l.cloud.(*fakeK8S).MockKubectlErr = nil
	l.cloud.(*fakeK8S).MockKubectlKube = &builder.FakeCommand{
		MockInteractiveError: errors.New("interactive error"),
	}

	if err := l.Execute(args); !errors.Is(err, l.cloud.(*fakeK8S).MockKubectlKube.(*builder.FakeCommand).MockInteractiveError) {
		t.Error("should get error on kubectl - interactive")
	}

	fakeKubectl := &builder.FakeCommand{}
	l.cloud.(*fakeK8S).MockKubectlKube = fakeKubectl
	l.Flags.Follow = true
	l.Flags.Tail = 25

	if err := l.Execute(args); err != nil {
		t.Error("unexpected error")
	}

	str := strings.Join(fakeKubectl.ArgsAppend, " ")

	if !strings.Contains(str, "logs -f --tail 25") {
		t.Error("bad kubectl command - missing logs -f : " + str)
	}
}
