package commands

import (
	"kool-dev/kool/core/presets"
	"kool-dev/kool/core/shell"
	"testing"
)

func TestNewKoolPreset(t *testing.T) {
	k := NewKoolPreset()

	if _, ok := k.shell.(*shell.DefaultShell); !ok {
		t.Errorf("unexpected shell.Shell on default KoolPreset instance")
	}

	if _, ok := k.presetsParser.(*presets.DefaultParser); !ok {
		t.Errorf("unexpected presets.Parser on default KoolPreset instance")
	}

	if _, ok := k.promptSelect.(*shell.DefaultPromptSelect); !ok {
		t.Errorf("unexpected shell.PromptSelect on default KoolPreset instance")
	}
}

func TestPresetNoInteractionRequiresArgument(t *testing.T) {
	f := NewKoolPreset()
	f.shell = &shell.FakeShell{MockIsTerminal: true}
	prompt := &shell.FakePromptSelect{}
	f.promptSelect = prompt
	cmd := NewPresetCommand(f)
	cmd.SetArgs([]string{"--no-interaction"})
	assertExecGotError(t, cmd, "please specify a preset")
	if prompt.CalledAsk {
		t.Fatal("non-interactive preset prompted for selection")
	}
}
