package cmd

import (
	"github.com/spf13/cobra"
	"kool-dev/kool/environment"
)

// CobraRunFN Cobra command run function
type CobraRunFN func(*cobra.Command, []string)

var version string = "0.0.0-dev"

var rootCmd = NewRootCmd(environment.NewEnvStorage())

// NewRootCmd creates the root command
func NewRootCmd(envStorage environment.EnvStorage) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "kool",
		Short: "kool - Kool stuff",
		Long: `An easy and robust software development environment
tool helping you from project creation until deployment.
Complete documentation is available at https://kool.dev/docs`,
		Version:           version,
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmf *cobra.Command, args []string) {
			if verbose := cmf.Flags().Lookup("verbose"); verbose != nil && verbose.Value.String() == "true" {
				envStorage.Set("KOOL_VERBOSE", verbose.Value.String())
			}
		},
	}

	cmd.PersistentFlags().Bool("verbose", false, "execute commands with verbosity")
	return
}

// Execute proxies the call to cobra root command
func Execute() error {
	return rootCmd.Execute()
}

// RootCmd exposes the root command
func RootCmd() *cobra.Command {
	return rootCmd
}

// DefaultCommandRunFunction default run function logic
func DefaultCommandRunFunction(services ...KoolService) CobraRunFN {
	return func(cmd *cobra.Command, args []string) {
		for _, service := range services {
			service.SetWriter(cmd.OutOrStdout())
			service.SetReader(cmd.InOrStdin())

			if err := service.Execute(args); err != nil {
				service.Error(err)
				service.Exit(1)
			}
		}
	}
}
