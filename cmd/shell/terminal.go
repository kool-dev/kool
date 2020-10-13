package shell

import "github.com/moby/term"

// TerminalChecker holds logic to check if environment is a terminal
type TerminalChecker interface {
	IsTerminal(interface{}, interface{}) bool
}

// DefaultTerminalChecker holds data to check if environment is a terminal
type DefaultTerminalChecker struct{}

// IsTerminal checks if input is a terminal
func (t *DefaultTerminalChecker) IsTerminal(in interface{}, out interface{}) (isTerminal bool) {
	_, inIsTerminal := term.GetFdInfo(in)
	_, outIsTerminal := term.GetFdInfo(out)

	isTerminal = inIsTerminal && outIsTerminal
	return
}

// NewTerminalChecker creates a new terminal checker
func NewTerminalChecker() TerminalChecker {
	return &DefaultTerminalChecker{}
}
