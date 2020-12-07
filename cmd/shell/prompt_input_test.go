package shell

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewPromptInput(t *testing.T) {
	p := NewPromptInput()

	if _, ok := p.(*DefaultPromptInput); !ok {
		t.Errorf("unexpected PromptInput on NewPromptInput")
	}
}

func TestAskPromptInput(t *testing.T) {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	p := NewPromptInput()

	_, _ = p.Ask("testing_question")

	w.Close()
	out, err := ioutil.ReadAll(r)
	os.Stdout = oldStdout

	if err != nil {
		t.Fatal(err)
	}

	output := string(out)

	if !strings.Contains(output, "testing_question") {
		t.Error("failed to render the question")
	}
}
