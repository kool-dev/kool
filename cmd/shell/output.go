package shell

import (
	"fmt"

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
