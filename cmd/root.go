package cmd

import (
	"io"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"

	"github.com/spf13/cobra"
)

// CobraRunFN Cobra command run function
type CobraRunFN func(*cobra.Command, []string)

// AddCommandsFN function to add subcommands
type AddCommandsFN func(*cobra.Command)

var hasWarnedDevelopmentVersion = false

var AddCommands AddCommandsFN = func(root *cobra.Command) {
	AddKoolCompletion(root)
	AddKoolCreate(root)
	AddKoolDeploy(root)
	AddKoolDocker(root)
	AddKoolExec(root)
	AddKoolInfo(root)
	AddKoolInit(root)
	AddKoolLogs(root)
	AddKoolPreset(root)
	AddKoolRestart(root)
	AddKoolRun(root)
	AddKoolSelfUpdate(root)
	AddKoolShare(root)
	AddKoolStart(root)
	AddKoolStatus(root)
	AddKoolStop(root)
}

// DEV_VERSION holds the static version shown for development time builds
const DEV_VERSION = "0.0.0-dev"

var version string = DEV_VERSION

var rootCmd = NewRootCmd(environment.NewEnvStorage())

func init() {
	AddCommands(rootCmd)
}

// NewRootCmd creates the root command
func NewRootCmd(env environment.EnvStorage) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "kool",
		Short: "kool - Kool stuff",
		Long: `An easy and robust software development environment
tool helping you from project creation through deployment.
Complete documentation is available at https://kool.dev/docs`,
		Version:           version,
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose := cmd.Flags().Lookup("verbose"); verbose != nil && verbose.Value.String() == "true" {
				env.Set("KOOL_VERBOSE", verbose.Value.String())
			}

			if !hasWarnedDevelopmentVersion && version == DEV_VERSION && shell.NewTerminalChecker().IsTerminal(cmd.OutOrStdout()) {
				shell.NewShell().Warning("Warning: you are executing a development version of kool.")
				hasWarnedDevelopmentVersion = true
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
		childRoot := NewRootCmd(environment.NewEnvStorage())

		childRoot.SetArgs(args)

		childRoot.SetIn(in)
		childRoot.SetOut(out)
		childRoot.SetErr(err)

		AddCommands(childRoot)

		return childRoot.Execute()
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
