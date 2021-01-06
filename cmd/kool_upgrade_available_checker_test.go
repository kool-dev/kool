package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/cmd/updater"
	"kool-dev/kool/environment"

	"errors"
	"fmt"
	"testing"
)

func newFakeUpdateAwareService(start *KoolStart, koolFakeUpdater *updater.FakeUpdater) *UpdateAwareService {
	return &UpdateAwareService{
		start,
		koolFakeUpdater,
	}
}

func TestStartWithUpdaterWrapper(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
	}

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "0.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    nil,
	}

	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.Run = DefaultCommandRunFunction(fakeUpdateAwareService)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledGetCurrentVersion {
		t.Errorf("did not call GetCurrentVersion")
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledCheckForUpdates {
		t.Errorf("did not call CheckForUpdates")
	}

	expected := "There's a new version available! Run kool self-update to upgrade!"

	if output := fmt.Sprint(koolStart.shell.(*shell.FakeShell).WarningOutput...); output != expected {
		t.Errorf("expecting warning message '%s', got '%s'", expected, output)
	}
}

func TestStartWithUpdaterWrapperError(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
	}

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "0.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    errors.New("error"),
	}

	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.Run = DefaultCommandRunFunction(fakeUpdateAwareService)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledGetCurrentVersion {
		t.Errorf("did not call GetCurrentVersion")
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledCheckForUpdates {
		t.Errorf("did not call CheckForUpdates")
	}

	if koolStart.exiter.(*shell.FakeExiter).Code() != 0 {
		t.Errorf("did not expect KoolStart service to have exit code different than 0; got '%d", koolStart.exiter.(*shell.FakeExiter).Code())
	}
}

func TestStartWithUpdaterWrapperSameVersion(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
	}

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "1.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    nil,
	}

	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.Run = DefaultCommandRunFunction(fakeUpdateAwareService)

	if _, err := execStartCommand(cmd); err != nil {
		t.Fatal(err)
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledGetCurrentVersion {
		t.Errorf("did not call GetCurrentVersion")
	}

	if !fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledCheckForUpdates {
		t.Errorf("did not call CheckForUpdates")
	}

	if koolStart.exiter.(*shell.FakeExiter).Code() != 0 {
		t.Errorf("did not expect KoolStart service to have exit code different than 0; got '%d", koolStart.exiter.(*shell.FakeExiter).Code())
	}
}

func TestDontCheckForUpdatesWhenNonTerminal(t *testing.T) {
	koolStart := &KoolStart{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockCmd: "start"},
	}

	koolUpdater := &updater.FakeUpdater{
		MockCurrentVersion: "0.0.0",
		MockLatestVersion:  "1.0.0",
		MockErrorUpdate:    nil,
	}
	koolStart.term.(*shell.FakeTerminalChecker).MockIsTerminal = false

	cmd := NewStartCommand(koolStart)
	fakeUpdateAwareService := newFakeUpdateAwareService(koolStart, koolUpdater)

	cmd.Run = DefaultCommandRunFunction(fakeUpdateAwareService)

	if fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledGetCurrentVersion {
		t.Errorf("called GetCurrentVersion")
	}

	if fakeUpdateAwareService.updater.(*updater.FakeUpdater).CalledCheckForUpdates {
		t.Errorf("called CheckForUpdates")
	}

	if koolStart.exiter.(*shell.FakeExiter).Code() != 0 {
		t.Errorf("did not expect KoolStart service to have exit code different than 0; got '%d", koolStart.exiter.(*shell.FakeExiter).Code())
	}
}
