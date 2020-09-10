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

func prepareInputRedirect(path string) (input io.ReadCloser, err error) {
	input, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	return
}
