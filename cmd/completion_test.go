package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"kool-dev/kool/cmd/shell"
	"strings"
	"testing"
)

func newFakeKoolCompletion() *KoolCompletion {
	return &KoolCompletion{
		DefaultKoolService{
			&shell.FakeExiter{},
			&shell.DefaultOutputWriter{},
			&shell.FakeInputReader{},
			&shell.FakeTerminalChecker{MockIsTerminal: true},
			shell.NewShell(),
		},
		rootCmd,
	}
}

func readOutput(r io.Reader) (output string, err error) {
	var out []byte

	if out, err = ioutil.ReadAll(r); err != nil {
		return
	}

	output = strings.TrimSpace(string(out))
	return
}

func execCompletionCommand(shellType string) (output string, err error) {
	f := newFakeKoolCompletion()
	cmd := NewCompletionCommand(f)

	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{shellType})

	if err = cmd.Execute(); err != nil {
		return
	}

	output, err = readOutput(b)
	return
}

func expectedCompletionOutput(shellType string) (expected string, err error) {
	b := bytes.NewBufferString("")

	switch shellType {
	case "bash":
		err = rootCmd.GenBashCompletion(b)
	case "zsh":
		err = rootCmd.GenZshCompletion(b)
	case "fish":
		err = rootCmd.GenFishCompletion(b, true)
	case "powershell":
		err = rootCmd.GenPowerShellCompletion(b)
	}

	if err != nil {
		return
	}

	expected, err = readOutput(b)
	return
}

func TestNewKoolCompletion(t *testing.T) {
	k := NewKoolCompletion()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Error("unexpected shell.OutputWriter on default KoolCompletion instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Error("unexpected shell.Exiter on default KoolCompletion instance")
	}

	if _, ok := k.DefaultKoolService.in.(*shell.DefaultInputReader); !ok {
		t.Error("unexpected shell.InputReader on default KoolCompletion instance")
	}

	if k.rootCmd.Name() != rootCmd.Name() {
		t.Error("unexpected cobra root Command on default KoolCompletion instance")
	}
}

func TestBashNewCompletionCommand(t *testing.T) {
	var (
		output   string
		expected string
		err      error
	)

	if output, err = execCompletionCommand("bash"); err != nil {
		t.Fatal(err)
	}

	if expected, err = expectedCompletionOutput("bash"); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Error("unexpected bash output for completion command")
	}
}

func TestZshNewCompletionCommand(t *testing.T) {
	var (
		output   string
		expected string
		err      error
	)

	if output, err = execCompletionCommand("zsh"); err != nil {
		t.Fatal(err)
	}

	if expected, err = expectedCompletionOutput("zsh"); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Error("unexpected zsh output for completion command")
	}
}

func TestFishNewCompletionCommand(t *testing.T) {
	var (
		output   string
		expected string
		err      error
	)

	if output, err = execCompletionCommand("fish"); err != nil {
		t.Fatal(err)
	}

	if expected, err = expectedCompletionOutput("fish"); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Error("unexpected fish output for completion command")
	}
}

func TestPowershellNewCompletionCommand(t *testing.T) {
	var (
		output   string
		expected string
		err      error
	)

	if output, err = execCompletionCommand("powershell"); err != nil {
		t.Fatal(err)
	}

	if expected, err = expectedCompletionOutput("powershell"); err != nil {
		t.Fatal(err)
	}

	if output != expected {
		t.Error("unexpected powershell output for completion command")
	}
}
