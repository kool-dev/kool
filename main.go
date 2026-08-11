package main

import (
	"log"
	"os"

	"kool-dev/kool/commands"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if err := commands.Execute(); err != nil {
		shell.NewShell().Println(err)
		code := 1
		if ex, ok := err.(shell.ErrExitable); ok {
			code = ex.Code
		}
		environment.CleanupWorkspace()
		os.Exit(code)
	}

	environment.CleanupWorkspace()
	os.Exit(0)
}
