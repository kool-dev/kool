package cmd

import (
	"os"
	"os/exec"
	"strings"
)

func shellExec(exe string, args ...string) (outStr string, err error) {
	var (
		cmd *exec.Cmd
		out []byte
	)

	cmd = exec.Command(exe, args...)
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	outStr = strings.TrimSpace(string(out))
	return
}
