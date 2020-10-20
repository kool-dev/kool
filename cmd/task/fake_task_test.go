package task

import "testing"

func TestFakeTask(t *testing.T) {
	f := &FakeTask{}

	_, _ = f.Run("testing", func() (interface{}, error) { return nil, nil })

	if !f.CalledRun {
		t.Errorf("failed to use mocked Run function on FakeCommand")
	}
}
