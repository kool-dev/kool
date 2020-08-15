package main

import (
	"fmt"
	"kool-dev/kool/cmd"
	"log"

	"github.com/spf13/cobra/doc"
)

func main() {
	fmt.Println("Going to generate cobra docs in markdown...")
	err := doc.GenMarkdownTree(cmd.RootCmd(), "./4-Commands/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")
}
