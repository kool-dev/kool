package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/compose"

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

	check checker.Checker
	down  builder.Command
	rm    builder.Command
}

func init() {
	var (
		stop    = NewKoolStop()
		stopCmd = NewStopCommand(stop)
	)

	rootCmd.AddCommand(stopCmd)
}

// NewKoolStop creates a new handler for stop logic with default dependencies
func NewKoolStop() *KoolStop {
	defaultKoolService := newDefaultKoolService()
	return &KoolStop{
		*defaultKoolService,
		&KoolStopFlags{false},
		checker.NewChecker(defaultKoolService.shell),
		compose.NewDockerCompose("down"),
		compose.NewDockerCompose("rm"),
	}
}

// Execute runs the stop logic with incoming arguments.
func (s *KoolStop) Execute(args []string) (err error) {
	var stopCommand builder.Command

	if err = s.check.Check(); err != nil {
		return
	}

	if len(args) == 0 {
		// no specific services passed in, so we gonna 'docker-compose down'
		if s.Flags.Purge {
			s.down.AppendArgs("--volumes", "--remove-orphans")
		}

		stopCommand = s.down
	} else {
		// we should only stop some services!
		s.rm.AppendArgs("-s", "-f") // stops containers; no interactive
		if s.Flags.Purge {
			s.rm.AppendArgs("-v") // removes volumes

			s.Warning("Attention: when stopping specific services, only anonymous volumes will be removed.")
		}

		s.rm.AppendArgs(args...)

		stopCommand = s.rm
	}

	err = s.Interactive(stopCommand)
	return
}

// NewStopCommand initializes new kool stop command
func NewStopCommand(stop *KoolStop) (stopCmd *cobra.Command) {
	stopCmd = &cobra.Command{
		Use:   "stop [service...]",
		Short: "Stop and destroy running [service] containers started with the 'kool start' command. If no [service] is provided, all containers will be stopped.",
		Run:   DefaultCommandRunFunction(stop),
	}

	stopCmd.Flags().BoolVarP(&stop.Flags.Purge, "purge", "", false, "Remove all persistent data from volume mounts on containers")
	return
}
