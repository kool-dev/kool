package task

import (
	"errors"
	"testing"
)

func TestNewTaskRunner(t *testing.T) {
	r := NewRunner()

	if _, ok := r.(*DefaultRunner); !ok {
		t.Error("NewRunner() did not return a *DefaultRunner")
	}
}

func TestRunTaskRunner(t *testing.T) {
	var ran bool = false
	r := NewRunner()

	_ = r.Run("testing", func() error {
		ran = true
		return nil
	})

	if !ran {
		t.Error("failed to run task on Runner")
	}

	ran = false

	err := r.Run("testing", func() error {
		return errors.New("run error")
	})

	if err != nil && err.Error() != "run error" {
		t.Error("failed to get task error")
	}
}
