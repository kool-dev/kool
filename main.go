package main

import (
	"log"
	"os"

	"kool-dev/kool/cmd"
	"kool-dev/kool/enviroment"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	enviroment.InitEnvironmentVariables()

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
