package cmd

import "io"

// FakeKoolService is a mock to be used on testing/replacement for KoolService interface
type FakeKoolService struct {
	ArgsExecute     []string
	ExitCode        int
	CalledExecute   bool
	CalledExit      bool
	CalledSetWriter bool
	CalledError     bool
	CalledWarning   bool
}

// Execute mocks the function for testing
func (f *FakeKoolService) Execute(args []string) (err error) {
	f.ArgsExecute = args
	f.CalledExecute = true
	return
}

// Exit mocks the function for testing
func (f *FakeKoolService) Exit(code int) {
	f.CalledExit = true
	f.ExitCode = code
}

// SetWriter mocks the function for testing
func (f *FakeKoolService) SetWriter(w io.Writer) {
	f.CalledSetWriter = true
}

// Error mocks the function for testing
func (f *FakeKoolService) Error(err error) {
	f.CalledError = true
}

// Warning mocks the function for testing
func (f *FakeKoolService) Warning(out ...interface{}) {
	f.CalledWarning = true
}
