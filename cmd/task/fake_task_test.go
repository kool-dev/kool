package task

import "testing"

func TestFakeTask(t *testing.T) {
	f := &FakeTask{}

	_ = f.Run("testing", func() error { return nil })

	if !f.CalledRun {
		t.Errorf("failed to use mocked Run function on FakeCommand")
	}
}
