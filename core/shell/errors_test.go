package shell

import (
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
