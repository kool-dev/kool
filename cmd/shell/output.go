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
	SetWriter(io.Writer)
	Error(error)
	Warning(...interface{})
}

// NewOutputWriter creates a new output writer
func NewOutputWriter() OutputWriter {
	return &DefaultOutputWriter{os.Stdout}
}

// SetWriter set default writer
func (w *DefaultOutputWriter) SetWriter(wr io.Writer) {
	w.w = wr
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

// ExecError error output
func ExecError(out string, err error) {
	if err != nil {
		color.New(color.BgRed, color.FgWhite).Println("error:", err)
	}

	if out != "" {
		fmt.Println("Output:", out)
	}
}

// Error error message
func Error(out ...interface{}) {
	color.New(color.BgRed, color.FgLightWhite).Println(out...)
}

// Warning warning message
func Warning(out ...interface{}) {
	color.New(color.Yellow).Println(out...)
}

// Success Success message
func Success(out ...interface{}) {
	color.New(color.Green).Println(out...)
}
