// +build ignore

package main

import (
	"bytes"
	"fmt"
	"kool-dev/kool/cmd"
	"kool-dev/kool/cmd/shell"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra/doc"
)

func main() {
	var (
		err        error
		koolOutput *bytes.Buffer
		cmdFile    *os.File
		koolFile   *os.File
	)

	fmt.Println("Going to generate cobra docs in markdown...")

	koolOutput = new(bytes.Buffer)

	err = doc.GenMarkdown(cmd.RootCmd(), koolOutput)

	if err != nil {
		log.Fatal(err)
	}

	koolMarkdown := koolOutput.String()

	for _, childCmd := range cmd.RootCmd().Commands() {
		var cmdName string

		if cmdName = strings.Replace(childCmd.CommandPath(), " ", "_", -1); cmdName == "kool_deploy" || cmdName == "kool_help" {
			continue
		}

		newName := strings.Replace(childCmd.CommandPath(), " ", "-", -1)
		koolMarkdown = strings.Replace(koolMarkdown, cmdName, newName, -1)

		cmdOutput := new(bytes.Buffer)

		err = doc.GenMarkdown(childCmd, cmdOutput)

		if err != nil {
			log.Fatal(err)
		}

		cmdFile, err = CreateFile(newName, "docs/4-Commands")

		if err != nil {
			log.Fatal(err)
		}

		defer cmdFile.Close()

		_, err = cmdOutput.WriteTo(cmdFile)

		if err != nil {
			log.Fatal(err)
		}
	}

	re := regexp.MustCompile("(?m)[\r\n]+^.*kool_deploy.*$")
	koolMarkdown = re.ReplaceAllString(koolMarkdown, "")

	koolFile, err = CreateFile("0-kool", "docs/4-Commands")

	if err != nil {
		log.Fatal(err)
	}

	defer koolFile.Close()

	koolOutput = new(bytes.Buffer)
	koolOutput.WriteString(koolMarkdown)

	_, err = koolOutput.WriteTo(koolFile)

	if err != nil {
		log.Fatal(err)
	}

	shell.NewOutputWriter().Success("Success!")
}

// CreateFile Create file to write markdown content
func CreateFile(filename string, dir string) (file *os.File, err error) {
	basename := fmt.Sprintf("%s.md", filename)

	file, err = os.Create(filepath.Join(dir, basename))

	return
}
