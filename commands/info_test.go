package commands

import (
	"encoding/json"
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

func TestInfoJSONOutput(t *testing.T) {
	f := fakeKoolInfo()
	f.shell.(*shell.FakeShell).MockIsJSONOutput = true
	f.cmdDocker.(*builder.FakeCommand).MockExecOut = "Docker version 29.0.0"
	f.cmdDocker.(*builder.FakeCommand).MockCmd = "docker"
	f.cmdDockerCompose.(*builder.FakeCommand).MockExecOut = "Docker Compose version v2.30.0"
	f.cmdDockerCompose.(*builder.FakeCommand).MockCmd = "docker"

	setupInfoTest(f)
	f.envStorage.Set("KOOL_OUTPUT", "json")

	cmd := NewInfoCmd(f)

	output, err := execInfoCommand(cmd, f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var info infoOutputJSON
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		t.Fatalf("failed to parse json output: %v\nraw: %s", err, output)
	}

	if info.KoolVersion == "" {
		t.Error("expected non-empty kool_version")
	}
	if info.DockerVersion != "Docker version 29.0.0" {
		t.Errorf("expected docker_version 'Docker version 29.0.0', got '%s'", info.DockerVersion)
	}
	if info.DockerComposeVersion != "Docker Compose version v2.30.0" {
		t.Errorf("expected docker_compose_version, got '%s'", info.DockerComposeVersion)
	}
	if info.Env["KOOL_TESTING"] != "1" {
		t.Errorf("expected env KOOL_TESTING=1, got '%s'", info.Env["KOOL_TESTING"])
	}
	if info.Env["KOOL_OUTPUT"] != "json" {
		t.Errorf("expected env KOOL_OUTPUT=json, got '%s'", info.Env["KOOL_OUTPUT"])
	}
}

func TestInfoJSONOutputRedactsAPIToken(t *testing.T) {
	f := fakeKoolInfo()
	f.shell.(*shell.FakeShell).MockIsJSONOutput = true
	f.cmdDocker.(*builder.FakeCommand).MockExecOut = "Docker version 29.0.0"
	f.cmdDocker.(*builder.FakeCommand).MockCmd = "docker"
	f.cmdDockerCompose.(*builder.FakeCommand).MockExecOut = "Docker Compose version v2.30.0"
	f.cmdDockerCompose.(*builder.FakeCommand).MockCmd = "docker"

	f.envStorage.Set("KOOL_API_TOKEN", "super-secret-token")

	cmd := NewInfoCmd(f)

	output, err := execInfoCommand(cmd, f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(output, "super-secret-token") {
		t.Error("KOOL_API_TOKEN value should be redacted in JSON output")
	}

	var info infoOutputJSON
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		t.Fatalf("failed to parse json output: %v", err)
	}

	if info.Env["KOOL_API_TOKEN"] == "super-secret-token" {
		t.Error("KOOL_API_TOKEN should be redacted")
	}
}
