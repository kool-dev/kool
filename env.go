package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	homedir "github.com/mitchellh/go-homedir"
)

func initEnvironmentVariables() {
	var (
		files   []string
		homeDir string
		err     error
	)

	files = []string{
		".env",
		".fwd",
	}

	homeDir, err = homedir.Dir()
	if err != nil {
		log.Fatal("Could not evaluate HOME directory - ", err)
	}

	files = append(files, fmt.Sprintf("%s%s.fwd", homeDir, string(filepath.Separator)))

	for _, file := range files {
		if _, err = os.Stat(file); os.IsNotExist(err) {
			continue
		}

		err = godotenv.Load(file)
		if err != nil {
			log.Fatal("Failure loading environment file ", file, " ERROR: '", err, "'")
		}
	}
}
