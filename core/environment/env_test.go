package environment

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitEnvironmentVariables(t *testing.T) {
	f := NewFakeEnvStorage()

	originalEnvFiles := envFiles
	defer func() { envFiles = originalEnvFiles }()

	testEnvFile := filepath.Join(t.TempDir(), ".env.test")
	envFiles = []string{testEnvFile}

	if err := os.WriteFile(testEnvFile, []byte("FOO=bar\n"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	InitEnvironmentVariables(f)

	homeDir, _ := os.UserHomeDir()

	if envHomeDir := f.Envs["HOME"]; envHomeDir != homeDir {
		t.Errorf("expecting $HOME value '%s', got '%s'", homeDir, envHomeDir)
	}

	UID := uid()

	if envUID := f.Envs["UID"]; envUID != UID {
		t.Errorf("expecting $UID value '%s', got '%s'", UID, envUID)
	}

	workDir, _ := os.Getwd()

	if envWorkDir := f.Envs["PWD"]; envWorkDir != workDir {
		t.Errorf("expecting $PWD value '%s', got '%s'", workDir, envWorkDir)
	}

	if !f.CalledLoad {
		t.Error("did not call Load on EnvSotrage")
	}

	if foo := f.Envs["FOO"]; foo != "bar" {
		t.Errorf("expected FOO to be bar: %v", foo)
	}

	pieces := strings.Split(workDir, string(os.PathSeparator))
	koolName := pieces[len(pieces)-1]

	if envKoolName := f.Envs["KOOL_NAME"]; envKoolName != koolName {
		t.Errorf("expecting $KOOL_NAME value '%s', got '%s'", koolName, envKoolName)
	}

	if envKoolNet := f.Envs["KOOL_GLOBAL_NETWORK"]; envKoolNet != "kool_global" {
		t.Errorf("expecting $KOOL_GLOBAL_NETWORK value 'kool_global', got '%s'", envKoolNet)
	}
}

func TestInitEnvironmentVariablesOverridesStalePWD(t *testing.T) {
	f := NewFakeEnvStorage()

	// Simulate a stale PWD value (as might happen when spawned with cwd option)
	f.Envs["PWD"] = "/some/stale/path"

	originalEnvFiles := envFiles
	defer func() { envFiles = originalEnvFiles }()
	envFiles = []string{} // no env files needed for this test

	InitEnvironmentVariables(f)

	workDir, _ := os.Getwd()

	// PWD should be overridden with the actual working directory
	if envWorkDir := f.Envs["PWD"]; envWorkDir != workDir {
		t.Errorf("expecting $PWD to be overridden to '%s', got '%s'", workDir, envWorkDir)
	}
}

func TestInitSourceProjectUsesComposeName(t *testing.T) {
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "compose.yml"), []byte("name: ${PROJECT_NAME:-custom-project}\nservices: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := NewFakeEnvStorage()

	initSourceProject(env, workDir)

	if got := env.Get("KOOL_WORKSPACE_SOURCE_PROJECT"); got != "custom-project" {
		t.Errorf("expected Compose name custom-project, got %q", got)
	}
}

func TestInitEnvironmentSkipsWorkspaceDetectionWithoutConfiguration(t *testing.T) {
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	workDir := t.TempDir()
	if err = os.Chdir(workDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalDirectory) })

	originalGitWorktreeOutput := gitWorktreeOutput
	gitCalled := false
	gitWorktreeOutput = func(string, ...string) ([]byte, error) {
		gitCalled = true
		return nil, nil
	}
	t.Cleanup(func() { gitWorktreeOutput = originalGitWorktreeOutput })

	env := NewFakeEnvStorage()
	InitEnvironmentVariables(env)

	if gitCalled {
		t.Error("expected Git workspace detection to be skipped")
	}
	if env.IsTrue("KOOL_WORKSPACES_ENABLED") || env.IsTrue("KOOL_PROXY_ENABLED") {
		t.Error("expected workspace and proxy features to remain disabled")
	}
}
