package shell

import "os"

type DefaultExiter struct{}

type Exiter interface {
	Exit(int)
}

func NewExiter() Exiter {
	return &DefaultExiter{}
}

func (e *DefaultExiter) Exit(code int) {
	os.Exit(1)
}
