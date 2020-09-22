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
	DependenciesChecker   checker.Checker
	NetworkHandler        network.Handler
	StartContainersRunner builder.Runner
	Exiter                shell.Exiter
}

var startCmdOutputWriter shell.OutputWriter = shell.NewOutputWriter()

// NewStartCommand initializes new kool start command
func NewStartCommand(startCmd *DefaultStartCmd) *cobra.Command {
	return &cobra.Command{
		Use:   "start [service]",
		Short: "Start the specified Kool environment containers. If no service is specified, start all.",
		Run: func(cmd *cobra.Command, args []string) {
			startCmdOutputWriter.SetWriter(cmd.OutOrStdout())

			if err := startCmd.checkDependencies(); err != nil {
				startCmdOutputWriter.ExecError("", err)
				startCmd.Exiter.Exit(1)
			}

			startCmd.startContainers(args)
		},
	}
}

func init() {
	defaultStartCmd := &DefaultStartCmd{
		checker.NewChecker(),
		network.NewHandler(),
		builder.NewCommand("docker-compose", "up", "-d", "--force-recreate"),
		shell.NewExiter(),
	}
	rootCmd.AddCommand(NewStartCommand(defaultStartCmd))
}

func (s *DefaultStartCmd) checkDependencies() (err error) {
	if err = s.DependenciesChecker.VerifyDependencies(); err != nil {
		return
	}

	if err = s.NetworkHandler.HandleGlobalNetwork(os.Getenv("KOOL_GLOBAL_NETWORK")); err != nil {
		return
	}

	return
}

func (s *DefaultStartCmd) startContainers(services []string) {
	err := s.StartContainersRunner.Interactive(services...)

	if err != nil {
		startCmdOutputWriter.ExecError("", err)
		s.Exiter.Exit(1)
	}
}
