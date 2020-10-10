package shell

// FakePromptSelect holds data for fake prompt select behavior
type FakePromptSelect struct {
	CalledAsk  bool
	MockAnswer string
	MockError  error
}

// Ask fake behavior for prompting a select question
func (f *FakePromptSelect) Ask(question string, options []string) (answer string, err error) {
	f.CalledAsk = true
	answer = f.MockAnswer
	err = f.MockError
	return
}
