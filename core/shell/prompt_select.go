package shell

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// PromptSelect contract that holds logic for prompt a select question
type PromptSelect interface {
	Ask(string, []string) (string, error)

	Confirm(string, ...any) (bool, error)
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
	if err = survey.AskOne(prompt, &answer); err != nil && err == terminal.InterruptErr {
		err = ErrUserCancelled
	}
	return
}

// Confirm prompts to the user a Yes/No confirm question
func (p *DefaultPromptSelect) Confirm(question string, args ...any) (confirmed bool, err error) {
	if args != nil {
		question = fmt.Sprintf(question, args...)
	}

	var answer string

	if answer, err = p.Ask(question, []string{"Yes", "No"}); err != nil {
		return
	}

	confirmed = answer == "Yes"
	return
}
