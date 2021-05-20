// +build !windows

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"os"
	"strings"
	"testing"

	"github.com/creack/pty"
	"github.com/spf13/cobra"
)

func assertServiceAfterExecutingDefaultRun(service *FakeKoolService) (errMessage string) {
	if !service.CalledSetOutStream {
		errMessage = "did not call SetOutStream on kool service"
		return
	}

	if !service.CalledSetInStream {
		errMessage = "did not call SetInStream on kool service"
		return
	}

	if !service.CalledSetErrStream {
		errMessage = "did not call SetErrStream on kool service"
		return
	}

	if !service.CalledExecute {
		errMessage = "did not call Execute on kool service"
		return
	}

	return
}

func assertFailingServiceAfterExecutingDefaultRun(service *FakeKoolService) (errMessage string) {
	if !service.CalledSetOutStream {
		errMessage = "did not call SetOutStream on kool service"
		return
	}

	if !service.CalledSetInStream {
		errMessage = "did not call SetInStream on kool service"
		return
	}

	if !service.CalledSetErrStream {
		errMessage = "did not call SetErrStream on kool service"
		return
	}

	if !service.CalledExecute {
		errMessage = "did not call Execute on kool service"
		return
	}

	if !service.CalledError {
		errMessage = "did not call Error on kool service"
		return
	}

	if !service.CalledExit {
		errMessage = "did not call Exit on kool service"
		return
	}

	if service.ExitCode != 1 {
		errMessage = fmt.Sprintf("should exit with status 1, got %v", service.ExitCode)
		return
	}

	return
}

func TestNewRootCmd(t *testing.T) {
	fakeEnv := environment.NewFakeEnvStorage()
	cmd := NewRootCmd(fakeEnv)

	if cmd.Name() != rootCmd.Name() {
		t.Errorf("expecting RootCmd to return '%s', got '%s'", rootCmd.Name(), cmd.Name())
	}
}

func TestRootCmd(t *testing.T) {
	cmd := RootCmd()

	if cmd.Name() != rootCmd.Name() {
		t.Errorf("expecting RootCmd to return '%s', got '%s'", rootCmd.Name(), cmd.Name())
	}
}

func TestRootCmdExecute(t *testing.T) {
	_, w, err := os.Pipe()

	if err != nil {
		t.Fatal(err)
	}

	originalOutput := os.Stdout
	os.Stdout = w

	defer func(originalOutput *os.File) {
		os.Stdout = originalOutput
	}(originalOutput)

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

	if out, err = io.ReadAll(b); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("kool version %s", version)
	output := strings.TrimSpace(string(out))

	if expected != output {
		t.Errorf("expecting rootCmd with Version Flag to return '%s', got '%s'", expected, output)
	}
}

func TestDefaultCommandRunFunction(t *testing.T) {
	f := &FakeKoolService{}

	cmd := &cobra.Command{
		Use:   "fake-command",
		Short: "fake - fake command",
		Run:   DefaultCommandRunFunction(f),
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing root command; error: %v", err)
	}

	if errMessage := assertServiceAfterExecutingDefaultRun(f); errMessage != "" {
		t.Error(errMessage)
	}
}

func TestFailingDefaultCommandRunFunction(t *testing.T) {
	f := &FakeKoolService{MockExecError: fmt.Errorf("execute error")}

	cmd := &cobra.Command{
		Use:   "fake-command",
		Short: "fake - fake command",
		Run:   DefaultCommandRunFunction(f),
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing root command; error: %v", err)
	}

	if errMessage := assertFailingServiceAfterExecutingDefaultRun(f); errMessage != "" {
		t.Error(errMessage)
	}
}

func TestMultipleServicesDefaultCommandRunFunction(t *testing.T) {
	var services []*FakeKoolService

	services = append(services, &FakeKoolService{})
	services = append(services, &FakeKoolService{})

	cmd := &cobra.Command{
		Use:   "fake-command",
		Short: "fake - fake command",
		Run:   DefaultCommandRunFunction(services[0], services[1]),
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing root command; error: %v", err)
	}

	for _, service := range services {
		if errMessage := assertServiceAfterExecutingDefaultRun(service); errMessage != "" {
			t.Error(errMessage)
		}
	}
}

