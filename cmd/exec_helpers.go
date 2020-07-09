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

func shellInteractive(exe string, args ...string) (err error) {
	var (
		cmd *exec.Cmd
	)

	if exe == "docker-compose" {
		args = append(dockerComposeDefaultArgs(), args...)
	}

	if verbose := os.Getenv("KOOL_VERBOSE"); verbose == "1" || verbose == "true" {
		fmt.Println("$", exe, strings.Join(args, " "))
	}

	cmd = exec.Command(exe, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if numArgs := len(args); numArgs >= 2 && args[numArgs-2] == "<" {
		var (
			file *os.File
			path string
		)

		// we have an input redirection attempt!
		path = args[numArgs-1]
		args = args[:numArgs-2]

		file, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		cmd.Stdin = file
		defer file.Close()
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
