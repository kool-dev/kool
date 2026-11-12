package main

import (
	"errors"
	"log"
	"os"

	"kool-dev/kool/commands"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if err := commands.Execute(); err != nil {
		s := shell.NewShell()
		if s.IsJSONOutput() {
			s.Error(err)
		} else {
			s.Println(err)
		}
		code := 1
		var ex shell.ErrExitable
		if errors.As(err, &ex) {
			code = ex.Code
		}
		environment.CleanupWorkspace()
		os.Exit(code)
	}

	environment.CleanupWorkspace()
	os.Exit(0)
}
