package shell

import (
	"errors"
	"fmt"
)

var ErrUserCancelled = fmt.Errorf("user cancelled the operation")

func IsUserCancelledError(err error) bool {
	return errors.Is(err, ErrUserCancelled)
}

type ErrExitable struct {
	Err  error
	Code int
}

func (e ErrExitable) Error() string {
	return e.Err.Error()
}
