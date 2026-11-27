package commands

import (
	"errors"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/api"
	"kool-dev/kool/services/cloud/setup"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"
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

// fakeDeployStatuses returns a fetch function that replies with the given
// statuses in order; an empty status means a failed status check
func fakeDeployStatuses(statuses ...string) (fetch func() (*api.DeployStatusResponse, error), calls *int) {
	calls = new(int)
	fetch = func() (*api.DeployStatusResponse, error) {
		status := statuses[*calls]
		*calls++
		if status == "" {
			return nil, errors.New("\n503 - unavailable\n")
		}
		return &api.DeployStatusResponse{Status: status, Error: "error for " + status}, nil
	}
	return
}

func TestPollDeployStatusFinalStates(t *testing.T) {
	for _, final := range []string{"success", "failed", "timedout"} {
		fake := fakeKoolDeploy("")
		fetch, calls := fakeDeployStatuses("building", "deploying", "deploying", final)
		var waits []time.Duration

		status, err := fake.pollDeployStatus(fetch, func(d time.Duration) { waits = append(waits, d) })

		if err != nil {
			t.Errorf("unexpected error polling until %s: %v", final, err)
		}
		if status == nil || status.Status != final {
			t.Errorf("expected final status %s; got %+v", final, status)
		}
		if *calls != 4 {
			t.Errorf("expected 4 status checks until %s; got %d", final, *calls)
		}
		if expected := []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second}; !reflect.DeepEqual(waits, expected) {
			t.Errorf("expected waits %v; got %v", expected, waits)
		}

		lines := fake.shell.(*shell.FakeShell).OutLines
		if expected := []string{"> deploy: building", "> deploy: deploying", "> deploy: " + final}; !reflect.DeepEqual(lines, expected) {
			t.Errorf("expected status changes %v; got %v", expected, lines)
		}
	}
}

func TestPollDeployStatusRetriesFailedChecks(t *testing.T) {
	fake := fakeKoolDeploy("")
	fetch, calls := fakeDeployStatuses("deploying", "", "", "", "deploying", "", "failed")
	var waits []time.Duration

	status, err := fake.pollDeployStatus(fetch, func(d time.Duration) { waits = append(waits, d) })

	if err != nil {
		t.Errorf("unexpected error when status checks recover: %v", err)
	}
	if status == nil || status.Status != "failed" {
		t.Errorf("expected failed status; got %+v", status)
	}
	if *calls != 7 {
		t.Errorf("expected 7 status checks; got %d", *calls)
	}

	// backoff doubles while checks fail and resets after one succeeds
	expected := []time.Duration{2 * time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 2 * time.Second, 2 * time.Second}
	if !reflect.DeepEqual(waits, expected) {
		t.Errorf("expected waits %v; got %v", expected, waits)
	}

	fakeShell := fake.shell.(*shell.FakeShell)
	if !fakeShell.CalledWarning {
		t.Error("expected a warning for failed status checks")
	} else if warning := fakeShell.WarningOutput[0].(string); !strings.Contains(warning, "attempt 1 of 5") || !strings.Contains(warning, "503 - unavailable") {
		t.Errorf("unexpected warning for failed status check: %q", warning)
	}
}

func TestPollDeployStatusGivesUp(t *testing.T) {
	fake := fakeKoolDeploy("")
	fetch, calls := fakeDeployStatuses("deploying", "", "", "", "", "", "success")
	var waits []time.Duration

	status, err := fake.pollDeployStatus(fetch, func(d time.Duration) { waits = append(waits, d) })

	if err == nil || err.Error() != "failed checking deploy status 5 times in a row: 503 - unavailable" {
		t.Errorf("unexpected error after failed status checks: %v", err)
	}
	if status == nil || status.Status != "deploying" {
		t.Errorf("expected last known status; got %+v", status)
	}
	if *calls != 6 {
		t.Errorf("expected 6 status checks; got %d", *calls)
	}

	expected := []time.Duration{2 * time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 10 * time.Second}
	if !reflect.DeepEqual(waits, expected) {
		t.Errorf("expected waits %v; got %v", expected, waits)
	}
}

