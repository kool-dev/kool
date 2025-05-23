package k8s

import (
	"errors"
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/api"
	"os"
	"path/filepath"
)

type K8S interface {
	Authenticate(string, string, string) (string, error)
	Kubectl(shell.PathChecker) (builder.Command, error)
	Cleanup(shell.OutputWritter)
}

type DefaultK8S struct {
	deployExec api.DeployExec
	resp       *api.DeployExecResponse
}

var authTempPath = "/tmp"

// NewDefaultK8S returns a new pointer for DefaultK8S with dependencies
func NewDefaultK8S() *DefaultK8S {
	return &DefaultK8S{
		deployExec: *api.NewDeployExec(),
	}
}

func (k *DefaultK8S) Authenticate(domain, service, cluster string) (cloudService string, err error) {
	k.deployExec.Body().Set("domain", domain)
	k.deployExec.Body().Set("service", service)
	k.deployExec.Body().Set("cluster", cluster)

	if k.resp, err = k.deployExec.Call(); err != nil {
		return
	}

	if k.resp.Token == "" {
		err = fmt.Errorf("failed to generate access credentials to cloud deploy")
		return
	}

	cloudService = k.resp.Path

	err = os.WriteFile(k.getTempCAPath(), []byte(k.resp.CA), os.ModePerm)
	return
}

func (k *DefaultK8S) Kubectl(looker shell.PathChecker) (kube builder.Command, err error) {
	if k.resp == nil {
		err = errors.New("calling kubectl but did not authenticate")
		return
	}

	kube = builder.NewCommand("kubectl")

	kube.AppendArgs("--server", k.resp.Server)
	kube.AppendArgs("--token", k.resp.Token)
	kube.AppendArgs("--namespace", k.resp.Namespace)
	kube.AppendArgs("--certificate-authority", k.getTempCAPath())

	if looker.LookPath(kube) != nil {
		// we do not have 'kubectl' on current path... let's use a container!
		kool := builder.NewCommand("kool")
		kool.AppendArgs(
			"docker", "--",
			"-v", fmt.Sprintf("%s:%s", k.getTempCAPath(), k.getTempCAPath()),
			"kooldev/toolkit:full",
			kube.Cmd(),
		)
		kool.AppendArgs(kube.Args()...)
		kube = kool
	}

	return
}

func (k *DefaultK8S) Cleanup(out shell.OutputWritter) {
	if err := os.Remove(k.getTempCAPath()); err != nil {
		out.Warning("failed to clear up temporary file; error:", err.Error())
	}
}

func (k *DefaultK8S) getTempCAPath() string {
	return filepath.Join(authTempPath, ".kool-cluster-CA")
}
