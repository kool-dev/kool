package commands

import (
	"kool-dev/kool/core/shell"
)

// FakeKoolService is a mock to be used on testing/replacement for KoolService interface
type FakeKoolService struct {
	shell *shell.FakeShell

	ArgsExecute    []string
	CalledExecute  bool
	MockExecuteErr error
}

func newFakeKoolService() *FakeKoolService {
	return &FakeKoolService{
		&shell.FakeShell{MockIsTerminal: true}, nil, false, nil,
	}
}

// Shell returns the stored (possibly fake) Shell implementation
func (f *FakeKoolService) Shell() shell.Shell {
	if f.shell == nil {
		// workaround interface nil-nil requirement
		// for type/value
		return nil
	}

	return f.shell
}

// Execute mocks the function for testing
func (f *FakeKoolService) Execute(args []string) (err error) {
	f.ArgsExecute = args
	f.CalledExecute = true
	err = f.MockExecuteErr
	return
}
