package cmd

import (
	"kool-dev/kool/cmd/shell"
	"os"

	"github.com/spf13/cobra"
)

// ExecFlags holds the flags for the start command
type ExecFlags struct {
	DisableTty   bool
	EnvVariables []string
	Detach       bool
}

var execCmd = &cobra.Command{
	Use:   "exec [options] [service] [command]",
	Short: "Execute a command within a running service container",
	Args:  cobra.MinimumNArgs(2),
	Run:   runExec,
}

var execFlags = &ExecFlags{false, []string{}, false}

var execCmdOutputWriter shell.OutputWriter = shell.NewOutputWriter()

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().BoolVarP(&execFlags.DisableTty, "disable-tty", "T", false, "Disables TTY")
	execCmd.Flags().StringArrayVarP(&execFlags.EnvVariables, "env", "e", []string{}, "Environment variables")
	execCmd.Flags().BoolVarP(&execFlags.Detach, "detach", "d", false, "Detached mode: Run command in the background")

	//After a non-flag arg, stop parsing flags
	execCmd.Flags().SetInterspersed(false)
}

func runExec(cmd *cobra.Command, args []string) {
	var service string = args[0]

	execCmdOutputWriter.SetWriter(cmd.OutOrStdout())
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

	if execFlags.Detach {
		args = append(args, "--detach")
	}

	args = append(args, service)
	args = append(args, command...)

	err = shell.Interactive("docker-compose", args...)

	if err != nil {
		execCmdOutputWriter.ExecError("", err)
		os.Exit(1)
	}
}
