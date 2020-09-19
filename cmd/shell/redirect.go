package shell

import (
	"io"
	"os"
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

func parseRedirects(originalArgs []string, defOut io.WriteCloser) (
	args []string,
	in io.ReadCloser,
	out io.WriteCloser,
	closeStdin bool,
	closeStdout bool,
	err error) {
	var numArgs int

	args = originalArgs

	if numArgs = len(args); numArgs < 2 {
		in = os.Stdin
		out = defOut
		return
	}

	in = os.Stdin
	out = defOut

	// check the before-last position of the command
	// for some redirect key and properly handle them.
	switch args[numArgs-2] {
	case InputRedirect:
		{
			if in, err = os.OpenFile(args[numArgs-1], os.O_RDONLY, os.ModePerm); err != nil {
				return
			}
			closeStdin = true
		}
	case OutputRedirect, OutputRedirectAppend:
		{
			var mode int = os.O_CREATE | os.O_WRONLY
			if args[numArgs-2] == OutputRedirectAppend {
				mode |= os.O_APPEND
			} else {
				mode |= os.O_TRUNC
			}

			if out, err = os.OpenFile(args[numArgs-1], mode, os.ModePerm); err != nil {
				return
			}
			closeStdout = true
		}
	}

	if closeStdin || closeStdout {
		// fix arguments removing the redirect
		args = args[:numArgs-2]
	}

	return
}
