package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"

	"github.com/spf13/cobra"
)

// KoolExecFlags holds the flags for the exec command
type KoolExecFlags struct {
	DisableTty   bool
	EnvVariables []string
	Detach       bool
}

// KoolExec holds handlers and functions to implement the exec command logic
type KoolExec struct {
	DefaultKoolService
	Flags *KoolExecFlags

	terminal    shell.TerminalChecker
	envStorage  environment.EnvStorage
	composeExec builder.Command
}

func init() {
	var (
		exec    = NewKoolExec()
		execCmd = NewExecCommand(exec)
	)

	rootCmd.AddCommand(execCmd)
}

// NewKoolExec creates a new handler for exec logic
func NewKoolExec() *KoolExec {
	return &KoolExec{
		*newDefaultKoolService(),
		&KoolExecFlags{false, []string{}, false},
		shell.NewTerminalChecker(),
		environment.NewEnvStorage(),
		builder.NewCommand("docker-compose", "exec"),
	}
}

// Execute runs the exec logic with incoming arguments.
func (e *KoolExec) Execute(args []string) (err error) {
	if !e.terminal.IsTerminal(e.GetWriter()) {
		e.composeExec.AppendArgs("-T")
	}

	if asuser := e.envStorage.Get("KOOL_ASUSER"); asuser != "" {
		e.composeExec.AppendArgs("--user", asuser)
	}

	if len(e.Flags.EnvVariables) > 0 {
		for _, envVar := range e.Flags.EnvVariables {
			e.composeExec.AppendArgs("--env", envVar)
		}
	}

	if e.Flags.Detach {
		e.composeExec.AppendArgs("--detach")
	}

	err = e.composeExec.Interactive(args...)
	return
}

// NewExecCommand initializes new kool exec command
func NewExecCommand(exec *KoolExec) (execCmd *cobra.Command) {
	execCmd = &cobra.Command{
		Use:   "exec [options] [service] [command]",
		Short: "Execute a command within a running service container",
		Args:  cobra.MinimumNArgs(2),
		Run:   DefaultCommandRunFunction(exec),
	}

	execCmd.Flags().BoolVarP(&exec.Flags.DisableTty, "disable-tty", "T", false, "Deprecated - no effect")
	execCmd.Flags().StringArrayVarP(&exec.Flags.EnvVariables, "env", "e", []string{}, "Environment variables")
	execCmd.Flags().BoolVarP(&exec.Flags.Detach, "detach", "d", false, "Detached mode: Run command in the background")

	//After a non-flag arg, stop parsing flags
	execCmd.Flags().SetInterspersed(false)
	return
}
