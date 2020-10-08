package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"io/ioutil"
	"kool-dev/kool/environment"
	"sort"
	"strings"
	"testing"
)

const testingEnv string = `
KOOL_FILTER_TESTING=1
KOOL_TESTING=1
`

func setup(f *KoolInfo) {
	f.envStorage.Set("KOOL_FILTER_TESTING", "1")
	f.envStorage.Set("KOOL_TESTING", "1")
}

func TestInfo(t *testing.T) {
	f := &KoolInfo{
		*newDefaultKoolService(),
		environment.NewFakeEnvStorage(),
	}

	setup(f)

	output, err := execInfoCommand(NewInfoCmd(f))

	if err != nil {
		t.Fatal(err)
	}

	expected := strings.Trim(testingEnv, "\n")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFilteredInfo(t *testing.T) {
	f := &KoolInfo{
		*newDefaultKoolService(),
		environment.NewFakeEnvStorage(),
	}

	setup(f)

	cmd := NewInfoCmd(f)
	cmd.SetArgs([]string{"FILTER"})

	output, err := execInfoCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := "KOOL_FILTER_TESTING=1"

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func execInfoCommand(cmd *cobra.Command) (output string, err error) {
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	if err = cmd.Execute(); err != nil {
		return
	}

	var out []byte
	if out, err = ioutil.ReadAll(b); err != nil {
		return
	}

	envs := strings.Split(strings.Trim(string(out), "\n"), "\n")
	sort.Strings(envs)

	output = strings.Join(envs, "\n")
	return
}
