package cmd

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func shellExec(exe string, args ...string) (outStr string, err error) {
	var (
		cmd *exec.Cmd
		out []byte
	)

	cmd = exec.Command(exe, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	out, err = cmd.CombinedOutput()

	outStr = strings.TrimSpace(string(out))
	return
}

func shellInteractive(exe string, args ...string) (err error) {
	var (
		cmd *exec.Cmd
	)

	cmd = exec.Command(exe, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
				// Not clear how we can hit this, but probably not
				// worth terminating the child.
				// fmt.Println("error sending signal", sig, err)
			}
		}
	}
}

func execError(out string, err error) {
	log.Println("ERROR: ", err)
	log.Println("Output:")
	fmt.Println(out)
}
