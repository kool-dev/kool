package cmd

import (
	"kool-dev/kool/cmd/shell"
	"os"
	"os/exec"
)

func checkKoolDependencies() {
	var err error

	if _, err = exec.LookPath("docker"); err != nil {
		shell.ExecError("Docker doesn't seem to be installed, install it first and retry.", err)
		os.Exit(1)
	}

	if _, err = exec.LookPath("docker-compose"); err != nil {
		shell.ExecError("Docker-compose doesn't seem to be installed, install it first and retry.", err)
		os.Exit(1)
	}

	if _, err = shell.Exec("docker", "info"); err != nil {
		shell.ExecError("Docker daemon doesn't seem to be running, run it first and retry.", err)
		os.Exit(1)
	}
}
