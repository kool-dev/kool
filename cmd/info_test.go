package cmd

import (
	"bytes"
	"github.com/fireworkweb/godotenv"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"
)

const testingEnv string = `
KOOL_FILTER_TESTING=1
KOOL_TESTING=1
`

func setup() {
	testingEnv, _ := godotenv.Unmarshal(testingEnv)
	for k, v := range testingEnv {
		os.Setenv(k, v)
	}
}

func TestInfo(t *testing.T) {
	setup()

	output, err := execCommand(NewInfoCmd())

	if err != nil {
		t.Fatal(err)
	}

	expected := strings.Trim(testingEnv, "\n")

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFilteredInfo(t *testing.T) {
	setup()

	cmd := NewInfoCmd()
	cmd.SetArgs([]string{"FILTER"})

	output, err := execCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	expected := "KOOL_FILTER_TESTING=1"

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func execCommand(cmd *cobra.Command) (output string, err error) {
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
