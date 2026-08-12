package commands

import (
	"errors"
	"fmt"
	"io"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/parser"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/api"
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
	AddKoolCloud(root)
	AddKoolDocker(root)
	AddKoolExec(root)
	AddKoolInfo(root)
	AddKoolLogs(root)
	AddKoolPreset(root)
	AddKoolProxy(root)
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

var rootCmd = newRootCmd(environment.NewEnvStorage(), true)

var originalWorkingDir = ""

func init() {
	AddCommands(rootCmd)

	// pass in the version
	api.SetCliVersion(version)
}

// NewRootCmd creates the root command
func NewRootCmd(env environment.EnvStorage) (cmd *cobra.Command) {
	return newRootCmd(env, false)
}

func newRootCmd(env environment.EnvStorage, initializeEnvironment bool) (cmd *cobra.Command) {
	environmentInitialized := false
	cmd = &cobra.Command{
		Args:          cobra.ArbitraryArgs,
		Use:           "kool",
		SilenceErrors: true,
		Short:         "Cloud native environments made easy",
		Long: `From development to production, a robust and easy-to-use developer tool
that makes Docker container adoption quick and easy for building and deploying cloud native
applications.

Complete documentation is available at https://kool.dev/docs`,
		Version:               version,
		DisableAutoGenTag:     true,
		DisableFlagsInUseLine: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			cmd.SilenceUsage = true
			workDirFlag := cmd.Flags().Lookup("working_dir")
			if workDirFlag != nil && workDirFlag.Value.String() != "" {
				workDir := workDirFlag.Value.String()

				currentWorkingDir, getwdErr := os.Getwd()
				if getwdErr != nil {
					return getwdErr
				}
				if originalWorkingDir == "" {
					originalWorkingDir = currentWorkingDir
				}
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
				if workDir, err = os.Getwd(); err != nil {
					return
				}
				env.Set("PWD", workDir)

			}
			if initializeEnvironment && !environmentInitialized {
				environment.InitEnvironmentVariables(env)
				environmentInitialized = true
			}
			if workspaceError := env.Get("KOOL_WORKSPACE_ERROR"); workspaceError != "" {
				return errors.New(workspaceError)
			}

			if verbose := cmd.Flags().Lookup("verbose"); verbose != nil && verbose.Value.String() == "true" {
				env.Set("KOOL_VERBOSE", verbose.Value.String())
			}

			if !hasWarnedDevelopmentVersion && version == DEV_VERSION && shell.NewTerminalChecker().IsTerminal(cmd.OutOrStdout()) {
				shell.NewShell().Warning("Warning: you are executing a development version of kool.")
				hasWarnedDevelopmentVersion = true
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

				err = fmt.Errorf("command not found")
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
		currentDirectory, getwdErr := os.Getwd()
		if getwdErr != nil {
			return getwdErr
		}
		currentPWD := os.Getenv("PWD")
		currentOriginalWorkingDir := originalWorkingDir
		defer func() {
			_ = os.Chdir(currentDirectory)
			_ = os.Setenv("PWD", currentPWD)
			originalWorkingDir = currentOriginalWorkingDir
		}()

		initializeEnvironment := hasWorkingDirArg(args)
		if initializeEnvironment {
			restoreEnvironment := clearDirectoryEnvironment(currentDirectory)
			defer restoreEnvironment()
			restoreWorkspace := environment.IsolateWorkspace()
			defer restoreWorkspace()
			originalWorkingDir = currentDirectory
		}
		childRoot := newRootCmd(environment.NewEnvStorage(), initializeEnvironment)

		childRoot.SetArgs(args)

		childRoot.SetIn(in)
		childRoot.SetOut(out)
		childRoot.SetErr(err)

		AddCommands(childRoot)

		return childRoot.Execute()
	}
}

func hasWorkingDirArg(args []string) bool {
	for _, arg := range args {
		if arg == "-w" || arg == "--working_dir" || strings.HasPrefix(arg, "--working_dir=") {
			return true
		}
	}
	return false
}

func clearDirectoryEnvironment(directory string) func() {
	previous := make(map[string]string)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		previous[parts[0]] = parts[1]
	}
	keys := []string{
		"PWD", "KOOL_NAME", "KOOL_WORKSPACES_ENABLED", "KOOL_PROXY_ENABLED",
		"KOOL_WORKSPACE", "KOOL_WORKSPACE_PROVIDER", "KOOL_WORKSPACE_ERROR",
		"KOOL_WORKSPACE_SOURCE", "KOOL_WORKSPACE_SOURCE_PROJECT", "KOOL_WORKSPACE_NAME",
		"KOOL_WORKSPACE_PATH", "KOOL_WORKSPACE_PROJECT", "KOOL_WORKSPACE_SERVICES",
		"KOOL_PROXY_DOMAIN", "KOOL_PROXY_HOST", "COMPOSE_PROJECT_NAME", "COMPOSE_FILE",
	}
	keys = append(keys, environment.LoadedEnvKeys(directory)...)
	for _, key := range keys {
		_ = os.Unsetenv(key)
	}
	return func() {
		for _, entry := range os.Environ() {
			key := strings.SplitN(entry, "=", 2)[0]
			if _, exists := previous[key]; !exists {
				_ = os.Unsetenv(key)
			}
		}
		for key, value := range previous {
			_ = os.Setenv(key, value)
		}
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
