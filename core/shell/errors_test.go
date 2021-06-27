package shell

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsUserCancelledError(t *testing.T) {
	err := fmt.Errorf("error")

	if IsUserCancelledError(err) {
		t.Error("method IsUserCancelledError should return false for non user cancelled errors")
	}

	err = ErrUserCancelled

	if !IsUserCancelledError(err) {
		t.Error("method IsUserCancelledError should return true for user cancelled errors")
	}
}

func TestExitable(t *testing.T) {
	err := errors.New("some error")

	exitable := ErrExitable{Err: err}

	if exitable.Error() != err.Error() {
		t.Error("error should be the same")
	}
}
