package shell

import (
	"kool-dev/kool/core/builder"
	"os"
	"path/filepath"
	"testing"
)

type fakeReaderWriterCloser struct {
	calledClose bool
}

func (f *fakeReaderWriterCloser) Close() error {
	f.calledClose = true
	return nil
}

func (f *fakeReaderWriterCloser) Read(p []byte) (n int, err error) {
	return
}

func (f *fakeReaderWriterCloser) Write(p []byte) (n int, err error) {
	return
}

func TestParseRedirectParseNoRedirects(t *testing.T) {
	f := &FakeShell{}
	// test no redirects
	p, err := parseRedirects(builder.NewCommand("foo", "bar"), f)

	if err != nil {
		t.Error("unexpected error parsing redirects")
	}

	if p.hasCustomStdin || p.hasCustomStdout {
		t.Error("bad parse - should not custom in/outs")
	}

	// test input redirect
	input := filepath.Join(t.TempDir(), "input")
	file, _ := os.Create(input)
	file.Close()

	s := NewShell()
	s.SetInStream(file)
	p, err = parseRedirects(builder.NewCommand("foo", "<", input), s)

	if err != nil {
		t.Error("unexpected error parsing redirects")
	}

	if !p.hasCustomStdin || p.hasCustomStdout {
		t.Errorf("bad parse - should have custom in; not out")
	}

	p.Close()

	// test output
	output := filepath.Join(t.TempDir(), "output")

	s = NewShell()
	s.SetOutStream(file)
	p, err = parseRedirects(builder.NewCommand("foo", ">", output), s)

	if err != nil {
		t.Errorf("unexpected error parsing redirects")
	}

	if p.hasCustomStdin || !p.hasCustomStdout {
		t.Errorf("bad parse - should have custom in; not out")
	}

	p.Close()
}

func TestCommandWithPointers(t *testing.T) {
	ptr := &CommandWithPointers{
		Command: builder.NewCommand("foo", "bar"),
		in:      &fakeReaderWriterCloser{},
		out:     &fakeReaderWriterCloser{},
	}

	cmd := ptr.Cmd()

	if cmd == nil {
		t.Errorf("failed to create command")
		return
	}

	if cmd.Args == nil || cmd.Args[0] != "foo" || cmd.Args[1] != "bar" {
		t.Errorf("bad command/arguments for created Commands")
	}

	// calls close - should not clode in/out
	ptr.Close()

	if ptr.in.(*fakeReaderWriterCloser).calledClose {
		t.Errorf("did not expect to call close on Stdin")
	}
	if ptr.out.(*fakeReaderWriterCloser).calledClose {
		t.Errorf("did not expect to call close on Stdout")
	}

	ptr.hasCustomStdin = true
	ptr.hasCustomStdout = true

	// calls close - should have custom in/out
	ptr.Close()

	if !ptr.in.(*fakeReaderWriterCloser).calledClose {
		t.Errorf("did not get expected call to close on Stdin")
	}
	if !ptr.out.(*fakeReaderWriterCloser).calledClose {
		t.Errorf("did not get expected call to close on Stdout")
	}
}

func TestHasRedirect(t *testing.T) {
	if !hasRedirect(builder.NewCommand("cmd", ">", "out")) {
		t.Error("failed telling redirection >")
	}
	if !hasRedirect(builder.NewCommand("cmd", "<", "in")) {
		t.Error("failed telling redirection <")
	}
	if !hasRedirect(builder.NewCommand("cmd", ">>", "in")) {
		t.Error("failed telling redirection >>")
	}

	if hasRedirect(builder.NewCommand("cmd")) || hasRedirect(builder.NewCommand("cmd", "arg")) {
		t.Error("false positive for redirection")
	}
}
