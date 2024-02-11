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

func fakeKoolDeploy(pwd string) *KoolDeploy {
	c := NewCloud()
	c.Fake()
	return &KoolDeploy{
		*(newDefaultKoolService().Fake()),
		c,
		setup.NewDefaultCloudSetupParser(pwd),
		&KoolCloudDeployFlags{},
		environment.NewFakeEnvStorage(),
		nil,
	}
}

func TestCreateReleaseFileNoConfig(t *testing.T) {
	fake := fakeKoolDeploy("")

	if _, err := fake.createReleaseFile(); err == nil || !strings.Contains(err.Error(), "no kool.cloud.yml config files found") {
		t.Errorf("expected error on createReleaseFile when no kool.deploy.yml exists in current working directory; got: %v", err)
	}
}

func TestCreateReleaseFileCreatesTgz(t *testing.T) {
	tmpDir := t.TempDir()
	mockConfig(tmpDir, t, nil)

	fake := fakeKoolDeploy(tmpDir)
	fake.env.Set("PWD", tmpDir)

	if err := fake.loadAndValidateConfig(); err != nil {
		t.Errorf("unexpected error on loadAndValidateConfig; got: %v", err)
	}

	if tg, err := fake.createReleaseFile(); err != nil {
		t.Errorf("unexpected error on createReleaseFile; got: %v", err)
	} else if _, err := os.Stat(tg); err != nil {
		t.Errorf("expected tgz file to be created; got: %v", err)
	}
}

func TestCreateReleaseFileCreatesTgzWithEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	mockConfig(tmpDir, t, []byte("services:\n  foo:\n    image: bar\n    env:\n      source: 'foo.env'\n"))

	if err := os.WriteFile(filepath.Join(tmpDir, "foo.env"), []byte("FOO=BAR"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fake := fakeKoolDeploy(tmpDir)
	fake.env.Set("PWD", tmpDir)

	if err := fake.loadAndValidateConfig(); err != nil {
		t.Errorf("unexpected error on loadAndValidateConfig; got: %v", err)
	}

	if tg, err := fake.createReleaseFile(); err != nil {
		t.Errorf("unexpected error on createReleaseFile; got: %v", err)
	} else if _, err := os.Stat(tg); err != nil {
		t.Errorf("expected tgz file to be created; got: %v", err)
	}

	if !fake.shell.(*shell.FakeShell).CalledPrintln {
		t.Error("expected Println to have been called on shell")
	} else if !strings.Contains(fake.shell.(*shell.FakeShell).OutLines[0], "Compressing files:") {
		t.Error("expected to print 'Compressing files:'")
	} else if !strings.Contains(fake.shell.(*shell.FakeShell).OutLines[1], "- "+filepath.Join(tmpDir, "kool.cloud.yml")) {
		t.Error("expected to print '- " + filepath.Join(tmpDir, "kool.cloud.yml") + "'")
	} else if !strings.Contains(fake.shell.(*shell.FakeShell).OutLines[2], "- "+filepath.Join(tmpDir, "docker-compose.yml")) {
		t.Error("expected to print '- " + filepath.Join(tmpDir, "docker-compose.yml") + "'")
	} else if !strings.Contains(fake.shell.(*shell.FakeShell).OutLines[3], "- "+filepath.Join(tmpDir, "foo.env")) {
		t.Error("expected to print '- " + filepath.Join(tmpDir, "foo.env") + "'")
	}
}

func TestCreateReleaseFileCreatesTgzWithEnvironmentFile(t *testing.T) {
	tmpDir := t.TempDir()
	mockConfig(tmpDir, t, []byte("services:\n  foo:\n    image: bar\n    environment: 'bar.env'\n"))

	if err := os.WriteFile(filepath.Join(tmpDir, "bar.env"), []byte("BAR=FOO"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fake := fakeKoolDeploy(tmpDir)
	fake.env.Set("PWD", tmpDir)

	if err := fake.loadAndValidateConfig(); err != nil {
		t.Errorf("unexpected error on loadAndValidateConfig; got: %v", err)
	}

	if tg, err := fake.createReleaseFile(); err != nil {
		t.Errorf("unexpected error on createReleaseFile; got: %v", err)
	} else if _, err := os.Stat(tg); err != nil {
		t.Errorf("expected tgz file to be created; got: %v", err)
	}

	if !fake.shell.(*shell.FakeShell).CalledPrintln {
		t.Error("expected Println to have been called on shell")
	} else if !strings.Contains(fake.shell.(*shell.FakeShell).OutLines[0], "Compressing files:") {
		t.Error("expected to print 'Compressing files:'")
	} else if !strings.Contains(fake.shell.(*shell.FakeShell).OutLines[1], "- "+filepath.Join(tmpDir, "kool.cloud.yml")) {
		t.Error("expected to print '- " + filepath.Join(tmpDir, "kool.cloud.yml") + "'")
	} else if !strings.Contains(fake.shell.(*shell.FakeShell).OutLines[2], "- "+filepath.Join(tmpDir, "docker-compose.yml")) {
		t.Error("expected to print '- " + filepath.Join(tmpDir, "docker-compose.yml") + "'")
	} else if !strings.Contains(fake.shell.(*shell.FakeShell).OutLines[3], "- "+filepath.Join(tmpDir, "bar.env")) {
		t.Error("expected to print '- " + filepath.Join(tmpDir, "bar.env") + "'")
	}
}

func TestCleanupReleaseFile(t *testing.T) {
	tmpDir := t.TempDir()
	mockConfig(tmpDir, t, nil)

	fake := fakeKoolDeploy("")
	fake.env.Set("PWD", tmpDir)

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
	fake := fakeKoolDeploy("")

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
