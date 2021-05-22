package shell

// FakeExiter holds a fake mock implementing Exiter interface
type FakeExiter struct {
	exited bool
	code   int
}

// Exit implements
func (f *FakeExiter) Exit(code int) {
	f.exited = true
	f.code = code
}

// Code returns the code given the last call to Exit
func (f *FakeExiter) Code() int {
	return f.code
}

// Exited tells whether the Exit method was called
func (f *FakeExiter) Exited() bool {
	return f.exited
}
