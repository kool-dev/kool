package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/presets"
	"kool-dev/kool/cmd/shell"
	"testing"
)

func newFakeKoolInit() *KoolInit {
	return &KoolInit{
		*newFakeKoolService(),
		&KoolInitFlags{false},
		&presets.FakeParser{},
	}
}

func TestNewKoolInit(t *testing.T) {
	k := NewKoolInit()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Errorf("unexpected shell.OutputWriter on default KoolInit instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolInit instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolInit instance")
	} else if k.Flags.Override {
		t.Errorf("bad default value for Override flag on default KoolInit instance")
	}

	if _, ok := k.parser.(*presets.DefaultParser); !ok {
		t.Errorf("unexpected presets.Parser on default KoolInit instance")
	}
}

func TestInitCommand(t *testing.T) {
	f := newFakeKoolInit()
	f.parser.(*presets.FakeParser).MockExists = true
	cmd := NewInitCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing init command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Error("did not call SetWriter")
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.out.(*shell.FakeOutputWriter).CalledPrintln {
		t.Error("did not call Println")
	}

	expected := "Preset laravel is initializing!\n"
	output := fmt.Sprintln(f.out.(*shell.FakeOutputWriter).Out...)

	if expected != output {
		t.Errorf("Expecting message '%s', got '%s'", expected, output)
	}

	if !f.parser.(*presets.FakeParser).CalledLookUpFiles {
		t.Error("did not call parser.LookUpFiles")
	}

	if !f.parser.(*presets.FakeParser).CalledWriteFiles {
		t.Error("did not call parser.WriteFiles")
	}

	if !f.out.(*shell.FakeOutputWriter).CalledSuccess {
		t.Error("did not call Success")
	}

	expected = "Preset laravel initialized!"
	output = fmt.Sprint(f.out.(*shell.FakeOutputWriter).SuccessOutput...)

	if expected != output {
		t.Errorf("Expecting success message '%s', got '%s'", expected, output)
	}
}

func TestInvalidScriptInitCommand(t *testing.T) {
	f := newFakeKoolInit()
	cmd := NewInitCommand(f)

	cmd.SetArgs([]string{"invalid"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing init command; error: %v", err)
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.out.(*shell.FakeOutputWriter).CalledError {
		t.Error("did not call Error")
	}

	expected := "Unknown preset invalid"
	output := f.out.(*shell.FakeOutputWriter).Err.Error()

	if expected != output {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestExistingFilesInitCommand(t *testing.T) {
	f := newFakeKoolInit()
	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockFoundFiles = []string{"kool.yml"}
	cmd := NewInitCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing init command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledWarning {
		t.Error("did not call Warning")
	}

	expected := "Some preset files already exist. In case you wanna override them, use --override."
	output := fmt.Sprint(f.out.(*shell.FakeOutputWriter).WarningOutput...)

	if output != expected {
		t.Errorf("expecting message '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestOverrideFilesInitCommand(t *testing.T) {
	f := newFakeKoolInit()
	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockFoundFiles = []string{"kool.yml"}

	cmd := NewInitCommand(f)

	cmd.SetArgs([]string{"--override", "laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing init command; error: %v", err)
	}

	if f.parser.(*presets.FakeParser).CalledLookUpFiles {
		t.Error("unexpected existing files checking")
	}

	if f.out.(*shell.FakeOutputWriter).CalledWarning {
		t.Error("unexpected existing files Warning")
	}

	if f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("unexpected program Exit")
	}

	if !f.out.(*shell.FakeOutputWriter).CalledSuccess {
		t.Error("did not call Success")
	}
}

func TestWriteErrorInitCommand(t *testing.T) {
	f := newFakeKoolInit()
	f.parser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockError = errors.New("write error")

	cmd := NewInitCommand(f)

	cmd.SetArgs([]string{"laravel"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing init command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledError {
		t.Error("did not call Error")
	}

	expected := "Failed to write preset file : write error"
	output := f.out.(*shell.FakeOutputWriter).Err.Error()

	if output != expected {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}
