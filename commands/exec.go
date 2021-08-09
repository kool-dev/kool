package commands

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/services/compose"
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

func AddKoolExec(root *cobra.Command) {
	var (
		exec    = NewKoolExec()
		execCmd = NewExecCommand(exec)
	)

	root.AddCommand(execCmd)
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
	if !e.IsTerminal() {
		e.composeExec.AppendArgs("-T")
	}

	if asuser := e.env.Get("KOOL_ASUSER"); asuser != "" {
		// we have a KOOL_ASUSER env; now we need to know whether
		// the image of the target service have such user
		var passwd string
		if passwd, err = e.Exec(e.composeExec, args[0], "cat", "/etc/passwd"); err != nil {
			e.Warning("failed to check running container for kool user; not setting a user (err: %s)", err.Error())
			err = nil
		} else if strings.Contains(passwd, fmt.Sprintf("kool:x:%s", asuser)) {
			// since user (kool:x:UID) exists within the container, we set it
			e.composeExec.AppendArgs("--user", asuser)
		}
	}

	if aware, ok := e.composeExec.(compose.TtyAware); ok {
		// let DockerCompose know about whether we are under TTY or not
		aware.SetIsTTY(e.IsTerminal())
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
		Use:   "exec [OPTIONS] SERVICE COMMAND [--] [ARG...]",
		Short: "Execute a command inside a running service container",
		Long:  `Execute a COMMAND inside the specified SERVICE container (similar to an SSH session).`,
		Args:  cobra.MinimumNArgs(2),
		RunE:  DefaultCommandRunFunction(exec),

		DisableFlagsInUseLine: true,
	}

	execCmd.Flags().BoolVarP(&exec.Flags.DisableTty, "disable-tty", "T", false, "Deprecated - no effect.")
	execCmd.Flags().StringArrayVarP(&exec.Flags.EnvVariables, "env", "e", []string{}, "Environment variables.")
	execCmd.Flags().BoolVarP(&exec.Flags.Detach, "detach", "d", false, "Detached mode: Run command in the background.")

	//After a non-flag arg, stop parsing flags
	execCmd.Flags().SetInterspersed(false)
	return
}
