package commands

import (
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/services/compose"
	"os"
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

func (e *KoolExec) detectTTY() {
	var isTerminal = e.Shell().IsTerminal()

	if !isTerminal {
		e.composeExec.AppendArgs("-T")
	}

	if aware, ok := e.composeExec.(compose.TtyAware); ok {
		// let DockerCompose know about whether we are under TTY or not
		aware.SetIsTTY(isTerminal)
	}
}

func (e *KoolExec) checkUser(service string) {
	var (
		asuser      string
		err         error
		actualInput io.Reader
	)

	if asuser = e.env.Get("KOOL_ASUSER"); asuser == "" {
		return
	}

	actualInput = e.Shell().InStream()
	defer e.Shell().SetInStream(actualInput) // return actualInput

	// avoid interference of Exec on actual input
	// by temporarily setting os.Stdin
	e.Shell().SetInStream(os.Stdin)

	// we have a KOOL_ASUSER env; now we need to know whether
	// the image of the target service have such user
	var passwd string

	if passwd, err = e.Shell().Exec(e.composeExec, service, "cat", "/etc/passwd"); err != nil {
		// for safety, let's write the warning message to os.Stderr
		// so we avoid getting cross-fire on in/out redirections
		actualOut := e.Shell().OutStream()
		defer e.Shell().SetOutStream(actualOut)

		e.Shell().SetOutStream(os.Stderr)
		e.Shell().Warning(fmt.Sprintf("failed to check running container for kool user; did you forget kool start? (err: %s)", err.Error()))
	} else if strings.Contains(passwd, fmt.Sprintf("kool:x:%s", asuser)) {
		// since user (kool:x:UID) exists within the container, we set it
		e.composeExec.AppendArgs("--user", asuser)
	}
}

// Execute runs the exec logic with incoming arguments.
func (e *KoolExec) Execute(args []string) (err error) {
	e.detectTTY()

	e.checkUser(args[0])

	if len(e.Flags.EnvVariables) > 0 {
		for _, envVar := range e.Flags.EnvVariables {
			e.composeExec.AppendArgs("--env", envVar)
		}
	}

	if e.Flags.Detach {
		e.composeExec.AppendArgs("--detach")
	}

	err = e.Shell().Interactive(e.composeExec, args...)
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
