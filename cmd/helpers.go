package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/shell"
	"os"
	"os/exec"
)

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

	if _, err = shell.Exec("docker", "info"); err != nil {
		execError("Docker daemon doesn't seem to be running, run it first and retry.", err)
		os.Exit(1)
	}
}

func execError(out string, err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	if out != "" {
		fmt.Println("Output:", out)
	}
}
