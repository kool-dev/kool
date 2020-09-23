package cmd

import (
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
