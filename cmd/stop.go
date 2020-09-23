package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"

	"github.com/spf13/cobra"
)

// KoolStopFlags holds the flags for the kool stop command
type KoolStopFlags struct {
	Purge bool
}

// KoolStop holds handlers and functions to implement the stop command logic
type KoolStop struct {
	DefaultKoolService
	Flags *KoolStopFlags

	check  checker.Checker
	doStop builder.Command
}

func init() {
	var (
		stop    = NewKoolStop()
		stopCmd = NewStopCommand(stop)
	)

	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVarP(&stop.Flags.Purge, "purge", "", false, "Remove all persistent data from volume mounts on containers")
}

// NewKoolStop creates a new handler for stop logic with default dependencies
func NewKoolStop() *KoolStop {
	return &KoolStop{
		*newDefaultKoolService(),
		&KoolStopFlags{false},
		checker.NewChecker(),
		builder.NewCommand("docker-compose", "down"),
	}
}

// Execute runs the stop logic with incoming arguments.
func (s *KoolStop) Execute(args []string) (err error) {
	if err = s.check.Check(); err != nil {
		return
	}

	if s.Flags.Purge {
		s.doStop.AppendArgs("--volumes", "--remove-orphans")
	}

	err = s.doStop.Interactive()
	return
}

// NewStopCommand initializes new kool stop command
func NewStopCommand(stop *KoolStop) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop all running containers started with 'kool start' command",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			stop.SetWriter(cmd.OutOrStdout())

			if err := stop.Execute(args); err != nil {
				stop.Error(err)
				stop.Exit(1)
			}
		},
	}
}
