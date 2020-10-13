package cmd

import (
	"io"
	"kool-dev/kool/cmd/shell"
)

// KoolService interface holds the contract for a
// general service which implements some bigger chunk
// of logic usually linked to a command.
type KoolService interface {
	Execute([]string) error

	shell.Exiter
	shell.OutputWriter
	shell.InputReader
}

// DefaultKoolService holds handlers and functions shared by all
// services, meant to be used on commands when executing the services.
type DefaultKoolService struct {
	exiter shell.Exiter
	out    shell.OutputWriter
	in     shell.InputReader
}

func newDefaultKoolService() *DefaultKoolService {
	return &DefaultKoolService{
		shell.NewExiter(),
		shell.NewOutputWriter(),
		shell.NewInputReader(),
	}
}

// Exit proxies the call the given Exiter
func (k *DefaultKoolService) Exit(code int) {
	k.exiter.Exit(code)
}

// GetWriter proxies the call to the given OutputWriter
func (k *DefaultKoolService) GetWriter() (w io.Writer) {
	w = k.out.GetWriter()
	return
}

// SetWriter proxies the call to the given OutputWriter
func (k *DefaultKoolService) SetWriter(w io.Writer) {
	k.out.SetWriter(w)
}

// GetReader proxies the call to the given InputReader
func (k *DefaultKoolService) GetReader() (r io.Reader) {
	r = k.in.GetReader()
	return
}

// SetReader proxies the call to the given InputReader
func (k *DefaultKoolService) SetReader(r io.Reader) {
	k.in.SetReader(r)
}

// Println proxies the call to the given OutputWriter
func (k *DefaultKoolService) Println(out ...interface{}) {
	k.out.Println(out...)
}

// Error proxies the call to the given OutputWriter
func (k *DefaultKoolService) Error(err error) {
	k.out.Error(err)
}

// Warning proxies the call to the given OutputWriter
func (k *DefaultKoolService) Warning(out ...interface{}) {
	k.out.Warning(out...)
}

// Success proxies the call to the given OutputWriter
func (k *DefaultKoolService) Success(out ...interface{}) {
	k.out.Success(out...)
}
