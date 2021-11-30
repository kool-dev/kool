package shell

import "github.com/moby/term"

const standardTermWidth int = 80

// TerminalChecker holds logic to check if environment is a terminal
type TerminalChecker interface {
	IsTerminal(...interface{}) bool
}

// DefaultTerminalChecker holds logic to check if file descriptors are a TTY
type DefaultTerminalChecker struct{}

// IsTerminal checks if input is a terminal
func (t *DefaultTerminalChecker) IsTerminal(fds ...interface{}) bool {
	for i := range fds {
		if _, isTerminal := term.GetFdInfo(fds[i]); !isTerminal {
			return false
		}
	}

	return true
}

// NewTerminalChecker creates a new terminal checker
func NewTerminalChecker() TerminalChecker {
	return &DefaultTerminalChecker{}
}
