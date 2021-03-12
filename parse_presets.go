// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const presetsTemplate string = `package presets

// auto generated file

`

func main() {
	var (
		folders []os.FileInfo
		files   []os.FileInfo
		err     error
	)

	fmt.Println("Started building cmd/presets/presets.go")

	presets, err := os.Create("cmd/presets/presets.go")

	if err != nil {
		log.Fatal(err)
	}

	defer presets.Close()

	folders, err = os.ReadDir("presets")

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

		files, err = os.ReadDir(filepath.Join("presets", folder.Name()))

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() || file.Name() == "preset-config.yml" {
				continue
			}

			presetFile, err := os.Open(filepath.Join("presets", folder.Name(), file.Name()))

			if err != nil {
				log.Fatal(err)
			}

			filebytes, err := io.ReadAll(presetFile)

			if err != nil {
				log.Fatal(err)
			}

			filecontent := string(filebytes)

			presets.WriteString(fmt.Sprintf("\t\t\"%s\": `%s`,\n", file.Name(), filecontent))
			fmt.Println("Parsed file:", file.Name())

			presetFile.Close()
		}

		presets.WriteString("\t}\n")
	}

	presets.WriteString("\treturn presets\n")
	presets.WriteString("}\n")

	fmt.Println("Finished building cmd/presets/presets.go")
}
