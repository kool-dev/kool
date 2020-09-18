package shell

import (
	"fmt"
	"io"

	"github.com/gookit/color"
)

// ExecError error output
func ExecError(out string, err error) {
	if err != nil {
		color.New(color.BgRed, color.FgWhite).Println("error:", err)
	}

	if out != "" {
		fmt.Println("Output:", out)
	}
}

// FexecError error output
func FexecError(w io.Writer, out string, err error) {
	if err != nil {
		errorMessage := color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err)
		fmt.Fprintf(w, "%v\n", errorMessage)
	}

	if out != "" {
		fmt.Fprintf(w, "Output: %s\n", out)
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

// Fwarning warning message
func Fwarning(w io.Writer, out ...interface{}) {
	warning := color.New(color.Yellow).Sprint(out...)
	fmt.Fprintf(w, "%s\n", warning)
}

// Success Success message
func Success(out ...interface{}) {
	color.New(color.Green).Println(out...)
}
