package shell

// FakePromptSelect holds data for fake prompt select behavior
type FakePromptSelect struct {
	CalledAsk  bool
	MockAnswer map[string]string
	MockError  map[string]error
}

// Ask fake behavior for prompting a select question
func (f *FakePromptSelect) Ask(question string, options []string) (answer string, err error) {
	f.CalledAsk = true
	answer = f.MockAnswer[question]
	err = f.MockError[question]
	return
}
