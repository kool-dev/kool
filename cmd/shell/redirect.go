package shell

import (
	"io"
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

// DefaultParsedRedirect holds parsed redirect data
type DefaultParsedRedirect struct {
	args        []string
	shell       Shell
	closeStdin  bool
	closeStdout bool
}

// ParsedRedirect holds logic for parsed redirect
type ParsedRedirect interface {
	Close()
}

// Close closes reader and writer if necessary
func (p *DefaultParsedRedirect) Close() {
	if p.closeStdin {
		p.shell.InStream().(*os.File).Close()
	}
	if p.closeStdout {
		p.shell.OutStream().(*os.File).Close()
	}
}

// CreateCommand creates a new *exec.Command for given executable
func (p *DefaultParsedRedirect) CreateCommand(executable string) (cmd *exec.Cmd) {
	cmd = exec.Command(executable, p.args...)
	cmd.Env = os.Environ()
	cmd.Stdout = p.shell.OutStream()
	cmd.Stderr = p.shell.ErrStream()
	cmd.Stdin = p.shell.InStream()
	return
}

func parseRedirects(originalArgs []string, originalShell Shell) (parsed *DefaultParsedRedirect, err error) {
	var (
		numArgs int
		inFile  io.ReadCloser
		outFile io.WriteCloser
	)

	parsed = &DefaultParsedRedirect{originalArgs, originalShell, false, false}

	if numArgs = len(parsed.args); numArgs < 2 {
		return
	}

	// check the before-last position of the command
	// for some redirect key and properly handle them.
	switch parsed.args[numArgs-2] {
	case InputRedirect:
		{
			if inFile, err = os.OpenFile(parsed.args[numArgs-1], os.O_RDONLY, os.ModePerm); err != nil {
				return
			}
			parsed.shell.SetInStream(inFile)
			parsed.closeStdin = true
		}
	case OutputRedirect, OutputRedirectAppend:
		{
			var mode int = os.O_CREATE | os.O_WRONLY
			if parsed.args[numArgs-2] == OutputRedirectAppend {
				mode |= os.O_APPEND
			} else {
				mode |= os.O_TRUNC
			}

			if outFile, err = os.OpenFile(parsed.args[numArgs-1], mode, os.ModePerm); err != nil {
				return
			}
			parsed.shell.SetOutStream(outFile)
			parsed.closeStdout = true
		}
	}

	if parsed.closeStdin || parsed.closeStdout {
		// fix arguments removing the redirect
		parsed.args = parsed.args[:numArgs-2]
	}

	return
}
