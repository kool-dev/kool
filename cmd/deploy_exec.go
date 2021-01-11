package cmd

import (
	"fmt"
	"io/ioutil"
	"kool-dev/kool/api"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"os"

	"github.com/spf13/cobra"
)

// KoolDeployExec holds handlers and functions for using Deploy API
type KoolDeployExec struct {
	DefaultKoolService

	kubectl, kool builder.Command
	env           environment.EnvStorage
	apiExec       api.ExecCall
}

// NewDeployExecCommand initializes new kool deploy Cobra command
func NewDeployExecCommand(deployExec *KoolDeployExec) *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Executes a command in your deployed application on Kool cloud",
		Run:   DefaultCommandRunFunction(deployExec),
	}
}

// NewKoolDeployExec creates a new pointer with default KoolDeployExec service dependencies
func NewKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newDefaultKoolService(),
		builder.NewCommand("kubectl"),
		builder.NewCommand("kool"),
		environment.NewEnvStorage(),
		api.NewDefaultExecCall(),
	}
}

// Execute runs the deploy exec logic - integrating with Deploy API
func (e *KoolDeployExec) Execute(args []string) (err error) {
	var (
		domain string
		resp   *api.ExecResponse
	)

	e.Println("kool deploy exec - start")

	if domain = e.env.Get("KOOL_DEPLOY_DOMAIN"); domain == "" {
		err = fmt.Errorf("missing deploy domain (env KOOL_DEPLOY_DOMAIN)")
		return
	}

	e.apiExec.Body().Set("domain", domain)

	if resp, err = e.apiExec.Call(); err != nil {
		return
	}

	if resp.Token == "" {
		err = fmt.Errorf("failed to generate access credentials to cloud deploy")
		return
	}

	CAPath := fmt.Sprintf("%s/.kool-cluster-CA", os.TempDir())
	if err = ioutil.WriteFile(CAPath, []byte(resp.CA), os.ModePerm); err != nil {
		return
	}

	e.kubectl.AppendArgs("--token", resp.Token, "-n", resp.Namespace, "exec", "-i")
	// e.kubectl.AppendArgs("--insecure-skip-tls-verify", "true", "--server", resp.Server)
	e.kubectl.AppendArgs("--certificate-authority", CAPath)
	if e.IsTerminal() {
		e.kubectl.AppendArgs("-t")
	}
	e.kubectl.AppendArgs(resp.Path, "--")
	if len(args) == 0 {
		args = []string{"bash"}
	}
	e.kubectl.AppendArgs(args...)

	if e.LookPath(e.kubectl) == nil {
		// the command is available on current PATH, so let's
		// just execute it
		err = e.Interactive(e.kubectl)
		return
	}

	// we do not have 'kubectl' on current path... let's use a container!
	e.kool.AppendArgs("docker", "--", "kooldev/toolkit:full", e.kubectl.Cmd())
	e.kool.AppendArgs(e.kubectl.Args()...)

	err = e.Interactive(e.kool)
	return
}
