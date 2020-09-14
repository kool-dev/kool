package builder

import (
	"fmt"
	"kool-dev/kool/cmd/shell"
	"os"
	"strings"
	"os/exec"

	"github.com/google/shlex"
)

// Command holds data and logic for an executable command.
type Command struct {
	command string
	args    []string
}

// Executor holds available methods for command.
type Executor interface {
	AppendArgs(args ...string)
	Interactive() error
	String() string
	Exec() (string, error)
	LookPath() error
}

// ParseCommand transforms a command line string into separated
// command name and arguments list, expanding environment variables
// if any.
func ParseCommand(line string) (command *Command, err error) {
	var parsed []string

	if parsed, err = shlex.Split(os.ExpandEnv(line)); err != nil {
		return
	}

	command = &Command{parsed[0], parsed[1:]}

	return
}

// AppendArgs allows to appending arguments to the command builder.
func (c *Command) AppendArgs(args ...string) {
	c.args = append(c.args, args...)
}

// Interactive will send the command to an interactive execution.
func (c *Command) Interactive() (err error) {
	err = shell.Interactive(c.command, c.args...)
	return
}

// Exec will send the command to shell execution.
func (c *Command) Exec() (outStr string, err error) {
	outStr, err = shell.Exec(c.command, c.args...)
	return
}

// String returns a string representation of the command.
func (c *Command) String() string {
	return strings.Trim(fmt.Sprintf("%s %s", c.command, strings.Join(c.args, " ")), " ")
}

// LookPath returns if the command exists
func (c *Command) LookPath() (err error) {
	_, err = exec.LookPath(c.command)
	return
}
