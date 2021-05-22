package commands

import (
	"io"
	"kool-dev/kool/core/builder"
)

// FakeKoolService is a mock to be used on testing/replacement for KoolService interface
type FakeKoolService struct {
	ArgsExecute        []string
	ExitCode           int
	CalledExecute      bool
	CalledExit         bool
	CalledPrintln      bool
	CalledPrintf       bool
	CalledError        bool
	CalledWarning      bool
	CalledSuccess      bool
	CalledIsTerminal   bool
	CalledInStream     bool
	CalledSetInStream  bool
	CalledOutStream    bool
	CalledSetOutStream bool
	CalledErrStream    bool
	CalledSetErrStream bool
	CalledExec         bool
	CalledInteractive  bool
	CalledLookPath     bool
	MockExecError      error
}

// Execute mocks the function for testing
func (f *FakeKoolService) Execute(args []string) (err error) {
	f.ArgsExecute = args
	f.CalledExecute = true
	err = f.MockExecError
	return
}

// Exit mocks the function for testing
func (f *FakeKoolService) Exit(code int) {
	f.CalledExit = true
	f.ExitCode = code
}

// Println mocks the function for testing
func (f *FakeKoolService) Println(out ...interface{}) {
	f.CalledPrintln = true
}

// Printf mocks the function for testing
func (f *FakeKoolService) Printf(format string, a ...interface{}) {
	f.CalledPrintf = true
}

// Error mocks the function for testing
func (f *FakeKoolService) Error(err error) {
	f.CalledError = true
}

// Warning mocks the function for testing
func (f *FakeKoolService) Warning(out ...interface{}) {
	f.CalledWarning = true
}

// Success mocks the function for testing
func (f *FakeKoolService) Success(out ...interface{}) {
	f.CalledSuccess = true
}

// IsTerminal mocks the function for testing
func (f *FakeKoolService) IsTerminal() (isTerminal bool) {
	f.CalledIsTerminal = true
	return
}

// InStream mocks the function for testing
func (f *FakeKoolService) InStream() (inStream io.Reader) {
	f.CalledInStream = true
	return
}

// SetInStream mocks the function for testing
func (f *FakeKoolService) SetInStream(inStream io.Reader) {
	f.CalledSetInStream = true
}

// OutStream mocks the function for testing
func (f *FakeKoolService) OutStream() (outStream io.Writer) {
	f.CalledOutStream = true
	return
}

// SetOutStream mocks the function for testing
func (f *FakeKoolService) SetOutStream(outStream io.Writer) {
	f.CalledSetOutStream = true
}

// ErrStream mocks the function for testing
func (f *FakeKoolService) ErrStream() (errStream io.Writer) {
	f.CalledErrStream = true
	return
}

// SetErrStream mocks the function for testing
func (f *FakeKoolService) SetErrStream(errStream io.Writer) {
	f.CalledSetErrStream = true
}

// Exec mocks the function for testing
func (f *FakeKoolService) Exec(command builder.Command, extraArgs ...string) (outStr string, err error) {
	f.CalledExec = true
	return
}

// Interactive mocks the function for testing
func (f *FakeKoolService) Interactive(command builder.Command, extraArgs ...string) (err error) {
	f.CalledInteractive = true
	return
}

// LookPath mocks the function for testing
func (f *FakeKoolService) LookPath(command builder.Command) (err error) {
	f.CalledLookPath = true
	return
}
