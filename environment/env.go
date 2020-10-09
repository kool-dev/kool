package environment

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fireworkweb/godotenv"
	homedir "github.com/mitchellh/go-homedir"
)

var envFile string = ".env"

// InitEnvironmentVariables handles the reading of .env files and
// setting up important environment variables necessary for kool
// to operate as expected.
func InitEnvironmentVariables(envStorage EnvStorage, defaultEnvValues string) {
	var (
		homeDir, workDir string
		err              error
	)

	homeDir, err = homedir.Dir()
	if err != nil {
		log.Fatal("Could not evaluate HOME directory - ", err)
	}
	if envStorage.Get("HOME") == "" {
		envStorage.Set("HOME", homeDir)
	}
	if envStorage.Get("UID") == "" {
		envStorage.Set("UID", fmt.Sprintf("%d", os.Getuid()))
	}

	if envStorage.Get("PWD") == "" {
		workDir, err = os.Getwd()
		if err != nil {
			log.Fatal("Could not evaluate working directory - ", err)
		}
		envStorage.Set("PWD", workDir)
	}

	if _, err = os.Stat(envFile); !os.IsNotExist(err) {
		err = envStorage.Load(envFile)
		if err != nil {
			log.Fatal("Failure loading environment file ", envFile, " error: '", err, "'")
		}
	}

	// After loading all files, we should complemente the non-overwritten
	// default variables to their expected distribution values
	allEnv := os.Environ()
	currentEnv := map[string]bool{}
	for _, env := range allEnv {
		currentEnv[strings.Split(env, "=")[0]] = true
	}
	defaultEnv, _ := godotenv.Unmarshal(defaultEnvValues)
	for k, v := range defaultEnv {
		if _, exists := currentEnv[k]; !exists {
			envStorage.Set(k, v)
		}
	}

	// Now that we loaded up the files, we will check for
	// missing variables that we need to fix
	if envStorage.Get("KOOL_NAME") == "" {
		pieces := strings.Split(envStorage.Get("PWD"), string(os.PathSeparator))
		envStorage.Set("KOOL_NAME", pieces[len(pieces)-1])
	}

	if envStorage.Get("KOOL_GLOBAL_NETWORK") == "" {
		envStorage.Set("KOOL_GLOBAL_NETWORK", "kool_global")
	}

	initAsuser(envStorage)
}
