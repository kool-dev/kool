package shell

import (
	"io"
	"kool-dev/kool/core/builder"
	"os"
	"os/exec"
)

// InputRedirect holds the key to indicate the right part
// of the command expression is meant to be the Stdin
// for the whole left part.
const InputRedirect string = "<"

// OutputRedirect holds the key to indicate the output from
// the left part of the command up to this key is meant to be
// written to the destiny pointed by the right part.
const OutputRedirect string = ">"

// OutputRedirectAppend holds the key to indicate the output from
// the left part of the command up to this key is meant to be
// written in append mode to the destiny pointed by the right part.
const OutputRedirectAppend string = ">>"

// OutputPipe holds the key to indicate the output from
// the left command up to this key is meant to be
// written as input the command in the right part of it.
const OutputPipe string = "|"

// CommandWithPointers holds the set of file descriptors to be used for
// executing the given command.
type CommandWithPointers struct {
	builder.Command

	in       io.Reader
	out, err io.Writer

	hasCustomStdin  bool
	hasCustomStdout bool
}

// ParsedRedirect holds logic for parsed redirect
type ParsedRedirect interface {
	Close()
}

// Close closes reader and writer if necessary
func (c *CommandWithPointers) Close() {
	if c.hasCustomStdin {
		if cl, ok := c.in.(io.Closer); ok {
			cl.Close()
		}
	}
	if c.hasCustomStdout {
		if cl, ok := c.out.(io.WriteCloser); ok {
			cl.Close()
		}
	}
}

// Cmd creates a new *exec.Command for given command
func (c *CommandWithPointers) Cmd() (cmd *exec.Cmd) {
	cmd = execCmdFn(c.Command.Cmd(), c.Command.Args()...)
	cmd.Env = os.Environ()
	cmd.Stdout = c.out
	cmd.Stderr = c.err
	cmd.Stdin = c.in
	return
}

func hasRedirect(command builder.Command) bool {
	var count int = len(command.Args())

	if count < 2 {
		return false
	}

	check := command.Args()[count-2]

	return check == InputRedirect || check == OutputRedirect || check == OutputRedirectAppend
}

func parseRedirects(command builder.Command, sh Shell) (cmdptr *CommandWithPointers, err error) {
	if !hasRedirect(command) {
		// lastly command which does not have a redirect
		cmdptr = &CommandWithPointers{
			Command: command,
			in:      sh.InStream(),
			out:     sh.OutStream(),
			err:     sh.ErrStream(),
		}
		return
	}

	if cmdptr, err = splitRedirect(command); err != nil {
		return
	}

	var chainedCmdptr *CommandWithPointers
	// process chained redirects like "cmd < input > output"
	if chainedCmdptr, err = parseRedirects(cmdptr.Command, sh); err != nil {
		return
	}

	if cmdptr.hasCustomStdin {
		chainedCmdptr.in = cmdptr.in
		chainedCmdptr.hasCustomStdin = cmdptr.hasCustomStdin
	}
	if cmdptr.hasCustomStdout {
		chainedCmdptr.out = cmdptr.out
		chainedCmdptr.hasCustomStdout = cmdptr.hasCustomStdout
	}

	cmdptr = chainedCmdptr

	return
}

func splitRedirect(cmd builder.Command) (cmdptr *CommandWithPointers, err error) {
	var (
		args    = cmd.Args()
		numArgs = len(args)
		inFile  io.ReadCloser
		outFile io.WriteCloser
	)

	// check the before-last position of the command
	// for some redirect key and properly handle them.

	switch args[numArgs-2] {
	case InputRedirect:
		{
			if inFile, err = os.OpenFile(args[numArgs-1], os.O_RDONLY, os.ModePerm); err != nil {
				return
			}
		}
	case OutputRedirect, OutputRedirectAppend:
		{
			var mode int = os.O_CREATE | os.O_WRONLY
			if args[numArgs-2] == OutputRedirectAppend {
				mode |= os.O_APPEND
			} else {
				mode |= os.O_TRUNC
			}

			if outFile, err = os.OpenFile(args[numArgs-1], mode, os.ModePerm); err != nil {
				return
			}
		}
	}

	cmdptr = &CommandWithPointers{
		Command: builder.NewCommand(cmd.Cmd(), args[:numArgs-2]...),
	}

	if inFile != nil {
		cmdptr.in = inFile
		cmdptr.hasCustomStdin = true
	}
	if outFile != nil {
		cmdptr.out = outFile
		cmdptr.hasCustomStdout = true
	}

	return
}
