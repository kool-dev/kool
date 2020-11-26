package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/presets"
	"kool-dev/kool/cmd/shell"
	"testing"
)

func newFakeKoolPreset() *KoolPreset {
	return &KoolPreset{
		*newFakeKoolService(),
		&KoolPresetFlags{false},
		&presets.FakeParser{},
		&shell.FakeTerminalChecker{},
		&shell.FakePromptSelect{},
	}
}

func TestNewKoolPreset(t *testing.T) {
	k := NewKoolPreset()

	if _, ok := k.DefaultKoolService.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolPreset instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolPreset instance")
	}

	if _, ok := k.DefaultKoolService.in.(*shell.DefaultInputReader); !ok {
		t.Errorf("unexpected shell.InputReader on default KoolPreset instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolPreset instance")
	} else if k.Flags.Override {
		t.Errorf("bad default value for Override flag on default KoolPreset instance")
	}

	if _, ok := k.parser.(*presets.DefaultParser); !ok {
		t.Errorf("unexpected presets.Parser on default KoolPreset instance")
	}

	if _, ok := k.promptSelect.(*shell.DefaultPromptSelect); !ok {
		t.Errorf("unexpected shell.PromptSelect on default KoolPreset instance")
	}

	if _, ok := k.terminal.(*shell.DefaultTerminalChecker); !ok {
		t.Errorf("unexpected shell.TerminalChecker on default KoolPreset instance")
	}
}

func TestPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.parser.(*presets.FakeParser).MockExists = true
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true
	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.shell.(*shell.FakeShell).CalledPrintln {
		t.Error("did not call Println")
	}

	expected := "Preset laravel is initializing!"
	output := f.shell.(*shell.FakeShell).OutLines[0]

	if expected != output {
		t.Errorf("Expecting message '%s', got '%s'", expected, output)
	}

	if !f.parser.(*presets.FakeParser).CalledLookUpFiles {
		t.Error("did not call parser.LookUpFiles")
	}

	if !f.parser.(*presets.FakeParser).CalledWriteFiles {
		t.Error("did not call parser.WriteFiles")
	}

	if !f.shell.(*shell.FakeShell).CalledSuccess {
		t.Error("did not call Success")
	}

	expected = "Preset laravel initialized!"
	output = fmt.Sprint(f.shell.(*shell.FakeShell).SuccessOutput...)

	if expected != output {
		t.Errorf("Expecting success message '%s', got '%s'", expected, output)
	}
}

func TestInvalidScriptPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true
	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"invalid"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	expected := "Unknown preset invalid"
	output := f.shell.(*shell.FakeShell).Err.Error()

	if expected != output {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestExistingFilesPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockFoundFiles = []string{"kool.yml"}
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true
	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledWarning {
		t.Error("did not call Warning")
	}

	expected := "Some preset files already exist. In case you wanna override them, use --override."
	output := fmt.Sprint(f.shell.(*shell.FakeShell).WarningOutput...)

	if output != expected {
		t.Errorf("expecting message '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestOverrideFilesPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockFoundFiles = []string{"kool.yml"}
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"--override", "laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if f.parser.(*presets.FakeParser).CalledLookUpFiles {
		t.Error("unexpected existing files checking")
	}

	if f.shell.(*shell.FakeShell).CalledWarning {
		t.Error("unexpected existing files Warning")
	}

	if f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("unexpected program Exit")
	}

	if !f.shell.(*shell.FakeShell).CalledSuccess {
		t.Error("did not call Success")
	}
}

func TestWriteErrorPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockError = errors.New("write error")
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true

	cmd := NewPresetCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	expected := "Failed to write preset file : write error"
	output := f.shell.(*shell.FakeShell).Err.Error()

	if output != expected {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestNoArgsPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	mockAnswer := make(map[string]string)
	mockAnswer["What language do you want to use"] = "php"
	mockAnswer["What preset do you want to use"] = "laravel"

	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = mockAnswer
	f.parser.(*presets.FakeParser).MockLanguages = []string{"php"}
	f.parser.(*presets.FakeParser).MockPresets = []string{"laravel"}
	f.parser.(*presets.FakeParser).MockExists = true
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	expected := "Preset laravel is initializing!"
	output := f.shell.(*shell.FakeShell).OutLines[0]

	if expected != output {
		t.Errorf("Expecting message '%s', got '%s'", expected, output)
	}
}

func TestFailingLanguageNoArgsPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.parser.(*presets.FakeParser).MockLanguages = []string{"php"}
	f.parser.(*presets.FakeParser).MockPresets = []string{"laravel"}

	mockError := make(map[string]error)
	mockError["What language do you want to use"] = errors.New("error prompt select language")

	f.promptSelect.(*shell.FakePromptSelect).MockError = mockError
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	expected := "error prompt select language"
	output := f.shell.(*shell.FakeShell).Err.Error()

	if output != expected {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestFailingPresetNoArgsPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.parser.(*presets.FakeParser).MockLanguages = []string{"php"}
	f.parser.(*presets.FakeParser).MockPresets = []string{"laravel"}

	mockAnswer := make(map[string]string)
	mockAnswer["What language do you want to use"] = "php"

	f.promptSelect.(*shell.FakePromptSelect).MockAnswer = mockAnswer

	mockError := make(map[string]error)
	mockError["What preset do you want to use"] = errors.New("error prompt select preset")

	f.promptSelect.(*shell.FakePromptSelect).MockError = mockError
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.promptSelect.(*shell.FakePromptSelect).CalledAsk {
		t.Error("did not call Ask on PromptSelect")
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	expected := "error prompt select preset"
	output := f.shell.(*shell.FakeShell).Err.Error()

	if output != expected {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestCancellingPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()

	mockError := make(map[string]error)
	mockError["What language do you want to use"] = shell.ErrPromptSelectInterrupted

	f.promptSelect.(*shell.FakePromptSelect).MockError = mockError
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = true

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledWarning {
		t.Error("did not call Warning")
	}

	expected := "Operation Cancelled\n"
	output := fmt.Sprintln(f.shell.(*shell.FakeShell).WarningOutput...)

	if output != expected {
		t.Errorf("expecting warning '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}

	if f.exiter.(*shell.FakeExiter).Code() != 0 {
		t.Error("did not call Exit with code 0")
	}
}

func TestNonTTYPresetCommand(t *testing.T) {
	f := newFakeKoolPreset()
	f.term.(*shell.FakeTerminalChecker).MockIsTerminal = false

	cmd := NewPresetCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.shell.(*shell.FakeShell).CalledError {
		t.Error("did not call Error")
	}

	err := f.shell.(*shell.FakeShell).Err

	if err == nil {
		t.Error("expecting an error, got none")
	} else if err.Error() != "the input device is not a TTY; for non-tty environments, please specify a preset argument" {
		t.Errorf("expecting error 'the input device is not a TTY; for non-tty environments, please specify a preset argument', got %v", err)
	}
}
