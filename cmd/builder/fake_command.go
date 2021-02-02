package builder

// FakeCommand implements the Command interface and is used for mocking on testing scenarios
type FakeCommand struct {
	ArgsAppend         []string
	CalledAppendArgs   bool
	CalledString       bool
	CalledCmd          bool
	CalledArgs         bool
	CalledReset        bool
	CalledParseCommand bool

	MockCmd              string
	MockExecOut          string
	MockError            error
	MockLookPathError    error
	MockExecError        error
	MockInteractiveError error
}

// AppendArgs mocked function for testing
func (f *FakeCommand) AppendArgs(args ...string) {
	f.ArgsAppend = append(f.ArgsAppend, args...)
	f.CalledAppendArgs = true
}

// String mocked function for testing
func (f *FakeCommand) String() string {
	f.CalledString = true
	return ""
}

// Args returns the command arguments
func (f *FakeCommand) Args() []string {
	f.CalledArgs = true
	return f.ArgsAppend
}

// Reset resets arguments
func (f *FakeCommand) Reset() {
	f.CalledReset = true
}

// Cmd returns the command executable
func (f *FakeCommand) Cmd() string {
	f.CalledCmd = true
	return f.MockCmd
}

// Parse call the ParseCommand function
func (f *FakeCommand) Parse(line string) (err error) {
	f.CalledParseCommand = true
	err = f.MockError
	return
}
