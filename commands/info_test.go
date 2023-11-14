package commands

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupInfoTest(f *KoolInfo) {
	f.envStorage.Set("KOOL_FILTER_TESTING", "1")
	f.envStorage.Set("KOOL_TESTING", "1")
}

func fakeKoolInfo() *KoolInfo {
	return &KoolInfo{
		*(newDefaultKoolService().Fake()),
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
		&builder.FakeCommand{},
	}
}

func TestInfo(t *testing.T) {
	f := fakeKoolInfo()

	setupInfoTest(f)

	output, err := execInfoCommand(NewInfoCmd(f), f)

	if err != nil {
		t.Fatal(err)
	}

	for _, expected := range []string{"KOOL_FILTER_TESTING=1", "KOOL_TESTING=1"} {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected '%s', got '%s'", expected, output)
		}
	}
}

func TestFilteredInfo(t *testing.T) {
	f := fakeKoolInfo()

	setupInfoTest(f)

	cmd := NewInfoCmd(f)
	cmd.SetArgs([]string{"FILTER"})

	output, err := execInfoCommand(cmd, f)

	if err != nil {
		t.Fatal(err)
	}

	expected := "KOOL_FILTER_TESTING=1"

	if !strings.Contains(output, expected) {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func execInfoCommand(cmd *cobra.Command, f *KoolInfo) (output string, err error) {
	if err = cmd.Execute(); err != nil {
		return
	}

	output = strings.Join(f.shell.(*shell.FakeShell).OutLines, "\n")
	return
}
