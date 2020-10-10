package shell

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// ErrPromptSelectInterrupted error throwed on signal interrupt
var ErrPromptSelectInterrupted error = terminal.InterruptErr

// PromptSelect contract that holds logic for prompt a select question
type PromptSelect interface {
	Ask(string, []string) (string, error)
}

// DefaultPromptSelect holds data for prompting a select question
type DefaultPromptSelect struct{}

// NewPromptSelect creates a new prompt select
func NewPromptSelect() PromptSelect {
	return &DefaultPromptSelect{}
}

// Ask prompt to the user a select question
func (p *DefaultPromptSelect) Ask(question string, options []string) (answer string, err error) {
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	err = survey.AskOne(prompt, &answer)
	return
}
