package shell

import (
	"fmt"
	"io"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gookit/color"
)

// DefaultShell holds data for handling a shell
type DefaultShell struct {
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
	lookedUp  map[string]bool
}

// Shell implements functions for handling a shell
type Shell interface {
	InStream() io.Reader
	SetInStream(io.Reader)
	OutStream() io.Writer
	SetOutStream(io.Writer)
	ErrStream() io.Writer
	SetErrStream(io.Writer)
	Exec(builder.Command, ...string) (string, error)
	Interactive(builder.Command, ...string) error
	LookPath(builder.Command) error
	Println(...interface{})
	Printf(string, ...interface{})
	Error(error)
	Warning(...interface{})
	Success(...interface{})
}

// NewShell creates a new shell
func NewShell() Shell {
	return &DefaultShell{
		inStream:  os.Stdin,
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
}

// InStream get input stream
func (s *DefaultShell) InStream() io.Reader {
	return s.inStream
}

// SetInStream set input stream
func (s *DefaultShell) SetInStream(inStream io.Reader) {
	s.inStream = inStream
}

// OutStream get output stream
func (s *DefaultShell) OutStream() io.Writer {
	return s.outStream
}

// SetOutStream set output stream
func (s *DefaultShell) SetOutStream(outStream io.Writer) {
	s.outStream = outStream
}

// ErrStream get error stream
func (s *DefaultShell) ErrStream() io.Writer {
	return s.errStream
}

// SetErrStream set error stream
func (s *DefaultShell) SetErrStream(errStream io.Writer) {
	s.errStream = errStream
}

// Exec will execute the given command silently and return the combined
// error/standard output, and an error if any.
func (s *DefaultShell) Exec(command builder.Command, extraArgs ...string) (outStr string, err error) {
	var (
		cmd  *exec.Cmd
		out  []byte
		args []string = command.Args()
		exe  string   = command.Cmd()
	)

	if exe == "docker-compose" {
		args = append(dockerComposeDefaultArgs(), args...)
	}

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	cmd = exec.Command(exe, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = s.InStream()
	out, err = cmd.CombinedOutput()
	outStr = strings.TrimSpace(string(out))
	return
}

// Interactive runs the given command proxying current Stdin/Stdout/Stderr
// which makes it interactive for running even something like `bash`.
func (s *DefaultShell) Interactive(command builder.Command, extraArgs ...string) (err error) {
	var (
		cmd            *exec.Cmd
		parsedRedirect *DefaultParsedRedirect
		exe            string   = command.Cmd()
		args           []string = command.Args()
	)

	if exe == "docker-compose" {
		args = append(dockerComposeDefaultArgs(), args...)
	}

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	if environment.NewEnvStorage().IsTrue("KOOL_VERBOSE") {
		fmt.Println("$", exe, strings.Join(args, " "))
	}

	// soon should refactor this onto a struct with methods
	// so we can remove this too long list of returned values.
	if parsedRedirect, err = parseRedirects(args, s); err != nil {
		return
	}

	defer parsedRedirect.Close()

	cmd = parsedRedirect.CreateCommand(exe)

	if err = s.LookPath(command); err != nil {
		s.Error(fmt.Errorf("failed to run %s error: %v", cmd.String(), err))
		os.Exit(2)
	}

	err = cmd.Start()

	if err != nil {
		return
	}

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
		close(waitCh)
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	// You need a for loop to handle multiple signals
	for {
		select {
		case err = <-waitCh:
			// Subprocess exited. Get the return code, if we can
			var waitStatus syscall.WaitStatus
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus = exitError.Sys().(syscall.WaitStatus)
				os.Exit(waitStatus.ExitStatus())
			}
			if err != nil {
				log.Fatal(err)
			}
			return
		case sig := <-sigChan:
			if err := cmd.Process.Signal(sig); err != nil {
				// check if it is something we should care about
				if err.Error() != "os: process already finished" {
					s.Error(fmt.Errorf("error sending signal to child process %v %v", sig, err))
				}
			}
		}
	}
}

// LookPath returns if the command exists
func (s *DefaultShell) LookPath(command builder.Command) (err error) {
	var exe string = command.Cmd()

	if s.lookedUp == nil {
		s.lookedUp = make(map[string]bool)
	}

	if exe != "kool" && !s.lookedUp[exe] && !strings.HasPrefix(exe, "./") && !strings.HasPrefix(exe, "/") {
		// non-kool and non-absolute/relative path... let's look it up
		_, err = exec.LookPath(exe)

		s.lookedUp[exe] = true
	}
	return
}

// Println execs Println on writer
func (s *DefaultShell) Println(out ...interface{}) {
	fmt.Fprintln(s.OutStream(), out...)
}

// Printf execs Printf on writer
func (s *DefaultShell) Printf(format string, a ...interface{}) {
	fmt.Fprintf(s.OutStream(), format, a...)
}

// Error error output
func (s *DefaultShell) Error(err error) {
	fmt.Fprintf(s.OutStream(), "%v\n", color.New(color.BgRed, color.FgWhite).Sprintf("error: %v", err))
}

// Warning warning message
func (s *DefaultShell) Warning(out ...interface{}) {
	warningMessage := color.New(color.Yellow).Sprint(out...)
	fmt.Fprintln(s.OutStream(), warningMessage)
}

// Success success message
func (s *DefaultShell) Success(out ...interface{}) {
	successMessage := color.New(color.Green).Sprint(out...)
	fmt.Fprintln(s.OutStream(), successMessage)
}

// Exec will execute the given command silently and return the combined
// error/standard output, and an error if any.
func Exec(exe string, args ...string) (outStr string, err error) {
	command := builder.NewCommand(exe, args...)
	s := NewShell()
	outStr, err = s.Exec(command)
	return
}

// Interactive runs the given command proxying current Stdin/Stdout/Stderr
// which makes it interactive for running even something like `bash`.
func Interactive(exe string, args ...string) (err error) {
	command := builder.NewCommand(exe, args...)
	s := NewShell()
	err = s.Interactive(command)
	return
}

// Println execs Println on writer
func Println(out ...interface{}) {
	NewShell().Println(out)
}

// Printf execs Printf on writer
func Printf(format string, a ...interface{}) {
	NewShell().Printf(format, a)
}

// Error error output
func Error(err error) {
	NewShell().Error(err)
}

// Warning warning message
func Warning(out ...interface{}) {
	NewShell().Warning(out)
}

// Success success message
func Success(out ...interface{}) {
	NewShell().Success(out)
}

func dockerComposeDefaultArgs() []string {
	return []string{"-p", environment.NewEnvStorage().Get("KOOL_NAME")}
}
