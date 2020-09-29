package shell

import (
	"fmt"
	"io"
	"os"

	"github.com/gookit/color"
)

// DefaultOutputWriter holds writer to put content
type DefaultOutputWriter struct {
	w io.Writer
}

// OutputWriter holds logic to output content
type OutputWriter interface {
	GetWriter() io.Writer
	SetWriter(io.Writer)
	Println(...interface{})
	Error(error)
	Warning(...interface{})
	Success(...interface{})
}

// NewOutputWriter creates a new output writer
func NewOutputWriter() OutputWriter {
	return &DefaultOutputWriter{os.Stdout}
}

// GetWriter get default writer
func (w *DefaultOutputWriter) GetWriter() io.Writer {
	return w.w
}

// SetWriter set default writer
func (w *DefaultOutputWriter) SetWriter(wr io.Writer) {
	w.w = wr
}

// Println execs Println on writer
func (w *DefaultOutputWriter) Println(out ...interface{}) {
	fmt.Fprintln(w.w, out...)
}

// Error error output
func (w *DefaultOutputWriter) Error(err error) {
	fmt.Fprintf(w.w, "%v\n", color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err))
}

// Warning warning message
func (w *DefaultOutputWriter) Warning(out ...interface{}) {
	warningMessage := color.New(color.Yellow).Sprint(out...)
	fmt.Fprintln(w.w, warningMessage)
}

// Success success message
func (w *DefaultOutputWriter) Success(out ...interface{}) {
	successMessage := color.New(color.Green).Sprint(out...)
	fmt.Fprintln(w.w, successMessage)
}
