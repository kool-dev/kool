package shell

import "github.com/moby/term"

// TerminalChecker holds logic to check if environment is a terminal
type TerminalChecker interface {
	IsTerminal(interface{}) bool
}

// DefaultTerminalChecker holds data to check if environment is a terminal
type DefaultTerminalChecker struct{}

// IsTerminal checks if input is a terminal
func (t *DefaultTerminalChecker) IsTerminal(in interface{}) (isTerminal bool) {
	_, isTerminal = term.GetFdInfo(in)
	return
}

// NewTerminalChecker creates a new output writer
func NewTerminalChecker() TerminalChecker {
	return &DefaultTerminalChecker{}
}
