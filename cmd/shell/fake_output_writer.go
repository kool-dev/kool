package shell

import "io"

// FakeOutputWriter is meant to be used for tests - a simple empty mock
// implementing the OutputWriter interface here defined.
type FakeOutputWriter struct {
	Err           error
	Out           []interface{}
	WarningOutput []interface{}
	SuccessOutput []interface{}

	CalledGetWriter, CalledSetWriter, CalledPrintln, CalledError, CalledWarning, CalledSuccess bool
}

// GetWriter is a mocked testing function
func (f *FakeOutputWriter) GetWriter() (w io.Writer) {
	f.CalledGetWriter = true
	return
}

// SetWriter is a mocked testing function
func (f *FakeOutputWriter) SetWriter(w io.Writer) {
	f.CalledSetWriter = true
}

// Println is a mocked testing function
func (f *FakeOutputWriter) Println(out ...interface{}) {
	f.CalledPrintln = true
	f.Out = out
}

// Error is a mocked testing function
func (f *FakeOutputWriter) Error(err error) {
	f.Err = err
	f.CalledError = true
}

// Warning is a mocked testing function
func (f *FakeOutputWriter) Warning(out ...interface{}) {
	f.CalledWarning = true
	f.WarningOutput = out
}

// Success is a mocked testing function
func (f *FakeOutputWriter) Success(out ...interface{}) {
	f.CalledSuccess = true
	f.SuccessOutput = out
}
