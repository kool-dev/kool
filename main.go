package main

import (
	"log"
	"os"

	"kool-dev/kool/commands"
	"kool-dev/kool/core/environment"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	environment.InitEnvironmentVariables(environment.NewEnvStorage())

	if err := commands.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
