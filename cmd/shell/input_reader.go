package shell

import (
	"io"
	"os"
)

// DefaultInputReader holds reader to read input
type DefaultInputReader struct {
	r io.Reader
}

// InputReader holds logic to read input
type InputReader interface {
	GetReader() io.Reader
	SetReader(io.Reader)
}

// NewInputReader creates a new input reader
func NewInputReader() InputReader {
	return &DefaultInputReader{os.Stdin}
}

// GetReader get default reader
func (r *DefaultInputReader) GetReader() io.Reader {
	return r.r
}

// SetReader set default reader
func (r *DefaultInputReader) SetReader(rd io.Reader) {
	r.r = rd
}
