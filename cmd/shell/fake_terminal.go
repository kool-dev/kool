package shell

// FakeTerminalChecker holds a fake mock implementing TerminalChecker interface
type FakeTerminalChecker struct {
	CalledIsTerminal bool
	MockIsTerminal   bool
}

// IsTerminal implements fake IsTerminal
func (f *FakeTerminalChecker) IsTerminal(fds ...interface{}) (isTerminal bool) {
	f.CalledIsTerminal = true
	isTerminal = f.MockIsTerminal
	return
}
