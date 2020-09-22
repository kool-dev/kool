package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
	"kool-dev/kool/cmd/shell"
	"os"

	"github.com/spf13/cobra"
)

// DefaultStartCmd holds interfaces for status command logic
type DefaultStartCmd struct {
	dependenciesChecker   checker.Checker
	networkHandler        network.Handler
	startContainersRunner builder.Runner
	exiter                shell.Exiter

	out shell.OutputWriter
}

// NewStartCommand initializes new kool start command
func NewStartCommand(startCmd *DefaultStartCmd) *cobra.Command {
	return &cobra.Command{
		Use:   "start [SERVICE]",
		Short: "Start the specified Kool environment containers. If no service is specified, start all.",
		Run: func(cmd *cobra.Command, args []string) {
			startCmd.out.SetWriter(cmd.OutOrStdout())

			if err := startCmd.checkDependencies(); err != nil {
				startCmd.out.Error(err)
				startCmd.exiter.Exit(1)
			}

			startCmd.startContainers(args)
		},
	}
}

var startCmd = NewStartCommand(&DefaultStartCmd{
	checker.NewChecker(),
	network.NewHandler(),
	builder.NewCommand("docker-compose", "up", "-d", "--force-recreate"),
	shell.NewExiter(),
	shell.NewOutputWriter(),
})

func init() {
	rootCmd.AddCommand(startCmd)
}

func (s *DefaultStartCmd) checkDependencies() (err error) {
	if err = s.dependenciesChecker.Check(); err != nil {
		return
	}

	if err = s.networkHandler.HandleGlobalNetwork(os.Getenv("KOOL_GLOBAL_NETWORK")); err != nil {
		return
	}

	return
}

func (s *DefaultStartCmd) startContainers(services []string) {
	err := s.startContainersRunner.Interactive(services...)

	if err != nil {
		s.out.Error(err)
		s.exiter.Exit(1)
	}
}
