package environment

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
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

	homeDir, _ := homedir.Dir()

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
