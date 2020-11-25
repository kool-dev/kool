package builder

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/shlex"
)

// DefaultCommand holds data and logic for an executable command.
type DefaultCommand struct {
	command string
	args    []string
}

// Command holds available methods for building commands.
type Command interface {
	Parser
	AppendArgs(...string)
	String() string
	Args() []string
	Cmd() string
}

// Parser holds available methods for parse commands
type Parser interface {
	Parse(string) error
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

// Args returns the command arguments
func (c *DefaultCommand) Args() []string {
	return c.args
}

// Cmd returns the command executable
func (c *DefaultCommand) Cmd() string {
	return c.command
}

// Parse calls the ParseCommand function
func (c *DefaultCommand) Parse(line string) (err error) {
	if parsed, err := ParseCommand(line); err == nil {
		*c = *parsed
	}

	return
}
