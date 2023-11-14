package shell

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewPromptSelect(t *testing.T) {
	p := NewPromptSelect()

	if _, ok := p.(*DefaultPromptSelect); !ok {
		t.Errorf("unexpected PromptSelect on NewPromptSelect")
	}
}

func TestAskPromptSelect(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p := NewPromptSelect()

	_, _ = p.Ask("testing_question", []string{"testing_option1", "testing_option2"})

	w.Close()
	out, err := io.ReadAll(r)
	os.Stdout = oldStdout

	if err != nil {
		t.Fatal(err)
	}

	output := string(out)

	if !strings.Contains(output, "testing_question") || !strings.Contains(output, "testing_option1") || !strings.Contains(output, "testing_option2") {
		t.Error("failed to render the question and its options")
	}
}

func TestConfirmPromptSelect(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p := NewPromptSelect()

	_, _ = p.Confirm("testing_question %s", "arg1")

	w.Close()
	out, err := io.ReadAll(r)
	os.Stdout = oldStdout

	if err != nil {
		t.Fatal(err)
	}

	output := string(out)

	if !strings.Contains(output, "testing_question arg1") || !strings.Contains(output, "Yes") || !strings.Contains(output, "No") {
		t.Error("failed to render the Confirm prompt and its options")
	}
}
