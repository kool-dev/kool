package shell

import "os"

// DefaultExiter holds the default Exit behaviour
type DefaultExiter struct{}

// Exiter interface allows for interchageable usage of implementations
// mainly for testing and extension purposes.
type Exiter interface {
	Exit(int)
}

// NewExiter creates a new DefaultExiter
func NewExiter() Exiter {
	return &DefaultExiter{}
}

// Exit implements the default Exit behaviour (proxy to OS)
func (e *DefaultExiter) Exit(code int) {
	os.Exit(code)
}
