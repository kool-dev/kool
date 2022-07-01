package commands

import (
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/updater"

	"errors"
	"fmt"
	"testing"
)

func newFakeUpdateAwareService(start *KoolStart, koolFakeUpdater *updater.FakeUpdater) *UpdateAwareService {
	return &UpdateAwareService{
		start,
		koolFakeUpdater,
		false,
	}
}

func TestStartWithUpdaterWrapper(t *testing.T) {
	koolStart := newFakeKoolStart()

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "0.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    nil,
	}

	koolStart.Fake()
	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.RunE = DefaultCommandRunFunction(fakeUpdateAwareService)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if !koolUpdater.CalledGetCurrentVersion {
		t.Errorf("did not call GetCurrentVersion")
	}

	if !koolUpdater.CalledCheckForUpdates {
		t.Errorf("did not call CheckForUpdates")
	}

	expected := "There's a new version available! Run kool self-update to upgrade!"

	if output := fmt.Sprint(koolStart.shell.(*shell.FakeShell).WarningOutput...); output != expected {
		t.Errorf("expecting warning message '%s', got '%s'", expected, output)
	}
}

func TestStartWithUpdaterWrapperTimeout(t *testing.T) {
	koolStart := newFakeKoolStart()

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "0.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    nil,
		MockTimeoutDelay:   true,
	}

	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.RunE = DefaultCommandRunFunction(fakeUpdateAwareService)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	expected := ""

	if output := fmt.Sprint(koolStart.shell.(*shell.FakeShell).WarningOutput...); output != expected {
		t.Errorf("expecting warning message '%s', got '%s'", expected, output)
	}
}

func TestStartWithUpdaterWrapperError(t *testing.T) {
	koolStart := newFakeKoolStart()
	koolStart.Fake()

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "0.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    errors.New("error"),
	}

	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.RunE = DefaultCommandRunFunction(fakeUpdateAwareService)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledGetCurrentVersion {
		t.Errorf("did not call GetCurrentVersion")
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledCheckForUpdates {
		t.Errorf("did not call CheckForUpdates")
	}
}

func TestStartWithUpdaterWrapperSameVersion(t *testing.T) {
	koolStart := newFakeKoolStart()
	koolStart.Fake()

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "1.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    nil,
	}

	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.RunE = DefaultCommandRunFunction(fakeUpdateAwareService)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledGetCurrentVersion {
		t.Errorf("did not call GetCurrentVersion")
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledCheckForUpdates {
		t.Errorf("did not call CheckForUpdates")
	}
}

func TestDontCheckForUpdatesWhenNonTerminal(t *testing.T) {
	koolStart := newFakeKoolStart()

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "0.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    nil,
	}
	koolStart.shell.(*shell.FakeShell).MockIsTerminal = false

	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.RunE = DefaultCommandRunFunction(fakeUpdateAwareService)

	if err := cmd.Execute(); err != nil {
		t.Errorf("error %v", err)
	}

	if !koolStart.shell.(*shell.FakeShell).CalledIsTerminal {
		t.Error("should have called IsTerminal")
	}

	if fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledGetCurrentVersion {
		t.Errorf("called GetCurrentVersion")
	}

	if fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledCheckForUpdates {
		t.Errorf("called CheckForUpdates")
	}
}
