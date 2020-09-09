package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// LogsFlags holds the flags for the logs command
type LogsFlags struct {
	Tail int
	Follow bool
}

var logsCmd = &cobra.Command{
	Use:   "logs [options] [service...]",
	Short: "Displays log output from services.",
	Run: runLogs,
}

var logsFlags = &LogsFlags{0, false}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().IntVarP(&logsFlags.Tail, "tail", "t", 0, "Number of lines to show from the end of the logs for each container")
	logsCmd.Flags().BoolVarP(&logsFlags.Follow, "follow", "f", false, "Follow log output.")
}

func runLogs(cmd *cobra.Command, originalArgs []string) {
	var args []string = []string{"logs"}

	if logsFlags.Tail > 0 {
		args = append(args, "--tail", strconv.Itoa(logsFlags.Tail))
	}

	if logsFlags.Follow {
		args = append(args, "--follow")
	}

	for _, service := range originalArgs {
		args = append(args, service)
	}

	err := shellInteractive("docker-compose", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
