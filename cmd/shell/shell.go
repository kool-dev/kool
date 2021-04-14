package shell

import (
	"errors"
	"fmt"
	"io"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/gookit/color"
)

type execCmdFnType func(string, ...string) *exec.Cmd
type execLookPathFnType func(string) (string, error)

var execCmdFn execCmdFnType = exec.Command
var execLookPathFn execLookPathFnType = exec.LookPath

// ErrLookPath error when we fail to find in PATH
var ErrLookPath = errors.New("command not found")

// RecursiveCall is used to proxy self-executaion commands internally
// instead of creating a whole new OS process
var RecursiveCall func([]string, io.Reader, io.Writer, io.Writer) error

// DefaultShell holds data for handling a shell
type DefaultShell struct {
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
	lookedUp  *lookupCache
	env       environment.EnvStorage
}

// OutputWritter implements basic output for CLIss
type OutputWritter interface {
	Println(...interface{})
	Printf(string, ...interface{})
	Warning(...interface{})
	Success(...interface{})
}

type PathChecker interface {
	LookPath(builder.Command) error
}

// Shell implements functions for handling a shell
type Shell interface {
	OutputWritter
	PathChecker
	InStream() io.Reader
	SetInStream(io.Reader)
	OutStream() io.Writer
	SetOutStream(io.Writer)
	ErrStream() io.Writer
	SetErrStream(io.Writer)
	Exec(builder.Command, ...string) (string, error)
	Interactive(builder.Command, ...string) error
	Error(error)
}

// NewShell creates a new shell
func NewShell() Shell {
	return &DefaultShell{
		inStream:  os.Stdin,
		outStream: os.Stdout,
		errStream: os.Stderr,
		env:       environment.NewEnvStorage(),
		lookedUp:  newLookupCache(),
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
		cmd     *exec.Cmd
		out     []byte
		args    []string = command.Args()
		exe     string   = command.Cmd()
		verbose bool     = s.env.IsTrue("KOOL_VERBOSE")
	)

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	if verbose {
		fmt.Fprintf(s.ErrStream(), "$ (exec) %s %v\n",
			exe,
			args,
		)
	}

	cmd = execCmdFn(exe, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = s.InStream()
	out, err = cmd.CombinedOutput()
	outStr = strings.TrimSpace(string(out))
	if err != nil && len(out) != 0 {
		// let's use the actual output for error, appending practical exec error
		// (most probably the later will be an non-zero exit status error)
		err = fmt.Errorf("%s (%s)", outStr, err.Error())
	}
	return
}

// Interactive runs the given command proxying current Stdin/Stdout/Stderr
// which makes it interactive for running even something like `bash`.
func (s *DefaultShell) Interactive(originalCmd builder.Command, extraArgs ...string) (err error) {
	var (
		cmdptr  *CommandWithPointers
		verbose bool = s.env.IsTrue("KOOL_VERBOSE")
		command      = originalCmd.Copy()
	)

	command.AppendArgs(extraArgs...)

	// soon should refactor this onto a struct with methods
	// so we can remove this too long list of returned values.
	if cmdptr, err = parseRedirects(command, s); err != nil {
		return
	}

	if verbose {
		checker := NewTerminalChecker()
		fmt.Fprintf(s.ErrStream(), "$ (TTY in: %v out: %v) %s %v\n",
			checker.IsTerminal(cmdptr.in),
			checker.IsTerminal(cmdptr.out),
			cmdptr.Command.Cmd(),
			cmdptr.Command.Args(),
		)
	}

	if cmdptr.Command.Cmd() == "kool" && RecursiveCall != nil {
		if verbose {
			fmt.Fprintln(s.ErrStream(), "[recursive call]")
		}
		err = RecursiveCall(cmdptr.Command.Args(), cmdptr.in, cmdptr.out, cmdptr.err)
	} else {
		if err = s.LookPath(cmdptr.Command); err != nil {
			err = ErrLookPath
			return
		}

		err = s.execute(cmdptr.Cmd())

		defer cmdptr.Close()
	}
	return
}

// LookPath returns if the command exists
func (s *DefaultShell) LookPath(command builder.Command) (err error) {
	var (
		exe       string = command.Cmd()
		hasLooked bool
	)

	if strings.HasPrefix(exe, "./") || strings.HasPrefix(exe, "/") || strings.HasPrefix(exe, "../") {
		// either absolute/relative path... don't need to check
		return
	}

	if hasLooked, err = s.lookedUp.fetch(exe); err != nil {
		return
	}

	if !hasLooked {
		_, err = execLookPathFn(exe)

		s.lookedUp.set(exe, err)
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

func (s *DefaultShell) execute(cmd *exec.Cmd) (err error) {
	if err = cmd.Start(); err != nil {
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
