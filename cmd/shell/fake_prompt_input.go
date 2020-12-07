package shell

// FakePromptInput holds data for fake prompt input behavior
type FakePromptInput struct {
	CalledAsk  bool
	MockAnswer map[string]string
	MockError  map[string]error
}

// Ask fake behavior for prompting a input question
func (f *FakePromptInput) Ask(question string) (answer string, err error) {
	f.CalledAsk = true
	answer = f.MockAnswer[question]
	err = f.MockError[question]
	return
}
