package environment

import (
	"log"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

var envFiles = []string{".env.local", ".env"}

// InitEnvironmentVariables handles the reading of .env files and
// setting up important environment variables necessary for kool
// to operate as expected.
func InitEnvironmentVariables(envStorage EnvStorage) {
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

	initUid(envStorage)

	if envStorage.Get("PWD") == "" {
		workDir, err = os.Getwd()
		if err != nil {
			log.Fatal("Could not evaluate working directory - ", err)
		}
		envStorage.Set("PWD", workDir)
	}

	for _, envFile := range envFiles {
		if _, err = os.Stat(envFile); os.IsNotExist(err) {
			continue
		}

		err = envStorage.Load(envFile)
		if err != nil {
			log.Fatal("Failure loading environment file ", envFile, " error: '", err, "'")
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
