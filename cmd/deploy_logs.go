package cmd

import (
	"fmt"
	"kool-dev/kool/api"
	"kool-dev/kool/cloud/k8s"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"

	"github.com/spf13/cobra"
)

// KoolDeployLogs holds handlers and functions for using Deploy API
type KoolDeployLogs struct {
	DefaultKoolService
	Flags *KoolLogsFlags
	env   environment.EnvStorage
	cloud k8s.K8S
}

// NewDeployLogsCommand inits Cobra command for kool deploy logs
func NewDeployLogsCommand(deployLogs *KoolDeployLogs) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "logs [OPTIONS] SERVICE",
		Short: "See the logs of running service container deployed to Kool Cloud",
		Long: `After deploying an application to Kool Cloud using 'kool deploy',
you can see the logs from the specified SERVICE container.
Must use a KOOL_API_TOKEN environment variable for authentication.`,
		Args: cobra.ExactArgs(1),
		Run:  DefaultCommandRunFunction(deployLogs),

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().IntVarP(&deployLogs.Flags.Tail, "tail", "t", 25, "Number of lines to show from the end of the logs for each container. A value equal to 0 will show all lines.")
	cmd.Flags().BoolVarP(&deployLogs.Flags.Follow, "follow", "f", false, "Follow log output.")
	return
}

// NewKoolDeployLogs creates a new pointer with default KoolDeployLogs service dependencies
func NewKoolDeployLogs() *KoolDeployLogs {
	return &KoolDeployLogs{
		*newDefaultKoolService(),
		&KoolLogsFlags{25, false},
		environment.NewEnvStorage(),
		k8s.NewDefaultK8S(),
	}
}

// Execute runs the deploy logs logic - integrating with API and K8S
func (e *KoolDeployLogs) Execute(args []string) (err error) {
	var (
		domain  string
		service string
		kubectl builder.Command
	)

	if len(args) == 0 {
		err = fmt.Errorf("KoolDeployLogs.Execute: required at least one argument")
		return
	}

	service = args[0]

	if url := e.env.Get("KOOL_API_URL"); url != "" {
		api.SetBaseURL(url)
	}

	if domain = e.env.Get("KOOL_DEPLOY_DOMAIN"); domain == "" {
		err = fmt.Errorf("missing deploy domain (env KOOL_DEPLOY_DOMAIN)")
		return
	}

	if err = e.cloud.Authenticate(domain, service); err != nil {
		return
	}

	defer e.cloud.Cleanup(e)

	if kubectl, err = e.cloud.Kubectl(e); err != nil {
		return
	}

	kubectl.AppendArgs("logs")
	if e.Flags.Follow {
		kubectl.AppendArgs("-f")
	}
	if e.Flags.Tail > 0 {
		kubectl.AppendArgs("--tail", fmt.Sprintf("%d", e.Flags.Tail))
	}
	kubectl.AppendArgs(e.cloud.CloudService())
	kubectl.AppendArgs("-c", "default")

	err = e.Interactive(kubectl)
	return
}
