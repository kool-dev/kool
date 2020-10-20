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
	var result interface{}

	r := NewRunner()

	result, _ = r.Run("testing", func() (interface{}, error) {
		return true, nil
	})

	if ran, ok := result.(bool); !ok || !ran {
		t.Error("failed to run task on Runner")
	}

	_, err := r.Run("testing", func() (interface{}, error) {
		return nil, errors.New("run error")
	})

	if err != nil && err.Error() != "run error" {
		t.Error("failed to get task error")
	}
}
