package commands

import (
	"kool-dev/kool/core/shell"
)

// KoolService interface holds the contract for a
// general service which implements some bigger chunk
// of logic usually linked to a command.
type KoolService interface {
	Execute([]string) error

	Shell() shell.Shell
}

// DefaultKoolService holds handlers and functions shared by all
// services, meant to be used on commands when executing the services.
type DefaultKoolService struct {
	shell shell.Shell
}

func newDefaultKoolService() *DefaultKoolService {
	return &DefaultKoolService{
		shell.NewShell(),
	}
}

// Shell exposes the attached shell implementation
func (k *DefaultKoolService) Shell() shell.Shell {
	return k.shell
}

// Fake changes the internal dependencies (most notably shell)
// to be the fake conterpart of the real implementation.
// Meant for tests only.
func (k *DefaultKoolService) Fake() *DefaultKoolService {
	k.shell = &shell.FakeShell{MockIsTerminal: true}
	return k
}
