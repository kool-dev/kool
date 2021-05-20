// +build ignore

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const templatesTemplate string = `package presets

// auto generated file

`

func main() {
	var (
		folders []os.DirEntry
		files   []os.DirEntry
		err     error
	)

	fmt.Println("Started building cmd/presets/templates.go")

	templates, err := os.Create("cmd/presets/templates.go")

	if err != nil {
		log.Fatal(err)
	}

	defer templates.Close()

	folders, err = os.ReadDir("templates")

	if err != nil {
		log.Fatal(err)
	}

	templates.WriteString(templatesTemplate)
	templates.WriteString("// GetTemplates get all templates\n")
	templates.WriteString("func GetTemplates() map[string]map[string]string {\n")
	templates.WriteString("\tvar templates = make(map[string]map[string]string)\n")

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		fmt.Println("Found folder", folder.Name())

		templates.WriteString(fmt.Sprintf("\ttemplates[\"%s\"] = map[string]string{\n", folder.Name()))

		files, err = os.ReadDir(filepath.Join("templates", folder.Name()))

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			templFile, err := os.Open(filepath.Join("templates", folder.Name(), file.Name()))

			if err != nil {
				log.Fatal(err)
			}

			filebytes, err := io.ReadAll(templFile)

			if err != nil {
				log.Fatal(err)
			}

			filecontent := string(filebytes)

			templates.WriteString(fmt.Sprintf("\t\t\"%s\": `%s`,\n", file.Name(), filecontent))
			fmt.Println("Parsed file:", file.Name())

			templFile.Close()
		}

		templates.WriteString("\t}\n")
	}

	templates.WriteString("\treturn templates\n")
	templates.WriteString("}\n")

	fmt.Println("Finished building cmd/presets/templates.go")
}
