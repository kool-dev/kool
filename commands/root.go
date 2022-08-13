package commands

import (
	"fmt"
	"io"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/parser"
	"kool-dev/kool/core/shell"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// CobraRunE Cobra command run function
type CobraRunE func(*cobra.Command, []string) error

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
	AddKoolLogs(root)
	AddKoolPreset(root)
	AddKoolRestart(root)
	AddKoolRun(root)
	AddKoolSelfUpdate(root)
	AddKoolShare(root)
	AddKoolStart(root)
	AddKoolStatus(root)
	AddKoolStop(root)
	AddKoolRecipe(root)
}

// DEV_VERSION holds the static version shown for development time builds
const DEV_VERSION = "0.0.0-dev"

var version string = DEV_VERSION

var rootCmd = NewRootCmd(environment.NewEnvStorage())

var originalWorkingDir = ""

func init() {
	AddCommands(rootCmd)
}

// NewRootCmd creates the root command
func NewRootCmd(env environment.EnvStorage) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Args:          cobra.ArbitraryArgs,
		Use:           "kool",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "Cloud native environments made easy",
		Long: `From development to production, a robust and easy-to-use developer tool
that makes Docker container adoption quick and easy for building and deploying cloud native
applications.

Complete documentation is available at https://kool.dev/docs`,
		Version:               version,
		DisableAutoGenTag:     true,
		DisableFlagsInUseLine: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if verbose := cmd.Flags().Lookup("verbose"); verbose != nil && verbose.Value.String() == "true" {
				env.Set("KOOL_VERBOSE", verbose.Value.String())
			}

			if !hasWarnedDevelopmentVersion && version == DEV_VERSION && shell.NewTerminalChecker().IsTerminal(cmd.OutOrStdout()) {
				shell.NewShell().Warning("Warning: you are executing a development version of kool.")
				hasWarnedDevelopmentVersion = true
			}

			workDirFlag := cmd.Flags().Lookup("working_dir")
			if workDirFlag != nil && workDirFlag.Value.String() != "" {
				workDir := workDirFlag.Value.String()

				if originalWorkingDir != "" {
					// having an original working dir set means we have
					// already changed the working dir before and we are in
					//  a recursive kool call. We need to restore the original
					// working dir before changing it again.
					if err = os.Chdir(originalWorkingDir); err != nil {
						return
					}
				}

				if !path.IsAbs(workDir) {
					if workDir, err = filepath.Abs(workDir); err != nil {
						return
					}
				}

				if err = os.Chdir(workDir); err != nil {
					return
				}

				if originalWorkingDir == "" {
					// we only set the original working dir if it is not set
					// yet. This is to avoid overriding the original working
					// dir in recursive calls.
					if originalWorkingDir, err = os.Getwd(); err != nil {
						return
					}
				}

				environment.NewEnvStorage().Set("PWD", workDir)
			}

			return
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 {
				err = cmd.Help()
				return
			}

			scriptParser := parser.NewParser()
			env := environment.NewEnvStorage()

			// look for kool.yml on current working directory
			_ = scriptParser.AddLookupPath(env.Get("PWD"))
			// look for kool.yml on kool folder within user home directory
			_ = scriptParser.AddLookupPath(path.Join(env.Get("HOME"), "kool"))

			_, err = scriptParser.Parse(args[0])

			if err == nil {
				// we did find a script for it!
				newDefaultKoolService().Shell().Info(fmt.Sprintf("Did you mean 'kool run %s'?", strings.Join(args, " ")))

				err = fmt.Errorf("Command not found.")
				return
			}

			err = fmt.Errorf("command %s not found", args[0])
			return
		},
	}

	cmd.PersistentFlags().Bool("verbose", false, "Increases output verbosity")
	cmd.PersistentFlags().StringP("working_dir", "w", "", "Changes the working directory for the command")
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
func DefaultCommandRunFunction(services ...KoolService) CobraRunE {
	return func(cmd *cobra.Command, args []string) (err error) {
		for _, service := range services {
			service.Shell().SetOutStream(cmd.OutOrStdout())
			service.Shell().SetInStream(cmd.InOrStdin())
			service.Shell().SetErrStream(cmd.ErrOrStderr())

			if err = service.Execute(args); err != nil {
				if shell.IsUserCancelledError(err) {
					service.Shell().Warning("Operation Cancelled")
					err = nil
				}
				return
			}
		}
		return
	}
}

// LongTaskCommandRunFunction long tasks run function logic
func LongTaskCommandRunFunction(tasks ...KoolTask) CobraRunE {
	return func(cmd *cobra.Command, args []string) (err error) {
		for _, task := range tasks {
			task.Shell().SetOutStream(cmd.OutOrStdout())
			task.Shell().SetInStream(cmd.InOrStdin())
			task.Shell().SetErrStream(cmd.ErrOrStderr())

			if err = task.Run(args); err != nil {
				return
			}
		}
		return
	}
}
