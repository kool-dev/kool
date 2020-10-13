package shell

import (
	"io"
)

// FakeInputReader is meant to be used for tests - a simple empty mock
// implementing the InputReader interface here defined.
type FakeInputReader struct {
	CalledGetReader, CalledSetReader bool
}

// GetReader is a mocked testing function
func (f *FakeInputReader) GetReader() (r io.Reader) {
	f.CalledGetReader = true
	return
}

// SetReader is a mocked testing function
func (f *FakeInputReader) SetReader(r io.Reader) {
	f.CalledSetReader = true
}
