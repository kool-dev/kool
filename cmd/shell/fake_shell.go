package shell

import (
	"io"
	"kool-dev/kool/cmd/builder"
)

// FakeShell fake shell data
type FakeShell struct {
	CalledInStream     bool
	CalledSetInStream  bool
	CalledOutStream    bool
	CalledSetOutStream bool
	CalledErrStream    bool
	CalledSetErrStream bool
	CalledExec         map[string]bool
	CalledInteractive  map[string]bool
	CalledLookPath     map[string]bool
	ArgsInteractive    map[string][]string

	MockOutStream io.Writer
	MockErrStream io.Writer
	MockInStream  io.Reader
}

// InStream get input stream
func (f *FakeShell) InStream() (inStream io.Reader) {
	f.CalledInStream = true
	return f.MockInStream
}

// SetInStream set input stream
func (f *FakeShell) SetInStream(inStream io.Reader) {
	f.CalledSetInStream = true
}

// OutStream get output stream
func (f *FakeShell) OutStream() (outStream io.Writer) {
	f.CalledOutStream = true
	return f.MockOutStream
}

// SetOutStream set output stream
func (f *FakeShell) SetOutStream(outStream io.Writer) {
	f.CalledSetOutStream = true
}

// ErrStream get error stream
func (f *FakeShell) ErrStream() (errStream io.Writer) {
	f.CalledErrStream = true
	return f.MockErrStream
}

// SetErrStream set error stream
func (f *FakeShell) SetErrStream(errStream io.Writer) {
	f.CalledSetErrStream = true
}

// Exec will execute the given command silently and return the combined
// error/standard output, and an error if any.
func (f *FakeShell) Exec(command builder.Command, extraArgs ...string) (outStr string, err error) {
	if f.CalledExec == nil {
		f.CalledExec = make(map[string]bool)
	}

	f.CalledExec[command.Cmd()] = true

	if _, ok := command.(*builder.FakeCommand); ok {
		err = command.(*builder.FakeCommand).MockExecError
		outStr = command.(*builder.FakeCommand).MockExecOut
	}
	return
}

// Interactive runs the given command proxying current Stdin/Stdout/Stderr
// which makes it interactive for running even something like `bash`.
func (f *FakeShell) Interactive(command builder.Command, extraArgs ...string) (err error) {
	if f.CalledInteractive == nil {
		f.CalledInteractive = make(map[string]bool)
	}

	if f.ArgsInteractive == nil {
		f.ArgsInteractive = make(map[string][]string)
	}

	f.CalledInteractive[command.Cmd()] = true
	f.ArgsInteractive[command.Cmd()] = extraArgs

	if _, ok := command.(*builder.FakeCommand); ok {
		err = command.(*builder.FakeCommand).MockError
	}

	return
}

// LookPath returns if the command exists
func (f *FakeShell) LookPath(command builder.Command) (err error) {
	if f.CalledLookPath == nil {
		f.CalledLookPath = make(map[string]bool)
	}

	f.CalledLookPath[command.Cmd()] = true

	if _, ok := command.(*builder.FakeCommand); ok {
		err = command.(*builder.FakeCommand).MockLookPathError
	}

	return
}
