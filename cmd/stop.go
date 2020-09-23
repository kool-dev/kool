package cmd

import (
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/shell"

	"github.com/spf13/cobra"
)

// KoolStop holds the logic
type KoolStop struct {
	exiter shell.Exiter
	out    shell.OutputWriter
}

// Execute runs the stop logic with incoming arguments.
func (s *KoolStop) Execute(args []string) (err error) {
	var dependenciesChecker = checker.NewChecker()

	if err = dependenciesChecker.Check(); err != nil {
		return
	}

	err = s.stop(stopFlags.Purge)
	return
}

// NewKoolStop creates a new handler for stop logic with default dependencies
func NewKoolStop() *KoolStop {
	return &KoolStop{
		shell.NewExiter(),
		shell.NewOutputWriter(),
	}
}

// StopFlags holds the flags for the stop command
type StopFlags struct {
	Purge bool
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all running containers started with 'kool start' command",
	Run: func(cmd *cobra.Command, args []string) {
		var stop = NewKoolStop()

		stop.out.SetWriter(cmd.OutOrStdout())

		if err := stop.Execute(args); err != nil {
			stop.out.Error(err)
			stop.exiter.Exit(1)
		}
	},
}

var stopFlags = &StopFlags{false}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVarP(&stopFlags.Purge, "purge", "", false, "Remove all persistent data from volume mounts on containers")
}

func (s *KoolStop) stop(purge bool) (err error) {
	var (
		args []string
	)

	args = []string{"down"}

	if purge {
		args = append(args, "--volumes", "--remove-orphans")
	}

	err = shell.Interactive("docker-compose", args...)
	return
}
