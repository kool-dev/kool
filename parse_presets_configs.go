// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const configsTemplate string = `package presets

// auto generated file

`

func main() {
	var (
		folders    []os.FileInfo
		configFile *os.File
		err        error
	)

	fmt.Println("Started building cmd/presets/configs.go")

	configs, err := os.Create("cmd/presets/configs.go")

	if err != nil {
		log.Fatal(err)
	}

	defer configs.Close()

	folders, err = os.ReadDir("presets")

	if err != nil {
		log.Fatal(err)
	}

	configs.WriteString(configsTemplate)
	configs.WriteString("// GetConfigs get all presets configs\n")
	configs.WriteString("func GetConfigs() map[string]string {\n")
	configs.WriteString("\tvar configs = make(map[string]string)\n")

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		fmt.Println("Found folder", folder.Name())
		configFilePath := filepath.Join("presets", folder.Name(), "preset-config.yml")

		if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
			continue
		}

		configFile, err = os.Open(configFilePath)

		if err != nil {
			log.Fatal(err)
		}

		filebytes, err := io.ReadAll(configFile)

		if err != nil {
			log.Fatal(err)
		}

		filecontent := string(filebytes)

		configs.WriteString(fmt.Sprintf("\tconfigs[\"%s\"] = `%s`\n", folder.Name(), filecontent))
		fmt.Println("Parsed file:", "preset-config.yml")

		configFile.Close()
	}

	configs.WriteString("\treturn configs\n")
	configs.WriteString("}\n")

	fmt.Println("Finished building cmd/presets/configs.go")
}
