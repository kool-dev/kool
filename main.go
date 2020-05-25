package main

import (
	"log"
	"os"

	fwd "./cmd"
)

func main() {
	// log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	initEnvironmentVariables()

	if err := fwd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
