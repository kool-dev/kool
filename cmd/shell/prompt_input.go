package shell

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// ErrPromptInputInterrupted error throwed on signal interrupt
var ErrPromptInputInterrupted error = terminal.InterruptErr

// PromptInput contract that holds logic for prompt a input question
type PromptInput interface {
	Ask(string) (string, error)
}

// DefaultPromptInput holds data for prompting a input question
type DefaultPromptInput struct{}

// NewPromptInput creates a new prompt input
func NewPromptInput() PromptInput {
	return &DefaultPromptInput{}
}

// Ask prompt to the user a input question
func (p *DefaultPromptInput) Ask(question string) (answer string, err error) {
	prompt := &survey.Input{
		Message: question,
	}
	err = survey.AskOne(prompt, &answer)
	return
}
