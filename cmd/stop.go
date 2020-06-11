package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// StopFlags holds the flags for the start command
type StopFlags struct {
	Purge bool
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop fwd environment containers",
	Run:   runStop,
}

var stopFlags = &StopFlags{false}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVarP(&stopFlags.Purge, "purge", "", false, "Remove all persistent data from containers")
}

func runStop(cmd *cobra.Command, args []string) {
	stopContainers()
}

func stopContainers() {
	var (
		args []string
		err  error
		out  string
	)

	args = []string{"down"}

	if stopFlags.Purge {
		args = append(args, "--volumes", "--remove-orphans")
	}

	out, err = shellExec("docker-compose", args...)

	if err != nil {
		execError(out, err)
		os.Exit(1)
	}
}
