package commands

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/services/checker"
	"kool-dev/kool/services/compose"

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

func AddKoolStop(root *cobra.Command) {
	var (
		stop    = NewKoolStop()
		stopCmd = NewStopCommand(stop)
	)

	root.AddCommand(stopCmd)
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
		s.down.AppendArgs("--remove-orphans")

		// no specific services passed in, so we gonna 'docker-compose down'
		if s.Flags.Purge {
			s.down.AppendArgs("--volumes")
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
		Use:   "stop [SERVICE...]",
		Short: "Stop and destroy running service containers",
		Long: `Stop and destroy the specified [SERVICE] containers, which were started
using 'kool start'. If no [SERVICE] is provided, all running containers are stopped.`,
		Run: DefaultCommandRunFunction(stop),

		DisableFlagsInUseLine: true,
	}

	stopCmd.Flags().BoolVarP(&stop.Flags.Purge, "purge", "", false, "Remove all persistent data from volume mounts on containers")
	return
}
