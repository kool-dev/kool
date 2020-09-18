package builder

import (
	"fmt"
	"kool-dev/kool/cmd/shell"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
)

// DefaultCommand holds data and logic for an executable command.
type DefaultCommand struct {
	command string
	args    []string
}

// Builder holds available methods for building commands.
type Builder interface {
	AppendArgs(args ...string)
	String() string
}

// Runner holds available methods for running commands.
type Runner interface {
	Interactive(args ...string) error
	Exec(args ...string) (string, error)
	LookPath() error
}

// NewCommand Create a new command.
func NewCommand(command string, args ...string) *DefaultCommand {
	return &DefaultCommand{command, args}
}

// ParseCommand transforms a command line string into separated
// command name and arguments list, expanding environment variables
// if any.
func ParseCommand(line string) (command *DefaultCommand, err error) {
	var parsed []string

	if parsed, err = shlex.Split(os.ExpandEnv(line)); err != nil {
		return
	}

	command = &DefaultCommand{parsed[0], parsed[1:]}

	return
}

// AppendArgs allows to appending arguments to the command builder.
func (c *DefaultCommand) AppendArgs(args ...string) {
	c.args = append(c.args, args...)
}

// String returns a string representation of the command.
func (c *DefaultCommand) String() string {
	return strings.Trim(fmt.Sprintf("%s %s", c.command, strings.Join(c.args, " ")), " ")
}

// LookPath returns if the command exists
func (c *DefaultCommand) LookPath() (err error) {
	_, err = exec.LookPath(c.command)
	return
}

// Interactive will send the command to an interactive execution.
func (c *DefaultCommand) Interactive(args ...string) (err error) {
	var finalArgs []string = c.args

	if len(args) > 0 {
		finalArgs = append(finalArgs, args...)
	}

	err = shell.Interactive(c.command, finalArgs...)
	return
}

// Exec will send the command to shell execution.
func (c *DefaultCommand) Exec(args ...string) (outStr string, err error) {
	var finalArgs []string = c.args

	if len(args) > 0 {
		finalArgs = append(finalArgs, args...)
	}

	outStr, err = shell.Exec(c.command, finalArgs...)
	return
}
