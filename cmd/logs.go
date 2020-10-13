package cmd

import (
	"kool-dev/kool/cmd/builder"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// KoolLogsFlags holds the flags for the logs command
type KoolLogsFlags struct {
	Tail   int
	Follow bool
}

// KoolLogs holds handlers and functions to implement the logs command logic
type KoolLogs struct {
	DefaultKoolService
	Flags *KoolLogsFlags

	list builder.Command
	logs builder.Command
}

func init() {
	var (
		logs    = NewKoolLogs()
		logsCmd = NewLogsCommand(logs)
	)

	rootCmd.AddCommand(logsCmd)
}

// NewKoolLogs creates a new handler for logs logic
func NewKoolLogs() *KoolLogs {
	return &KoolLogs{
		*newDefaultKoolService(),
		&KoolLogsFlags{25, false},
		builder.NewCommand("docker-compose", "ps", "-aq"),
		builder.NewCommand("docker-compose", "logs"),
	}
}

// Execute runs the logs logic with incoming arguments.
func (l *KoolLogs) Execute(args []string) (err error) {
	var services string

	if services, err = l.list.Exec(args...); err != nil {
		return
	}

	if services = strings.TrimSpace(services); services == "" {
		l.Warning("There are no containers")
		return
	}

	if l.Flags.Tail == 0 {
		l.logs.AppendArgs("--tail", "all")
	} else {
		l.logs.AppendArgs("--tail", strconv.Itoa(l.Flags.Tail))
	}

	if l.Flags.Follow {
		l.logs.AppendArgs("--follow")
	}

	err = l.logs.Interactive(args...)
	return
}

// NewLogsCommand initializes new kool logs command
func NewLogsCommand(logs *KoolLogs) (logsCmd *cobra.Command) {
	logsCmd = &cobra.Command{
		Use:   "logs [options] [service...]",
		Short: "Displays log output from services.",
		Run:   DefaultCommandRunFunction(logs),
	}

	logsCmd.Flags().IntVarP(&logs.Flags.Tail, "tail", "t", 25, "Number of lines to show from the end of the logs for each container. For value equal to 0, all lines will be shown.")
	logsCmd.Flags().BoolVarP(&logs.Flags.Follow, "follow", "f", false, "Follow log output.")
	return
}
