package cmd

import (
	"fmt"
	"kool-dev/kool/api"
	"kool-dev/kool/cloud/k8s"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"

	"github.com/spf13/cobra"
)

// KoolDeployExec holds handlers and functions for using Deploy API
type KoolDeployExec struct {
	DefaultKoolService
	env   environment.EnvStorage
	cloud k8s.K8S
}

// NewDeployExecCommand inits Cobra command for kool deploy exec
func NewDeployExecCommand(deployExec *KoolDeployExec) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "exec SERVICE [COMMAND] [--] [ARG...]",
		Short: "Execute a command inside a running service container deployed to Kool Cloud",
		Long: `After deploying an application to Kool Cloud using 'kool deploy',
execute a COMMAND inside the specified SERVICE container (similar to an SSH session).
Must use a KOOL_API_TOKEN environment variable for authentication.`,
		Args: cobra.MinimumNArgs(1),
		Run:  DefaultCommandRunFunction(deployExec),

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().SetInterspersed(false)
	return
}

// NewKoolDeployExec creates a new pointer with default KoolDeployExec service dependencies
func NewKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newDefaultKoolService(),
		environment.NewEnvStorage(),
		k8s.NewDefaultK8S(),
	}
}

// Execute runs the deploy exec logic - integrating with Deploy API
func (e *KoolDeployExec) Execute(args []string) (err error) {
	var (
		domain, service, cloudService string

		kubectl builder.Command
	)

	if len(args) == 0 {
		err = fmt.Errorf("KoolDeployExec.Execute: required at least one argument")
		return
	}

	service = args[0]
	args = args[1:]

	if url := e.env.Get("KOOL_API_URL"); url != "" {
		api.SetBaseURL(url)
	}

	if domain = e.env.Get("KOOL_DEPLOY_DOMAIN"); domain == "" {
		err = fmt.Errorf("missing deploy domain (env KOOL_DEPLOY_DOMAIN)")
		return
	}

	if cloudService, err = e.cloud.Authenticate(domain, service); err != nil {
		return
	}

	defer e.cloud.Cleanup(e)

	if kubectl, err = e.cloud.Kubectl(e); err != nil {
		return
	}

	// finish building exec command
	kubectl.AppendArgs("exec", "-i")
	if e.IsTerminal() {
		kubectl.AppendArgs("-t")
	}
	kubectl.AppendArgs(cloudService, "-c", "default")
	kubectl.AppendArgs("--")
	if len(args) == 0 {
		args = []string{"bash"}
	}
	kubectl.AppendArgs(args...)

	err = e.Interactive(kubectl)
	return
}
