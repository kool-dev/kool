package environment

import (
	"os"
	"strings"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
)

const defaultEnvTesting string = `
VAR_DEFAULT_ENV=1
`

func TestInitEnvironmentVariables(t *testing.T) {
	f := NewFakeEnvStorage()

	envFile = ".env.testing"

	defer func() { envFile = ".env" }()

	InitEnvironmentVariables(f, defaultEnvTesting)

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

	if defaultEnv := f.Envs["VAR_DEFAULT_ENV"]; defaultEnv != "1" {
		t.Errorf("expecting $VAR_DEFAULT_ENV value '1', got '%s'", defaultEnv)
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
