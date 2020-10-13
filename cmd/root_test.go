package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func assertServiceAfterExecutingDefaultRun(service *FakeKoolService) (errMessage string) {
	if !service.CalledSetWriter {
		errMessage = "did not call SetWriter on kool service"
		return
	}

	if !service.CalledSetReader {
		errMessage = "did not call SetReader on kool service"
		return
	}

	if !service.CalledExecute {
		errMessage = "did not call Execute on kool service"
		return
	}

	return
}

func assertFailingServiceAfterExecutingDefaultRun(service *FakeKoolService) (errMessage string) {
	if !service.CalledSetWriter {
		errMessage = "did not call SetWriter on kool service"
		return
	}

	if !service.CalledSetReader {
		errMessage = "did not call SetReader on kool service"
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

	if out, err = ioutil.ReadAll(b); err != nil {
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
