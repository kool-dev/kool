package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// ExecFlags holds the flags for the start command
type ExecFlags struct {
}

var execCmd = &cobra.Command{
	Use:                "exec [service] [command]",
	Short:              "Execute a command within a running service container",
	Args:               cobra.MinimumNArgs(2),
	Run:                runExec,
	DisableFlagParsing: true,
}

var execFlags = &ExecFlags{}

func init() {
	rootCmd.AddCommand(execCmd)
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
	if disableTty := os.Getenv("KOOL_TTY_DISABLE"); disableTty == "1" || disableTty == "true" {
		args = append(args, "-T")
	}
	args = append(args, service)
	args = append(args, command...)

	err = shellInteractive("docker-compose", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
