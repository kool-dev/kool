package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fireworkweb/godotenv"
	homedir "github.com/mitchellh/go-homedir"
)

func initEnvironmentVariables() {
	var (
		files            []string
		homeDir, workDir string
		err              error
	)

	files = []string{
		".env",
		".fwd",
	}

	homeDir, err = homedir.Dir()
	if err != nil {
		log.Fatal("Could not evaluate HOME directory - ", err)
	}
	if os.Getenv("HOME") == "" {
		os.Setenv("HOME", homeDir)
	}
	if os.Getenv("UID") == "" {
		os.Setenv("UID", fmt.Sprintf("%d", os.Getuid()))
	}

	if os.Getenv("PWD") == "" {
		workDir, err = os.Getwd()
		if err != nil {
			log.Fatal("Could not evaluate working directory - ", err)
		}
		os.Setenv("PWD", workDir)
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

	// After loading all files, we should complemente the non-overwritten
	// default variables to their expected fwd distribution values
	allEnv := os.Environ()
	currentEnv := map[string]bool{}
	for _, env := range allEnv {
		currentEnv[strings.Split(env, "=")[0]] = true
	}
	allEnv = nil
	defaultEnv, _ := godotenv.Unmarshal(DefaultEnv)
	for k, v := range defaultEnv {
		if _, exists := currentEnv[k]; !exists {
			os.Setenv(k, v)
		}
	}

	// Now that we loaded up the files, we will check for
	// missing variables that we need to fix
	if os.Getenv("FWD_NAME") == "" {
		pieces := strings.Split(os.Getenv("PWD"), string(os.PathSeparator))
		os.Setenv("FWD_NAME", pieces[len(pieces)-1])
	}

	if os.Getenv("FWD_NETWORK") == "" {
		os.Setenv("FWD_NETWORK", "fwd_global")
	}
	if os.Getenv("FWD_ASUSER") == "" {
		os.Setenv("FWD_ASUSER", fmt.Sprintf("%d", os.Getuid()))
	}
}
