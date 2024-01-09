package commands

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/services/cloud/api"
	"kool-dev/kool/services/cloud/k8s"

	"github.com/spf13/cobra"
)

// KoolDeployExec holds handlers and functions for using Deploy API
type KoolDeployExec struct {
	DefaultKoolService
	Flags *KoolDeployExecFlags
	env   environment.EnvStorage
	cloud k8s.K8S
}

// KoolDeployExecFlags holds flags to kool deploy exec command
type KoolDeployExecFlags struct {
	Container string
}

// NewDeployExecCommand inits Cobra command for kool deploy exec
func NewDeployExecCommand(deployExec *KoolDeployExec) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "exec SERVICE [COMMAND] [--] [ARG...]",
		Short: "Execute a command inside a running service container deployed to Kool.dev Cloud",
		Long: `After deploying an application to Kool.dev Cloud using 'kool deploy',
execute a COMMAND inside the specified SERVICE container (similar to an SSH session).
Must use a KOOL_API_TOKEN environment variable for authentication.`,
		Args: cobra.MinimumNArgs(1),
		RunE: DefaultCommandRunFunction(deployExec),

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().SetInterspersed(false)
	cmd.Flags().StringVarP(&deployExec.Flags.Container, "container", "c", "default", "Container target.")
	return
}

// NewKoolDeployExec creates a new pointer with default KoolDeployExec service dependencies
func NewKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newDefaultKoolService(),
		&KoolDeployExecFlags{"default"},
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

	defer e.cloud.Cleanup(e.Shell())

	if kubectl, err = e.cloud.Kubectl(e.Shell()); err != nil {
		return
	}

	// finish building exec command
	kubectl.AppendArgs("exec", "-i")
	if e.Shell().IsTerminal() {
		kubectl.AppendArgs("-t")
	}
	kubectl.AppendArgs(cloudService, "-c", e.Flags.Container)
	kubectl.AppendArgs("--")
	if len(args) == 0 {
		args = []string{"bash"}
	}
	kubectl.AppendArgs(args...)

	err = e.Shell().Interactive(kubectl)
	return
}
