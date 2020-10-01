package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestRootCmd(t *testing.T) {
	cmd := RootCmd()

	if cmd.Name() != rootCmd.Name() {
		t.Errorf("expecting RootCmd to return '%s', got '%s'", rootCmd.Name(), cmd.Name())
	}
}

func TestRootCmdExecute(t *testing.T) {
	if err := Execute(); err != nil {
		t.Errorf("unexpected error executing root command; error: %v", err)
	}
}

func TestVersionFlagCommand(t *testing.T) {
	cmd := RootCmd()

	cmd.SetArgs([]string{"--version"})

	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing root command; error: %v", err)
	}

	var (
		out []byte
		err error
	)

	if out, err = ioutil.ReadAll(b); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("kool version %s", version)
	output := strings.TrimSpace(string(out))

	if expected != output {
		t.Errorf("expecting rootCmd with Version Flag to return '%s', got '%s'", expected, output)
	}
}