func TestDeployFailedError(t *testing.T) {
	tests := []struct {
		status   api.DeployStatusResponse
		logsURL  string
		expected string
	}{
		{api.DeployStatusResponse{Status: "failed"}, "", "deploy failed"},
		{api.DeployStatusResponse{Status: "failed"}, "https://kool.dev/deploys/1", "deploy failed\nLogs available at: https://kool.dev/deploys/1"},
		{api.DeployStatusResponse{Status: "failed", Error: " pods not ready \n"}, "https://kool.dev/deploys/1", "deploy failed: pods not ready\nLogs available at: https://kool.dev/deploys/1"},
		{api.DeployStatusResponse{Status: "timedout", Error: "build took too long"}, "", "deploy timed out: build took too long"},
		{
			api.DeployStatusResponse{Status: "failed", Error: "migrate hook failed", Hint: "check your migrations", Detail: "SQLSTATE[42S01]"},
			"https://kool.dev/deploys/1",
			"deploy failed: migrate hook failed\nHint: check your migrations\nDetail: SQLSTATE[42S01]\nLogs available at: https://kool.dev/deploys/1",
		},
		{
			// credentials quoted by the diagnostics never reach the output
			api.DeployStatusResponse{Status: "failed", Error: "connect mysql://app:hunter2@db/app failed", Hint: "check DB_PASSWORD='hunter2'", Detail: "DB_PASSWORD=hunter2 DB_HOST=db"},
			"https://kool.dev/deploys/1",
			"deploy failed: connect mysql://app:[REDACTED]@db/app failed\nHint: check DB_PASSWORD='[REDACTED]'\nDetail: DB_PASSWORD=[REDACTED] DB_HOST=db\nLogs available at: https://kool.dev/deploys/1",
		},
	}

	for _, tt := range tests {
		if err := deployFailedError(&tt.status, tt.logsURL); err.Error() != tt.expected {
			t.Errorf("expected error %q; got %q", tt.expected, err.Error())
		}
	}
}

func TestWithDeployLogs(t *testing.T) {
	timeoutErr := errors.New("timeout waiting deploy to finish")

	if err := withDeployLogs(timeoutErr, ""); err != timeoutErr {
		t.Errorf("expected the same error without logs URL; got %v", err)
	}

	err := withDeployLogs(timeoutErr, "https://kool.dev/deploys/1")
	if err.Error() != "timeout waiting deploy to finish\nLogs available at: https://kool.dev/deploys/1" {
		t.Errorf("unexpected error with logs URL: %v", err)
	}
	if !errors.Is(err, timeoutErr) {
		t.Error("expected the error with logs URL to wrap the original error")
	}
}

func fakeStartedDeploy() *api.DeployCreateResponse {
	deploy := &api.DeployCreateResponse{LogsUrl: "https://kool.dev/deploys/44417"}
	deploy.Deploy.ID = 44417
	deploy.Deploy.Url = "https://app.example.com"
	deploy.Deploy.Environment.Name = "app"
	return deploy
}

func TestWaitDeployInterrupted(t *testing.T) {
	for sig, code := range map[os.Signal]int{os.Interrupt: 130, syscall.SIGTERM: 143} {
		fake := fakeKoolDeploy("")
		interrupted := make(chan os.Signal, 1)
		interrupted <- sig

		var stopped, finished bool
		err := fake.waitDeploy(fakeStartedDeploy(), make(chan deployWaitResult), make(chan time.Time), interrupted, func(f bool) {
			stopped, finished = true, f
		})

		if !stopped || finished {
			t.Errorf("expected the loading indicator stopped as not finished; got stopped=%v finished=%v", stopped, finished)
		}

		var exitable shell.ErrExitable
		if !errors.As(err, &exitable) || exitable.Code != code {
			t.Errorf("expected exit code %d on %s; got %v", code, sig, err)
			continue
		}

		for _, expected := range []string{
			"while waiting for deploy 44417 to finish",
			"keeps rolling out on Kool.dev Cloud and was neither cancelled nor reported as failed",
			"https://kool.dev/deploys/44417",
		} {
			if !strings.Contains(err.Error(), expected) {
				t.Errorf("expected interrupt message to contain %q; got %q", expected, err.Error())
			}
		}

		if strings.Contains(err.Error(), "deploy failed") {
			t.Errorf("an interrupt must not read as a failed deploy; got %q", err.Error())
		}
	}
}

func TestWaitDeployOutcomes(t *testing.T) {
	fake := fakeKoolDeploy("")
	finishes := make(chan deployWaitResult, 1)
	finishes <- deployWaitResult{status: &api.DeployStatusResponse{Status: "success"}}

	var finished bool
	if err := fake.waitDeploy(fakeStartedDeploy(), finishes, make(chan time.Time), make(chan os.Signal), func(f bool) { finished = f }); err != nil {
		t.Errorf("unexpected error for a successful deploy: %v", err)
	}
	if !finished || !fake.shell.(*shell.FakeShell).CalledSuccess {
		t.Error("expected the successful deploy to be reported as finished")
	}

	fake = fakeKoolDeploy("")
	finishes <- deployWaitResult{status: &api.DeployStatusResponse{Status: "failed", Error: "pods not ready"}}
	if err := fake.waitDeploy(fakeStartedDeploy(), finishes, make(chan time.Time), make(chan os.Signal), func(bool) {}); err == nil || !strings.HasPrefix(err.Error(), "deploy failed: pods not ready") {
		t.Errorf("expected failed deploy error; got %v", err)
	}

	timeout := make(chan time.Time, 1)
	timeout <- time.Now()
	finished = true
	if err := fake.waitDeploy(fakeStartedDeploy(), make(chan deployWaitResult), timeout, make(chan os.Signal), func(f bool) { finished = f }); err == nil || !strings.HasPrefix(err.Error(), "timeout waiting deploy to finish") {
		t.Errorf("expected timeout error; got %v", err)
	}
	if finished {
		t.Error("expected the loading indicator stopped as not finished on timeout")
	}
}
