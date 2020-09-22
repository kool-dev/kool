package shell

import "os"

type exitFN func(code int)

// DefaultExiter holds exiter data
type DefaultExiter struct {
	osExit exitFN
}

// Exiter holds exiter logic
type Exiter interface {
	Exit(int)
}

// NewExiter initialize a new exiter
func NewExiter() Exiter {
	return &DefaultExiter{os.Exit}
}

// Exit exit with a exit code
func (e *DefaultExiter) Exit(code int) {
	e.osExit(code)
}
