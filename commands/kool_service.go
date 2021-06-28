package commands

import (
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/shell"
	"os/exec"
	"syscall"
)

// KoolService interface holds the contract for a
// general service which implements some bigger chunk
// of logic usually linked to a command.
type KoolService interface {
	Execute([]string) error
	IsTerminal() bool

	shell.Shell
}

// DefaultKoolService holds handlers and functions shared by all
// services, meant to be used on commands when executing the services.
type DefaultKoolService struct {
	term  shell.TerminalChecker
	shell shell.Shell
}

func newDefaultKoolService() *DefaultKoolService {
	return &DefaultKoolService{
		shell.NewTerminalChecker(),
		shell.NewShell(),
	}
}

// Println proxies the call to the given Shell
func (k *DefaultKoolService) Println(out ...interface{}) {
	k.shell.Println(out...)
}

// Printf proxies the call to the given Shell
func (k *DefaultKoolService) Printf(format string, a ...interface{}) {
	k.shell.Printf(format, a...)
}

// Error proxies the call to the given Shell
func (k *DefaultKoolService) Error(err error) {
	k.shell.Error(err)
}

// Warning proxies the call to the given Shell
func (k *DefaultKoolService) Warning(out ...interface{}) {
	k.shell.Warning(out...)
}

// Success proxies the call to the given Shell
func (k *DefaultKoolService) Success(out ...interface{}) {
	k.shell.Success(out...)
}

// IsTerminal checks if input/output is a terminal
func (k *DefaultKoolService) IsTerminal() bool {
	return k.term.IsTerminal(k.InStream(), k.OutStream())
}

// InStream proxies the call to the given Shell
func (k *DefaultKoolService) InStream() io.Reader {
	return k.shell.InStream()
}

// SetInStream proxies the call to the given Shell
func (k *DefaultKoolService) SetInStream(inStream io.Reader) {
	k.shell.SetInStream(inStream)
}

// OutStream proxies the call to the given Shell
func (k *DefaultKoolService) OutStream() io.Writer {
	return k.shell.OutStream()
}

// SetOutStream proxies the call to the given Shell
func (k *DefaultKoolService) SetOutStream(outStream io.Writer) {
	k.shell.SetOutStream(outStream)
}

// ErrStream proxies the call to the given Shell
func (k *DefaultKoolService) ErrStream() io.Writer {
	return k.shell.ErrStream()
}

// SetErrStream proxies the call to the given Shell
func (k *DefaultKoolService) SetErrStream(errStream io.Writer) {
	k.shell.SetErrStream(errStream)
}

// Exec proxies the call to the given Shell
func (k *DefaultKoolService) Exec(command builder.Command, extraArgs ...string) (outStr string, err error) {
	outStr, err = k.shell.Exec(command, extraArgs...)
	return
}

// Interactive proxies the call to the given Shell
func (k *DefaultKoolService) Interactive(command builder.Command, extraArgs ...string) (err error) {
	err = k.shell.Interactive(command, extraArgs...)

	if err == shell.ErrLookPath {
		err = fmt.Errorf("failed to run %s error: %v", command.String(), err)
		return
	}

	// Subprocess exited. Get the return code, if we can
	if exitError, ok := err.(*exec.ExitError); ok {
		err = shell.ErrExitable{
			Err:  err,
			Code: exitError.Sys().(syscall.WaitStatus).ExitStatus(),
		}
	}

	return
}

// LookPath proxies the call to the given Shell
func (k *DefaultKoolService) LookPath(command builder.Command) (err error) {
	err = k.shell.LookPath(command)
	return
}
