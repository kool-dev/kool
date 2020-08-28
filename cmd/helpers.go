package cmd

import (
	"os"
	"os/exec"
)

func dockerComposeDefaultArgs() []string {
	return []string{"-p", os.Getenv("KOOL_NAME")}
}

func checkKoolDependencies() {
	var err error

	if _, err = exec.LookPath("docker"); err != nil {
		execError("Docker doesn't seem to be installed, install it first and retry.", err)
		os.Exit(1)
	}

	if _, err = exec.LookPath("docker-compose"); err != nil {
		execError("Docker-compose doesn't seem to be installed, install it first and retry.", err)
		os.Exit(1)
	}

	if _, err = shellExec("docker", "info"); err != nil {
		execError("Docker daemon doesn't seem to be running, run it first and retry.", err)
		os.Exit(1)
	}
}
