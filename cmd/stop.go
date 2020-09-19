package cmd

import (
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/shell"
	"os"

	"github.com/spf13/cobra"
)

// StopFlags holds the flags for the stop command
type StopFlags struct {
	Purge bool
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop kool environment containers",
	Run:   runStop,
}

var stopFlags = &StopFlags{false}

var stopCmdOutputWriter shell.OutputWriter = shell.NewOutputWriter()

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVarP(&stopFlags.Purge, "purge", "", false, "Remove all persistent data from containers")
}

func runStop(cmd *cobra.Command, args []string) {
	stopCmdOutputWriter.SetWriter(cmd.OutOrStdout())
	var dependenciesChecker = checker.NewChecker()

	if err := dependenciesChecker.VerifyDependencies(); err != nil {
		stopCmdOutputWriter.ExecError("", err)
		os.Exit(1)
	}

	stopContainers(stopFlags.Purge)
}

func stopContainers(purge bool) {
	var (
		args []string
		err  error
	)

	args = []string{"down"}

	if purge {
		args = append(args, "--volumes", "--remove-orphans")
	}

	err = shell.NewCommander().Interactive("docker-compose", args...)

	if err != nil {
		stopCmdOutputWriter.ExecError("", err)
		os.Exit(1)
	}
}
