package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
	"kool-dev/kool/cmd/shell"
	"os"

	"github.com/spf13/cobra"
)

// KoolStart holds handlers and functions for starting containers logic
type KoolStart struct {
	dependenciesChecker   checker.Checker
	networkHandler        network.Handler
	startContainersRunner builder.Runner

	exiter shell.Exiter
	out    shell.OutputWriter
}

// NewStartCommand initializes new kool start command
func NewStartCommand(start *KoolStart) *cobra.Command {
	return &cobra.Command{
		Use:   "start [SERVICE]",
		Short: "Start the specified Kool environment containers. If no service is specified, start all.",
		Run: func(cmd *cobra.Command, args []string) {
			start.out.SetWriter(cmd.OutOrStdout())

			if err := start.Execute(args); err != nil {
				start.out.Error(err)
				start.exiter.Exit(1)
			}
		},
		DisableFlagsInUseLine: true,
	}
}

// NewKoolStart creates a new pointer with default KoolStart service
// dependencies.
func NewKoolStart() *KoolStart {
	return &KoolStart{
		checker.NewChecker(),
		network.NewHandler(),
		builder.NewCommand("docker-compose", "up", "-d", "--force-recreate"),
		shell.NewExiter(),
		shell.NewOutputWriter(),
	}
}

var startCmd = NewStartCommand(NewKoolStart())

func init() {
	rootCmd.AddCommand(startCmd)
}

// Execute runs the start logic with incoming arguments.
func (s *KoolStart) Execute(args []string) (err error) {
	if err = s.checkDependencies(); err != nil {
		return
	}

	err = s.startContainers(args)
	return
}

func (s *KoolStart) checkDependencies() (err error) {
	if err = s.dependenciesChecker.Check(); err != nil {
		return
	}

	if err = s.networkHandler.HandleGlobalNetwork(os.Getenv("KOOL_GLOBAL_NETWORK")); err != nil {
		return
	}

	return
}

func (s *KoolStart) startContainers(services []string) (err error) {
	err = s.startContainersRunner.Interactive(services...)
	return
}
