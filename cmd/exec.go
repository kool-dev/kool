package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/compose"
	"kool-dev/kool/environment"
	"strings"

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

	env         environment.EnvStorage
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
		environment.NewEnvStorage(),
		compose.NewDockerCompose("exec"),
	}
}

// Execute runs the exec logic with incoming arguments.
func (e *KoolExec) Execute(args []string) (err error) {
	e.composeExec.Reset()
	if asuser := e.env.Get("KOOL_ASUSER"); asuser != "" {
		// we have a KOOL_ASUSER env; now we need to know whether
		// the image of the target service have such user
		passwd, _ := e.Exec(e.composeExec, args[0], "cat", "/etc/passwd")
		// kool:x:UID
		if strings.Contains(passwd, fmt.Sprintf("kool:x:%s", asuser)) {
			// since user existing within the container, we use it
			e.composeExec.AppendArgs("--user", asuser)
		}
	}

	if !e.IsTerminal() {
		e.composeExec.AppendArgs("-T")
	}

	if _, assert := e.composeExec.(*compose.DockerCompose); assert {
		// let DockerCompose know about wheter we are under TTY or not
		e.composeExec.(*compose.DockerCompose).SetIsTTY(e.IsTerminal())
	}

	if len(e.Flags.EnvVariables) > 0 {
		for _, envVar := range e.Flags.EnvVariables {
			e.composeExec.AppendArgs("--env", envVar)
		}
	}

	if e.Flags.Detach {
		e.composeExec.AppendArgs("--detach")
	}

	err = e.Interactive(e.composeExec, args...)
	return
}

// NewExecCommand initializes new kool exec command
func NewExecCommand(exec *KoolExec) (execCmd *cobra.Command) {
	execCmd = &cobra.Command{
		Use:   "exec [options] [service] [command]",
		Short: "Execute a [command] inside the specified [service] container.",
		Args:  cobra.MinimumNArgs(2),
		Run:   DefaultCommandRunFunction(exec),
	}

	execCmd.Flags().BoolVarP(&exec.Flags.DisableTty, "disable-tty", "T", false, "Deprecated - no effect.")
	execCmd.Flags().StringArrayVarP(&exec.Flags.EnvVariables, "env", "e", []string{}, "Environment variables.")
	execCmd.Flags().BoolVarP(&exec.Flags.Detach, "detach", "d", false, "Detached mode: Run command in the background.")

	//After a non-flag arg, stop parsing flags
	execCmd.Flags().SetInterspersed(false)
	return
}
