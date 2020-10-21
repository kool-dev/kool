package shell

import (
	"errors"
	"testing"
)

func TestFakePromptSelect(t *testing.T) {
	f := &FakePromptSelect{}
	f.MockAnswer = make(map[string]string)
	f.MockAnswer["question"] = "answer"

	answer, err := f.Ask("question", []string{"option"})

	if err != nil {
		t.Errorf("unexpected error on Ask: %v", err)
	}

	if answer != "answer" {
		t.Errorf("expecting answer 'answer', got %s", answer)
	}

	f.MockError = make(map[string]error)
	f.MockError["question"] = errors.New("error")

	_, err = f.Ask("question", []string{"option"})

	if err == nil {
		t.Errorf("should throw an error on Ask")
	}
}
