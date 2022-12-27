package commands

import "github.com/spf13/cobra"

func AddKoolCloud(root *cobra.Command) {
	var (
		cloudCmd = NewCloudCommand()
	)

	cloudCmd.AddCommand(NewDeployCommand(NewKoolDeploy()))
	cloudCmd.AddCommand(NewDeployExecCommand(NewKoolDeployExec()))
	cloudCmd.AddCommand(NewDeployDestroyCommand(NewKoolDeployDestroy()))
	cloudCmd.AddCommand(NewDeployLogsCommand(NewKoolDeployLogs()))
	cloudCmd.AddCommand(NewSetupCommand(NewKoolCloudSetup()))

	root.AddCommand(cloudCmd)
}

// NewCloudCommand initializes new kool cloud command
func NewCloudCommand() (cloudCmd *cobra.Command) {
	cloudCmd = &cobra.Command{
		Use:     "cloud COMMAND [flags]",
		Short:   "Interact with Kool Cloud and manage your deployments.",
		Long:    "The cloud subcommand encapsulates a set of APIs to interact with Kool Cloud and deploy, access and tail logs from your deployments.",
		Example: `kool cloud deploy`,
		//	add cobra usage help content
		DisableFlagsInUseLine: true,
	}

	return
}