func TestMultipleServicesFailingDefaultCommandRunFunction(t *testing.T) {
	failing := &FakeKoolService{MockExecError: fmt.Errorf("execute error")}
	passing := &FakeKoolService{}

	cmd := &cobra.Command{
		Use:   "fake-command",
		Short: "fake - fake command",
		Run:   DefaultCommandRunFunction(failing, passing),
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing root command; error: %v", err)
	}

	if errMessage := assertFailingServiceAfterExecutingDefaultRun(failing); errMessage != "" {
		t.Error(errMessage)
	}
}

func TestVerboseFlagRootCommand(t *testing.T) {
	fakeEnv := environment.NewFakeEnvStorage()

	fInfo := &KoolInfo{
		*newFakeKoolService(),
		fakeEnv,
	}

	root := NewRootCmd(fakeEnv)
	info := NewInfoCmd(fInfo)
	root.AddCommand(info)

	root.SetArgs([]string{"--verbose", "info"})

	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error executing command; error: %v", err)
	}

	if verbose := fakeEnv.IsTrue("KOOL_VERBOSE"); !verbose {
		t.Error("expecting 'KOOL_VERBOSE' to be true, got false")
	}
}

func TestRecursiveCall(t *testing.T) {
	recursive := &cobra.Command{
		Use: "recursive",
		Run: func(cmd *cobra.Command, args []string) {
			_ = shell.Interactive("kool", "-v")
		},
	}

	rootCmd.AddCommand(recursive)

	rootCmd.SetArgs([]string{"recursive"})

	err := Execute()

	if err != nil {
		t.Errorf("fail calling recursive command: %v", err)
	}
}

func TestMultipleRecursiveCall(t *testing.T) {
	recursive := &cobra.Command{
		Use: "recursive",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = shell.Interactive("kool", "-v"); err != nil {
				return
			}
			err = shell.Interactive("kool", "-v")
			return
		},
	}

	rootCmd.AddCommand(recursive)

	rootCmd.SetArgs([]string{"recursive"})

	err := Execute()

	if err != nil {
		t.Errorf("fail calling recursive command: %v", err)
	}
}

func TestAddCommands(t *testing.T) {
	root := NewRootCmd(environment.NewFakeEnvStorage())

	AddCommands(root)

	var subcommands map[string]bool = map[string]bool{
		"completion":  false,
		"create":      false,
		"deploy":      false,
		"docker":      false,
		"exec":        false,
		"info":        false,
		"init":        false,
		"logs":        false,
		"preset":      false,
		"restart":     false,
		"run":         false,
		"self-update": false,
		"share":       false,
		"start":       false,
		"status":      false,
		"stop":        false,
	}

	for _, subCmd := range root.Commands() {
		name := subCmd.Name()
		if _, ok := subcommands[name]; !ok {
			t.Errorf("unexpected command was added: %s", name)
			continue
		}

		subcommands[name] = true
	}

	for cmd, added := range subcommands {
		if !added {
			t.Errorf("expected command is missing: %s", cmd)
		}
	}
}

func TestDevelopmentVersionWarning(t *testing.T) {
	fakeEnv := environment.NewFakeEnvStorage()
	root := NewRootCmd(fakeEnv)

	fakecmd := &cobra.Command{
		Use: "fakecmd",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	root.AddCommand(fakecmd)

	// default test NOT A TTY
	b := bytes.NewBufferString("")
	root.SetOut(b)

	root.SetArgs([]string{"fakecmd"})
	version = DEV_VERSION
	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error executing command; error: %v", err)
	}

	var (
		out []byte
		err error
	)

	if out, err = io.ReadAll(b); err != nil {
		t.Fatal(err)
	}

	expected := "you are executing a development version"
	output := strings.TrimSpace(string(out))

	if strings.Contains(output, expected) {
		t.Errorf("bad warning under non-TTY: %s", output)
	}

	if hasWarnedDevelopmentVersion {
		t.Error("bar warning under non-TTY")
	}

	if pty, tty, err := pty.Open(); err != nil {
		t.Fatalf("failed creting PTY for testing: %v", err)
	} else {
		root.SetOut(tty)

		defer pty.Close()
		defer tty.Close()
	}
	version = DEV_VERSION
	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error executing command; error: %v", err)
	}

	if !hasWarnedDevelopmentVersion {
		t.Error("failed to warn about development version")
	}

	hasWarnedDevelopmentVersion = false
	version = "100.100.100"
	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error executing command; error: %v", err)
	}

	if hasWarnedDevelopmentVersion {
		t.Error("should not have warned on non-dev version")
	}
}
