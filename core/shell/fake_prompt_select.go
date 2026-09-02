package shell

// FakePromptSelect holds data for fake prompt select behavior
type FakePromptSelect struct {
	CalledAsk  bool
	MockAnswer map[string]string
	MockError  map[string]error

	CalledConfirm []*struct {
		question string
		args     []any
	}
	MockConfirm      map[string]bool
	MockConfirmError map[string]error
}

// Ask mocked behavior for testing prompting a select question
func (f *FakePromptSelect) Ask(question string, options []string) (answer string, err error) {
	f.CalledAsk = true
	answer = f.MockAnswer[question]
	err = f.MockError[question]
	return
}

// AskMany mocked behavior for testing multiple option prompts.
func (f *FakePromptSelect) AskMany(question string, options []string) (answers []string, err error) {
	f.CalledAsk = true
	if answer := f.MockAnswer[question]; answer != "" {
		answers = []string{answer}
	}
	err = f.MockError[question]
	return
}

// Confirm mocked behavior for testing prompting a confirm question
func (f *FakePromptSelect) Confirm(question string, args ...any) (confirmed bool, err error) {
	f.CalledConfirm = append(f.CalledConfirm, &struct {
		question string
		args     []any
	}{question, args})
	confirmed = f.MockConfirm[question]
	err = f.MockConfirmError[question]
	return
}
