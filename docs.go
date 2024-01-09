//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"errors"
	"fmt"
	"kool-dev/kool/commands"
	"kool-dev/kool/core/shell"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	var (
		err        error
		koolOutput *bytes.Buffer
		koolFile   *os.File
	)

	fmt.Println("Going to generate cobra docs in markdown...")

	koolOutput = new(bytes.Buffer)

	err = doc.GenMarkdownCustom(commands.RootCmd(), koolOutput, linkHandler)

	if err != nil {
		log.Fatal(err)
	}

	koolMarkdown := koolOutput.String()

	for _, childCmd := range commands.RootCmd().Commands() {
		if err = exportCmdDocs(childCmd, &koolMarkdown); err != nil {
			if strings.HasPrefix(err.Error(), "skip") {
				log.Println(err)
			} else {
				log.Fatal(err)
			}
		}
	}

	re := regexp.MustCompile("(?m)[\r\n]+^.*kool_deploy.*$")
	koolMarkdown = re.ReplaceAllString(koolMarkdown, "")

	koolFile, err = CreateFile("0-kool", "docs/05-Commands-Reference")

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

	shell.Success("Success!")
}

// CreateFile Create file to write markdown content
func CreateFile(filename string, dir string) (file *os.File, err error) {
	basename := fmt.Sprintf("%s.md", filename)

	file, err = os.Create(filepath.Join(dir, basename))

	return
}

func exportCmdDocs(childCmd *cobra.Command, koolMarkdown *string) (err error) {
	var (
		cmdName string
		cmdFile *os.File
	)

	if cmdName = strings.Replace(childCmd.CommandPath(), " ", "_", -1); cmdName == "kool_help" {
		err = errors.New("skip kool_help")
		return
	}

	newName := strings.Replace(childCmd.CommandPath(), " ", "-", -1)
	*koolMarkdown = strings.Replace(*koolMarkdown, cmdName, newName, -1)

	cmdOutput := new(bytes.Buffer)

	if err = doc.GenMarkdownCustom(childCmd, cmdOutput, linkHandler); err != nil {
		return
	}

	if cmdFile, err = CreateFile(newName, "docs/05-Commands-Reference"); err != nil {
		return
	}

	defer cmdFile.Close()

	if _, err = cmdOutput.WriteTo(cmdFile); err != nil {
		return
	}

	for _, subCmd := range childCmd.Commands() {
		if err = exportCmdDocs(subCmd, koolMarkdown); err != nil {
			if strings.HasPrefix(err.Error(), "skip") {
				log.Println(err)
			} else {
				return
			}
		}
	}

	return
}

func linkHandler(filename string) string {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	return strings.ToLower(base)
}
