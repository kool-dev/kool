package shell

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/creack/pty"
	"github.com/gookit/color"
)

func TestJSONErrorIsSingleDocument(t *testing.T) {
	t.Setenv("KOOL_OUTPUT", "json")
	s := NewShell()
	var output bytes.Buffer
	s.SetErrStream(&output)
	s.Error(errors.New("broken"))
	s.Error(ReportedError{Err: errors.New("already emitted")})
	var payload map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &payload); err != nil || payload["error"] != "broken" {
		t.Fatalf("invalid diagnostic %q: %v", output.String(), err)
	}
}

func TestTaskPipeKeepsTerminalColorAndSpacing(t *testing.T) {
	master, terminal, err := pty.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer master.Close()
	defer terminal.Close()
	previous := color.Enable
	color.Enable = true
	t.Cleanup(func() { color.Enable = previous })
	t.Setenv("KOOL_OUTPUT", "text")
	s := NewShell().(*DefaultShell)
	s.SetOutStream(terminal)
	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()
	s.SetOutStream(writer)
	if !s.useColor() {
		t.Fatal("task pipe lost terminal color")
	}
	var plain bytes.Buffer
	s.SetOutStream(&plain)
	s.Info("\t", "script")
	if plain.String() != "\tscript\n" {
		t.Fatalf("unexpected spacing %q", plain.String())
	}
	t.Setenv("KOOL_OUTPUT", "json")
	s.SetOutStream(terminal)
	s.SetErrStream(&plain)
	if s.useColor() {
		t.Fatal("redirected diagnostics should not be colored")
	}
	s.Warning("warning")
	if strings.Contains(plain.String(), "\x1b") {
		t.Fatal("ANSI in redirected diagnostics")
	}
}
