package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
	"os"

	"github.com/spf13/cobra"
)

// KoolStart holds handlers and functions for starting containers logic
type KoolStart struct {
	DefaultKoolService

	check checker.Checker
	net   network.Handler
	start builder.Runner
}

// NewStartCommand initializes new kool start command
func NewStartCommand(start *KoolStart) *cobra.Command {
	return &cobra.Command{
		Use:                   "start [SERVICE]",
		Short:                 "Start the specified Kool environment containers. If no service is specified, start all.",
		Run:                   DefaultCommandRunFunction(start),
		DisableFlagsInUseLine: true,
	}
}

// NewKoolStart creates a new pointer with default KoolStart service
// dependencies.
func NewKoolStart() *KoolStart {
	return &KoolStart{
		*newDefaultKoolService(),
		checker.NewChecker(),
		network.NewHandler(),
		builder.NewCommand("docker-compose", "up", "-d", "--force-recreate"),
	}
}

func init() {
	rootCmd.AddCommand(NewStartCommand(NewKoolStart()))
}

// Execute runs the start logic with incoming arguments.
func (s *KoolStart) Execute(args []string) (err error) {
	if err = s.check.Check(); err != nil {
		return
	}

	if err = s.net.HandleGlobalNetwork(os.Getenv("KOOL_GLOBAL_NETWORK")); err != nil {
		return
	}

	err = s.start.Interactive(args...)
	return
}
