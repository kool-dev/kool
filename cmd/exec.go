package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// ExecFlags holds the flags for the start command
type ExecFlags struct {
	DisableTty bool
}

var execCmd = &cobra.Command{
	Use:                "exec [options] [service] [command]",
	Short:              "Execute a command within a running service container",
	Args:               cobra.MinimumNArgs(2),
	Run:                runExec,
	DisableFlagParsing: true,
}

var execFlags = &ExecFlags{false}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().BoolVarP(&execFlags.DisableTty, "disable-tty", "T", false, "Disables TTY")
}

func runExec(cmd *cobra.Command, originalArgs []string) {
	var (
		service string
		args    []string
	)

	if originalArgs[0] == "--disable-tty" || originalArgs[0] == "-T" {
		execFlags.DisableTty = true
		args = originalArgs[1:]
	} else {
		args = originalArgs
	}

	service = args[0]

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
	args = append(args, service)
	args = append(args, command...)

	err = shellInteractive("docker-compose", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
