package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fireworkweb/godotenv"
	homedir "github.com/mitchellh/go-homedir"
)

func initEnvironmentVariables() {
	var (
		homeDir, workDir string
		err              error
	)

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

	var file string = ".env"
	if _, err = os.Stat(file); !os.IsNotExist(err) {
		err = godotenv.Load(file)
		if err != nil {
			log.Fatal("Failure loading environment file ", file, " ERROR: '", err, "'")
		}
	}

	// After loading all files, we should complemente the non-overwritten
	// default variables to their expected distribution values
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
	if os.Getenv("KOOL_NAME") == "" {
		pieces := strings.Split(os.Getenv("PWD"), string(os.PathSeparator))
		os.Setenv("KOOL_NAME", pieces[len(pieces)-1])
	}

	if os.Getenv("KOOL_GLOBAL_NETWORK") == "" {
		os.Setenv("KOOL_GLOBAL_NETWORK", "kool_global")
	}
	if os.Getenv("KOOL_ASUSER") == "" {
		os.Setenv("KOOL_ASUSER", fmt.Sprintf("%d", os.Getuid()))
	}
}
