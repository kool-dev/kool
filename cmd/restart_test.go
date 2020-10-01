package cmd

import (
	"errors"
	"kool-dev/kool/cmd/shell"
	"testing"
)

// NewRestartCommand initializes new kool start command
func TestRestartCommand(t *testing.T) {
	fake := &KoolRestart{
		*newFakeKoolService(),
		&FakeKoolService{},
		&FakeKoolService{},
	}

	cmd := NewRestartCommand(fake)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if !fake.stop.(*FakeKoolService).CalledExecute {
		t.Errorf("restart command did not call Execute on stop service")
	}

	if !fake.start.(*FakeKoolService).CalledExecute {
		t.Errorf("restart command did not call Execute on start service")
	}
}

func TestFailingStartRestartCommand(t *testing.T) {
	fake := &KoolRestart{
		*newFakeKoolService(),
		&FakeKoolService{},
		&FakeKoolService{},
	}

	fake.start.(*FakeKoolService).MockExecError = errors.New("start error")

	cmd := NewRestartCommand(fake)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if !fake.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not exit command due to error on start service")
	}

	if err := fake.out.(*shell.FakeOutputWriter).Err; err.Error() != "start error" {
		t.Errorf("expecting error 'start error', got '%v'", err)
	}
}

func TestFailingStopRestartCommand(t *testing.T) {
	fake := &KoolRestart{
		*newFakeKoolService(),
		&FakeKoolService{},
		&FakeKoolService{},
	}

	fake.stop.(*FakeKoolService).MockExecError = errors.New("stop error")

	cmd := NewRestartCommand(fake)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing restart command; error: %v", err)
	}

	if fake.start.(*FakeKoolService).CalledExecute {
		t.Errorf("restart should not call Execute on start service after failing stop service")
	}

	if !fake.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not exit command due to error on stop service")
	}

	if err := fake.out.(*shell.FakeOutputWriter).Err; err.Error() != "stop error" {
		t.Errorf("expecting error 'stop error', got '%v'", err)
	}
}
