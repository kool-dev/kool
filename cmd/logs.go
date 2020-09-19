package cmd

import (
	"kool-dev/kool/cmd/shell"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// LogsFlags holds the flags for the logs command
type LogsFlags struct {
	Tail   int
	Follow bool
}

var logsCmd = &cobra.Command{
	Use:   "logs [options] [service...]",
	Short: "Displays log output from services.",
	Run:   runLogs,
}

var logsFlags = &LogsFlags{25, false}

var logsCmdOutputWriter shell.OutputWriter = shell.NewOutputWriter()

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().IntVarP(&logsFlags.Tail, "tail", "t", 25, "Number of lines to show from the end of the logs for each container. For value equal to 0, all lines will be shown.")
	logsCmd.Flags().BoolVarP(&logsFlags.Follow, "follow", "f", false, "Follow log output.")
}

func runLogs(cmd *cobra.Command, originalArgs []string) {
	var args []string = []string{"logs"}

	logsCmdOutputWriter.SetWriter(cmd.OutOrStdout())

	if logsFlags.Tail == 0 {
		args = append(args, "--tail", "all")
	} else {
		args = append(args, "--tail", strconv.Itoa(logsFlags.Tail))
	}

	if logsFlags.Follow {
		args = append(args, "--follow")
	}

	args = append(args, originalArgs...)

	err := shell.Interactive("docker-compose", args...)

	if err != nil {
		logsCmdOutputWriter.ExecError("", err)
		os.Exit(1)
	}
}
