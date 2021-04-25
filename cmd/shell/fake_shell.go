package shell

import (
	"fmt"
	"io"
	"kool-dev/kool/cmd/builder"
	"strings"
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

	Err           error
	OutLines      []string
	WarningOutput []interface{}
	SuccessOutput []interface{}
	FOutput       string

	CalledPrintln, CalledPrintf, CalledError, CalledWarning, CalledSuccess bool

	MockOutStream io.Writer
	MockErrStream io.Writer
	MockInStream  io.Reader
	MockLookPath  error
}

// InStream is a mocked testing function
func (f *FakeShell) InStream() (inStream io.Reader) {
	f.CalledInStream = true
	return f.MockInStream
}

// SetInStream is a mocked testing function
func (f *FakeShell) SetInStream(inStream io.Reader) {
	f.CalledSetInStream = true
}

// OutStream is a mocked testing function
func (f *FakeShell) OutStream() (outStream io.Writer) {
	f.CalledOutStream = true
	return f.MockOutStream
}

// SetOutStream is a mocked testing function
func (f *FakeShell) SetOutStream(outStream io.Writer) {
	f.CalledSetOutStream = true
}

// ErrStream is a mocked testing function
func (f *FakeShell) ErrStream() (errStream io.Writer) {
	f.CalledErrStream = true
	return f.MockErrStream
}

// SetErrStream is a mocked testing function
func (f *FakeShell) SetErrStream(errStream io.Writer) {
	f.CalledSetErrStream = true
}

// Exec is a mocked testing function
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

// Interactive is a mocked testing function
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
		err = command.(*builder.FakeCommand).MockInteractiveError
	}

	return
}

// LookPath is a mocked testing function
func (f *FakeShell) LookPath(command builder.Command) (err error) {
	if f.CalledLookPath == nil {
		f.CalledLookPath = make(map[string]bool)
	}

	f.CalledLookPath[command.Cmd()] = true

	if _, ok := command.(*builder.FakeCommand); ok {
		err = command.(*builder.FakeCommand).MockLookPathError
		return
	}

	err = f.MockLookPath

	return
}

// Println is a mocked testing function
func (f *FakeShell) Println(out ...interface{}) {
	f.CalledPrintln = true
	f.OutLines = append(f.OutLines, strings.TrimSpace(fmt.Sprintln(out...)))
}

// Printf is a mocked testing function
func (f *FakeShell) Printf(format string, a ...interface{}) {
	f.CalledPrintf = true
	f.FOutput = fmt.Sprintf(format, a...)
}

// Error is a mocked testing function
func (f *FakeShell) Error(err error) {
	f.Err = err
	f.CalledError = true
}

// Warning is a mocked testing function
func (f *FakeShell) Warning(out ...interface{}) {
	f.CalledWarning = true
	f.WarningOutput = out
}

// Success is a mocked testing function
func (f *FakeShell) Success(out ...interface{}) {
	f.CalledSuccess = true
	f.SuccessOutput = out
}
