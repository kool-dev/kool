package shell

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestShellExec(t *testing.T) {
	output, err := Exec("echo", "x")

	if err != nil {
		t.Fatal(err)
	}

	output = strings.TrimSpace(output)

	if output != "x" {
		t.Errorf("Exec failed; expected output 'x', got '%s'", output)
	}
}

func TestShellInteractive(t *testing.T) {
	r, w, err := os.Pipe()

	if err != nil {
		t.Fatal(err)
	}

	originalOutput := os.Stdout
	os.Stdout = w

	defer func(originalOutput *os.File) {
		os.Stdout = originalOutput
	}(originalOutput)

	err = Interactive("echo", "x")

	w.Close()

	if err != nil {
		t.Errorf("Interactive failed; expected no errors 'x', got '%v'", err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)

	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(buf.String())

	if output != "x" {
		t.Errorf("Interactive failed; expected output 'x', got '%s'", output)
	}
}
