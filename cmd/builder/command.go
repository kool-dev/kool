package builder

import (
	"fmt"
	"kool-dev/kool/cmd/shell"
	"os"
	"strings"

	"github.com/google/shlex"
)

// Command holds data and logic for an executable command.
type Command struct {
	command string
	args    []string
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

// String returns a string representation of the command.
func (c *Command) String() string {
	return strings.Trim(fmt.Sprintf("%s %s", c.command, strings.Join(c.args, " ")), " ")
}
