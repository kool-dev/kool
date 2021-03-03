package cmd

import (
	"io"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"

	"github.com/spf13/cobra"
)

// CobraRunFN Cobra command run function
type CobraRunFN func(*cobra.Command, []string)

var version string = "0.0.0-dev"

var rootCmd = NewRootCmd(environment.NewEnvStorage())

// NewRootCmd creates the root command
func NewRootCmd(env environment.EnvStorage) (cmd *cobra.Command) {
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
				env.Set("KOOL_VERBOSE", verbose.Value.String())
			}
		},
	}

	cmd.PersistentFlags().Bool("verbose", false, "increases output verbosity")
	return
}

// Execute proxies the call to cobra root command
func Execute() error {
	setRecursiveCall(rootCmd)

	return rootCmd.Execute()
}

func setRecursiveCall(root *cobra.Command) {
	shell.RecursiveCall = func(args []string, in io.Reader, out, err io.Writer) error {
		root.SetArgs(args)

		oldIn := root.InOrStdin()
		oldOut := root.OutOrStdout()
		oldErr := root.ErrOrStderr()

		root.SetIn(in)
		root.SetOut(out)
		root.SetErr(err)

		defer func() {
			root.SetIn(oldIn)
			root.SetOut(oldOut)
			root.SetErr(oldErr)
		}()

		return root.Execute()
	}
}

// RootCmd exposes the root command
func RootCmd() *cobra.Command {
	return rootCmd
}

// DefaultCommandRunFunction default run function logic
func DefaultCommandRunFunction(services ...KoolService) CobraRunFN {
	return func(cmd *cobra.Command, args []string) {
		for _, service := range services {
			service.SetOutStream(cmd.OutOrStdout())
			service.SetInStream(cmd.InOrStdin())
			service.SetErrStream(cmd.ErrOrStderr())

			if err := service.Execute(args); err != nil {
				service.Error(err)
				service.Exit(1)
			}
		}
	}
}

// LongTaskCommandRunFunction long tasks run function logic
func LongTaskCommandRunFunction(tasks ...KoolTask) CobraRunFN {
	return func(cmd *cobra.Command, args []string) {
		for _, task := range tasks {
			task.SetOutStream(cmd.OutOrStdout())
			task.SetInStream(cmd.InOrStdin())
			task.SetErrStream(cmd.ErrOrStderr())

			if err := task.Run(args); err != nil {
				task.Error(err)
				task.Exit(1)
			}
		}
	}
}
