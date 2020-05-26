package cmd

import (
	"fmt"
	"log"
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

func execError(out string, err error) {
	log.Println("ERROR: ", err)
	log.Println("Output:")
	fmt.Println(out)
}
