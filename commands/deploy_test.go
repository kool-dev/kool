package commands

import (
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewKoolDeploy(t *testing.T) {
	kd := NewKoolDeploy()

	if _, is := kd.env.(*environment.DefaultEnvStorage); !is {
		t.Error("failed asserting default env storage")
	}

	if _, is := kd.git.(*builder.DefaultCommand); !is {
		t.Error("failed asserting default git command")
	}
}

func fakeKoolDeploy() *KoolDeploy {
	return &KoolDeploy{
		*(newDefaultKoolService().Fake()),
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
	}
}

func TestHandleDeployEnv(t *testing.T) {
	fake := fakeKoolDeploy()

	files := []string{}

	tmpDir := t.TempDir()
	fake.env.Set("PWD", tmpDir)

	files = fake.handleDeployEnv(files)

	if len(files) != 0 {
		t.Errorf("expected files to continue empty - no kool.deploy.env exists")
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "kool.deploy.env"), []byte("FOO=BAR"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	files = fake.handleDeployEnv(files)

	if len(files) != 1 {
		t.Errorf("expected files to have added kool.deploy.env")
	}

	files = fake.handleDeployEnv(files)

	if len(files) != 1 {
		t.Errorf("expected files to continue since was already there kool.deploy.env")
	}
}

func TestValidate(t *testing.T) {
	fake := fakeKoolDeploy()

	tmpDir := t.TempDir()
	fake.env.Set("PWD", tmpDir)

	if err := fake.validate(); err == nil || !strings.Contains(err.Error(), "could not find required file") {
		t.Error("failed getting proper error out of validate when no kool.deploy.yml exists in current working directory")
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "kool.deploy.yml"), []byte("services:\n"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := fake.validate(); err != nil {
		t.Errorf("unexpcted error on validate when file exists: %v", err)
	}
}

func TestParseFilesListFromGIT(t *testing.T) {
	fake := fakeKoolDeploy()

	if files, err := fake.parseFilesListFromGIT([]string{}); err != nil {
		t.Errorf("unexpected error from parseFileListFromGIT: %v", err)
	} else if len(files) != 0 {
		t.Errorf("unexpected return of files: %#v", files)
	}

	fake.git.(*builder.FakeCommand).MockExecOut = strings.Join([]string{"foo", string(rune(0x00)), "bar"}, "")

	if files, err := fake.parseFilesListFromGIT([]string{}); err != nil {
		t.Errorf("unexpected error from parseFileListFromGIT: %v", err)
	} else if len(files) != 2 {
		t.Errorf("unexpected return of files: %#v", files)
	}

	fake.git.(*builder.FakeCommand).MockExecError = errors.New("error")

	if _, err := fake.parseFilesListFromGIT([]string{"foo", "bar"}); err == nil || !strings.Contains(err.Error(), "failed listing GIT") {
		t.Errorf("unexpected error from parseFileListFromGIT: %v", err)
	}
}
