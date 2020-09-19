package shell

import (
	"fmt"
	"io"
	"os"

	"github.com/gookit/color"
)

// DefaultOutputWriter holds writer to put content
type DefaultOutputWriter struct {
	Writer io.Writer
}

// OutputWriter holds logic to output content
type OutputWriter interface {
	GetWriter() io.Writer
	SetWriter(writer io.Writer)
	ExecError(out string, err error)
	Error(out ...interface{})
	Warning(out ...interface{})
	Success(out ...interface{})
}

// NewOutputWriter creates a new output writer
func NewOutputWriter() OutputWriter {
	return &DefaultOutputWriter{os.Stdout}
}

// SetWriter set default writer
func (w *DefaultOutputWriter) SetWriter(writer io.Writer) {
	w.Writer = writer
}

// GetWriter get default writer
func (w *DefaultOutputWriter) GetWriter() io.Writer {
	return w.Writer
}

// ExecError error output
func (w *DefaultOutputWriter) ExecError(out string, err error) {
	if err != nil {
		errorMessage := color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err)
		fmt.Fprintf(w.Writer, "%v\n", errorMessage)
	}

	if out != "" {
		fmt.Fprintf(w.Writer, "Output: %s\n", out)
	}
}

// ExecError error output
func ExecError(out string, err error) {
	outputwriter := NewOutputWriter()
	outputwriter.ExecError(out, err)
}

// Error error message
func (w *DefaultOutputWriter) Error(out ...interface{}) {
	errorMessage := color.New(color.BgRed, color.FgWhite).Sprint(out...)
	fmt.Fprintln(w.Writer, errorMessage)
}

// Error error message
func Error(out ...interface{}) {
	outputwriter := NewOutputWriter()
	outputwriter.Error(out...)
}

// Warning warning message
func (w *DefaultOutputWriter) Warning(out ...interface{}) {
	warningMessage := color.New(color.Yellow).Sprint(out...)
	fmt.Fprintln(w.Writer, warningMessage)
}

// Warning warning message
func Warning(out ...interface{}) {
	outputwriter := NewOutputWriter()
	outputwriter.Warning(out...)
}

// Success Success message
func (w *DefaultOutputWriter) Success(out ...interface{}) {
	successMessage := color.New(color.Green).Sprint(out...)
	fmt.Fprintln(w.Writer, successMessage)
}

// Success Success message
func Success(out ...interface{}) {
	outputwriter := NewOutputWriter()
	outputwriter.Success(out...)
}
