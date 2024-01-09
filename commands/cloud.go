package commands

import (
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/services/cloud/api"

	"github.com/spf13/cobra"
)

// KoolCloudDeployFlags holds the flags for the kool cloud deploy command
type KoolCloudFlags struct {
	Token        string // env: KOOL_API_TOKEN
	DeployDomain string // env: KOOL_DEPLOY_DOMAIN
}

type Cloud struct {
	DefaultKoolService

	flags *KoolCloudFlags
	env   environment.EnvStorage
}

func NewCloud() *Cloud {
	return &Cloud{
		*newDefaultKoolService(),
		&KoolCloudFlags{},
		environment.NewEnvStorage(),
	}
}

func AddKoolCloud(root *cobra.Command) {
	var (
		cloud    = NewCloud()
		cloudCmd = NewCloudCommand(cloud)
	)

	cloudCmd.AddCommand(NewDeployCommand(NewKoolDeploy(cloud)))
	cloudCmd.AddCommand(NewDeployExecCommand(NewKoolDeployExec()))
	cloudCmd.AddCommand(NewDeployDestroyCommand(NewKoolDeployDestroy()))
	cloudCmd.AddCommand(NewDeployLogsCommand(NewKoolDeployLogs()))
	cloudCmd.AddCommand(NewSetupCommand(NewKoolCloudSetup()))

	root.AddCommand(cloudCmd)
}

// NewCloudCommand initializes new kool cloud command
func NewCloudCommand(cloud *Cloud) (cloudCmd *cobra.Command) {
	cloudCmd = &cobra.Command{
		Use:     "cloud COMMAND [flags]",
		Short:   "Interact with Kool.dev Cloud and manage your deployments.",
		Long:    "The cloud subcommand encapsulates a set of APIs to interact with Kool.dev Cloud and deploy, access and tail logs from your deployments.",
		Example: `kool cloud deploy`,
		//	add cobra usage help content
		DisableFlagsInUseLine: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			// calls root PersistentPreRunE
			var (
				requiredFlags bool           = cmd.Use != "setup"
				root          *cobra.Command = cmd
			)

			for root.HasParent() {
				root = root.Parent()
			}
			if err = root.PersistentPreRunE(cmd, args); err != nil {
				return
			}

			if url := cloud.env.Get("KOOL_API_URL"); url != "" {
				api.SetBaseURL(url)
			}

			// if no domain is set, we try to get it from the environment
			if cloud.flags.DeployDomain == "" && cloud.env.Get("KOOL_DEPLOY_DOMAIN") == "" {
				if requiredFlags {
					err = fmt.Errorf("missing deploy domain - please set it via --domain or KOOL_DEPLOY_DOMAIN environment variable")
					return
				}
			} else if cloud.flags.DeployDomain != "" {
				// shares the flag via environment variable
				cloud.env.Set("KOOL_DEPLOY_DOMAIN", cloud.flags.DeployDomain)
			}

			// if no token is set, we try to get it from the environment
			if cloud.flags.Token == "" && cloud.env.Get("KOOL_API_TOKEN") == "" {
				if requiredFlags {
					err = fmt.Errorf("missing Kool.dev Cloud API token - please set it via --token or KOOL_API_TOKEN environment variable")
					return
				}
			} else if cloud.flags.Token != "" {
				cloud.env.Set("KOOL_API_TOKEN", cloud.flags.Token)
			}

			return
		},
	}

	cloudCmd.PersistentFlags().StringVarP(&cloud.flags.Token, "token", "", "", "Token to authenticate with Kool.dev Cloud API")
	cloudCmd.PersistentFlags().StringVarP(&cloud.flags.DeployDomain, "domain", "", "", "Environment domain name to deploy to")

	return
}
