package shell

import "io"

// FakeOutputWriter is meant to be used for tests - a simple empty mock
// implementing the OutputWriter interface here defined.
type FakeOutputWriter struct {
	Err error
	Out []interface{}

	CalledSetWriter, CalledError, CalledWarning bool
}

// SetWriter is a mocked testing function
func (f *FakeOutputWriter) SetWriter(w io.Writer) {
	f.CalledSetWriter = true
}

// Error is a mocked testing function
func (f *FakeOutputWriter) Error(err error) {
	f.Err = err
	f.CalledError = true
}

// Warning is a mocked testing function
func (f *FakeOutputWriter) Warning(out ...interface{}) {
	f.CalledWarning = true
	f.Out = out
}
