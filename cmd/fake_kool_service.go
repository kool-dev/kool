package cmd

import (
	"io"
	"kool-dev/kool/cmd/builder"
)

// FakeKoolService is a mock to be used on testing/replacement for KoolService interface
type FakeKoolService struct {
	ArgsExecute        []string
	ExitCode           int
	CalledExecute      bool
	CalledExit         bool
	CalledGetWriter    bool
	CalledSetWriter    bool
	CalledGetReader    bool
	CalledSetReader    bool
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

// GetWriter mocks the function for testing
func (f *FakeKoolService) GetWriter() (w io.Writer) {
	f.CalledGetWriter = true
	return
}

// SetWriter mocks the function for testing
func (f *FakeKoolService) SetWriter(w io.Writer) {
	f.CalledSetWriter = true
}

// GetReader mocks the function for testing
func (f *FakeKoolService) GetReader() (r io.Reader) {
	f.CalledGetReader = true
	return
}

// SetReader mocks the function for testing
func (f *FakeKoolService) SetReader(r io.Reader) {
	f.CalledSetReader = true
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

// InStream get input stream
func (f *FakeKoolService) InStream() (inStream io.Reader) {
	f.CalledInStream = true
	return
}

// SetInStream set input stream
func (f *FakeKoolService) SetInStream(inStream io.Reader) {
	f.CalledSetInStream = true
}

// OutStream get output stream
func (f *FakeKoolService) OutStream() (outStream io.Writer) {
	f.CalledOutStream = true
	return
}

// SetOutStream set output stream
func (f *FakeKoolService) SetOutStream(outStream io.Writer) {
	f.CalledSetOutStream = true
}

// ErrStream get error stream
func (f *FakeKoolService) ErrStream() (errStream io.Writer) {
	f.CalledErrStream = true
	return
}

// SetErrStream set error stream
func (f *FakeKoolService) SetErrStream(errStream io.Writer) {
	f.CalledSetErrStream = true
}

// Exec will execute the given command silently and return the combined
// error/standard output, and an error if any.
func (f *FakeKoolService) Exec(command builder.Command, extraArgs ...string) (outStr string, err error) {
	f.CalledExec = true
	return
}

// Interactive runs the given command proxying current Stdin/Stdout/Stderr
// which makes it interactive for running even something like `bash`.
func (f *FakeKoolService) Interactive(command builder.Command, extraArgs ...string) (err error) {
	f.CalledInteractive = true
	return
}

// LookPath returns if the command exists
func (f *FakeKoolService) LookPath(command builder.Command) (err error) {
	f.CalledLookPath = true
	return
}
