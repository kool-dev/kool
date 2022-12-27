package shell

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// PromptInput contract that holds logic for prompt an input string
type PromptInput interface {
	Input(string, string) (string, error)
}

// DefaultPromptInput holds data for prompting a select question
type DefaultPromptInput struct{}

// NewPromptInput creates a new prompt select
func NewPromptInput() PromptInput {
	return &DefaultPromptInput{}
}

// Ask prompt to the user a select question
func (p *DefaultPromptInput) Input(question string, defaultInput string) (input string, err error) {
	prompt := &survey.Input{
		Message: question,
		Default: defaultInput,
	}
	if err = survey.AskOne(prompt, &input); err != nil && err == terminal.InterruptErr {
		err = ErrUserCancelled
	}
	return
}
