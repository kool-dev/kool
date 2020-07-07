package main

import (
	"log"
	"os"

	"kool-dev/kool/cmd"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	initEnvironmentVariables()

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
