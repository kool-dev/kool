package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// ExecFlags holds the flags for the start command
type ExecFlags struct {
	DisableTty bool
	EnvVariables []string
}

var execCmd = &cobra.Command{
	Use:                "exec [options] [service] [command]",
	Short:              "Execute a command within a running service container",
	Args:               cobra.MinimumNArgs(2),
	Run:                runExec,
}

var execFlags = &ExecFlags{false, []string{}}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().BoolVarP(&execFlags.DisableTty, "disable-tty", "T", false, "Disables TTY")
	execCmd.Flags().StringArrayVarP(&execFlags.EnvVariables, "env", "e", []string{}, "Environment variables")

	//After a non-flag arg, stop parsing flags
	execCmd.Flags().SetInterspersed(false)
}

func runExec(cmd *cobra.Command, args []string) {
	var service string = args[0]

	dockerComposeExec(service, args[1:]...)
}

func dockerComposeExec(service string, command ...string) {
	var (
		err  error
		args []string
	)

	args = []string{"exec"}
	if disableTty := os.Getenv("KOOL_TTY_DISABLE"); execFlags.DisableTty || disableTty == "1" || disableTty == "true" {
		args = append(args, "-T")
	}
	if asuser := os.Getenv("KOOL_ASUSER"); asuser != "" {
		args = append(args, "--user", asuser)
	}

	if len(execFlags.EnvVariables) > 0 {
		for _, envVar := range execFlags.EnvVariables {
			args = append(args, "--env", envVar)
		}
	}

	args = append(args, service)
	args = append(args, command...)

	err = shellInteractive("docker-compose", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
