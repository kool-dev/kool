package commands

import (
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/setup"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewKoolDeploy(t *testing.T) {
	kd := NewKoolDeploy(NewCloud())

	if _, is := kd.env.(*environment.DefaultEnvStorage); !is {
		t.Error("failed asserting default env storage")
	}

	if _, is := kd.cloud.env.(*environment.DefaultEnvStorage); !is {
		t.Error("failed asserting default cloud.env storage")
	}

	if _, is := kd.setupParser.(*setup.DefaultCloudSetupParser); !is {
		t.Error("failed asserting default cloud setup parser")
	}
}

func fakeKoolDeploy() *KoolDeploy {
	c := NewCloud()
	c.Fake()
	return &KoolDeploy{
		*(newDefaultKoolService().Fake()),
		c,
		setup.NewDefaultCloudSetupParser(""),
		&KoolCloudDeployFlags{},
		environment.NewFakeEnvStorage(),
		nil,
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

func TestCreateReleaseFile(t *testing.T) {
	fake := fakeKoolDeploy()

	tmpDir := t.TempDir()
	fake.env.Set("PWD", tmpDir)

	if _, err := fake.createReleaseFile(); err == nil || !strings.Contains(err.Error(), "no deploy config files found") {
		t.Errorf("expected error on createReleaseFile when no kool.deploy.yml exists in current working directory; got: %v", err)
	}

	mockConfig(tmpDir, t, nil)

	if tg, err := fake.createReleaseFile(); err != nil {
		t.Errorf("unexpected error on createReleaseFile; got: %v", err)
	} else if _, err := os.Stat(tg); err != nil {
		t.Errorf("expected tgz file to be created; got: %v", err)
	}
}

func TestCleanupReleaseFile(t *testing.T) {
	fake := fakeKoolDeploy()

	tmpDir := t.TempDir()
	fake.env.Set("PWD", tmpDir)

	mockConfig(tmpDir, t, nil)

	f := filepath.Join(tmpDir, "kool.cloud.yml")
	fake.cleanupReleaseFile(f)
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Errorf("expected file to be removed")
	}

	fake.cleanupReleaseFile(f)
	if !fake.shell.(*shell.FakeShell).CalledError {
		t.Errorf("expected for Error to have been called on shell")
	}
	if !strings.Contains(fake.shell.(*shell.FakeShell).Err.Error(), "error trying to remove temporary tarball") {
		t.Errorf("expected to print proper error message if file removal fails")
	}
}

func TestLoadAndValidateConfig(t *testing.T) {
	fake := fakeKoolDeploy()

	tmpDir := t.TempDir()
	fake.env.Set("PWD", tmpDir)

	if err := fake.loadAndValidateConfig(); err == nil || !strings.Contains(err.Error(), "could not find required file") {
		t.Error("failed getting proper error out of loadAndValidateConfig when no kool.cloud.yml exists in current working directory")
	}

	mockConfig(tmpDir, t, []byte("services:\n\tfoo:\n"))

	if err := fake.loadAndValidateConfig(); err == nil || !strings.Contains(err.Error(), "found character that cannot start") {
		t.Errorf("unexpcted error on loadAndValidateConfig with bad config: %v", err)
	}

	mockConfig(tmpDir, t, nil)

	if err := fake.loadAndValidateConfig(); err != nil {
		t.Errorf("unexpcted error on loadAndValidateConfig when file exists: %v", err)
	}

	if fake.cloudConfig.Cloud.Services == nil {
		t.Error("failed loading cloud config")
	}

	if len(fake.cloudConfig.Cloud.Services) != 1 {
		t.Error("service count mismatch - should be 1")
	} else if *fake.cloudConfig.Cloud.Services["foo"].Image != "bar" {
		t.Error("failed loading service foo image 'bar'")
	}
}

func mockConfig(tmpDir string, t *testing.T, mock []byte) {
	if mock == nil {
		mock = []byte("services:\n  foo:\n    image: bar\n")
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "kool.cloud.yml"), mock, os.ModePerm); err != nil {
		t.Fatal(err)
	}
}
