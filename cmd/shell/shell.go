package shell

import (
	"fmt"
	"io"
	"kool-dev/kool/environment"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var (
	lookedUp map[string]bool
)

// Exec will execute the given command silently and return the combined
// error/standard output, and an error if any.
func Exec(exe string, args ...string) (outStr string, err error) {
	var (
		cmd *exec.Cmd
		out []byte
	)

	if exe == "docker-compose" {
		args = append(dockerComposeDefaultArgs(), args...)
	}

	cmd = exec.Command(exe, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	out, err = cmd.CombinedOutput()
	outStr = strings.TrimSpace(string(out))
	return
}

// Interactive runs the given command proxying current Stdin/Stdout/Stderr
// which makes it interactive for running even something like `bash`.
func Interactive(exe string, args ...string) (err error) {
	var (
		cmd         *exec.Cmd
		cmdStdin    io.ReadCloser
		cmdStdout   io.WriteCloser
		closeStdin  bool
		closeStdout bool
	)

	if lookedUp == nil {
		lookedUp = make(map[string]bool)
	}

	if exe == "docker-compose" {
		args = append(dockerComposeDefaultArgs(), args...)
	}

	if environment.IsTrue("KOOL_VERBOSE") {
		fmt.Println("$", exe, strings.Join(args, " "))
	}

	// soon should refactor this onto a struct with methods
	// so we can remove this too long list of returned values.
	if args, cmdStdin, cmdStdout, closeStdin, closeStdout, err = parseRedirects(args); err != nil {
		return
	}

	if closeStdin {
		defer cmdStdin.Close()
	}
	if closeStdout {
		defer cmdStdout.Close()
	}

	cmd = exec.Command(exe, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = cmdStdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = cmdStdin

	if exe != "kool" && !lookedUp[exe] && !strings.HasPrefix(exe, "./") && !strings.HasPrefix(exe, "/") {
		// non-kool and non-absolute/relative path... let's look it up
		_, err = exec.LookPath(exe)

		if err != nil {
			Error("Failed to run ", cmd.String(), "error:", err)
			os.Exit(2)
		}

		lookedUp[exe] = true
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
					Error("error sending signal to child process", sig, err)
				}
			}
		}
	}
}

func dockerComposeDefaultArgs() []string {
	return []string{"-p", os.Getenv("KOOL_NAME")}
}
