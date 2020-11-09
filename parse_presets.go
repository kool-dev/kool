// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const presetsTemplate string = `package presets

// auto generated file

`

func main() {
	var (
		folders []os.FileInfo
		files []os.FileInfo
		err error
	)
	presets, err := os.Create("cmd/presets/presets.go")

	if err != nil {
		log.Fatal(err)
	}

	defer presets.Close()

	folders, err = ioutil.ReadDir("presets")

	if err != nil {
		log.Fatal(err)
	}

	presets.WriteString(presetsTemplate)
	presets.WriteString("// GetAll get all presets\n")
	presets.WriteString("func GetAll() map[string]map[string]string {\n")
	presets.WriteString("\tvar presets = make(map[string]map[string]string)\n")

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		fmt.Println("Found folder", folder.Name())

		presets.WriteString(fmt.Sprintf("\tpresets[\"%s\"] = map[string]string{\n", folder.Name()))

		files, err = ioutil.ReadDir(filepath.Join("presets", folder.Name()))

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			presetFile, err := os.Open(filepath.Join("presets", folder.Name(), file.Name()))

			if err != nil {
				log.Fatal(err)
			}

			filebytes, err := ioutil.ReadAll(presetFile)

			if err != nil {
				log.Fatal(err)
			}

			filecontent := string(filebytes)

			if file.Name() == ".preset" {
				lines := strings.Split(strings.TrimSpace(filecontent), "\n")

				for _, line := range lines {
					metadata := strings.Split(line, "=")
					presets.WriteString(fmt.Sprintf("\t\t\"preset_%s\": \"%s\",\n", metadata[0], metadata[1]))
				}
			} else {
				presets.WriteString(fmt.Sprintf("\t\t\"%s\": `%s`,\n", file.Name(), filecontent))
				fmt.Println("Parsed file:", file.Name())
			}

			presetFile.Close()
		}

		presets.WriteString("\t}\n")
	}

	presets.WriteString("\treturn presets\n")
	presets.WriteString("}\n")

	presets.WriteString("// GetTemplates get all templates\n")
	presets.WriteString("func GetTemplates() map[string]map[string]string {\n")
	presets.WriteString("\tvar templates = make(map[string]map[string]string)\n")

	folders, err = ioutil.ReadDir("templates")

	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		fmt.Println("Found folder", folder.Name())

		presets.WriteString(fmt.Sprintf("\ttemplates[\"%s\"] = map[string]string{\n", folder.Name()))

		files, err = ioutil.ReadDir(filepath.Join("templates", folder.Name()))

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			templFile, err := os.Open(filepath.Join("templates", folder.Name(), file.Name()))

			if err != nil {
				log.Fatal(err)
			}

			filebytes, err := ioutil.ReadAll(templFile)

			if err != nil {
				log.Fatal(err)
			}

			filecontent := string(filebytes)

			presets.WriteString(fmt.Sprintf("\t\t\"%s\": `%s`,\n", file.Name(), filecontent))
			fmt.Println("Parsed file:", file.Name())

			templFile.Close()
		}

		presets.WriteString("\t}\n")
	}

	presets.WriteString("\treturn templates\n")
	presets.WriteString("}\n")

	fmt.Println("Finished building cmd/presets/presets.go")
}
