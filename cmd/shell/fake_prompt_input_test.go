package shell

import (
	"errors"
	"testing"
)

func TestFakePromptInput(t *testing.T) {
	f := &FakePromptInput{}
	f.MockAnswer = make(map[string]string)
	f.MockAnswer["question"] = "answer"

	answer, err := f.Ask("question")

	if err != nil {
		t.Errorf("unexpected error on Ask: %v", err)
	}

	if answer != "answer" {
		t.Errorf("expecting answer 'answer', got %s", answer)
	}

	f.MockError = make(map[string]error)
	f.MockError["question"] = errors.New("ask error")

	_, err = f.Ask("question")

	if err == nil {
		t.Error("should throw an error on Ask")
	} else if err.Error() != "ask error" {
		t.Errorf("expecting error 'ask error', got %v", err)
	}
}
