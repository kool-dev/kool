package environment

import (
	"kool-dev/kool/core/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var envFiles = []string{".env.local", ".env"}
var loadedEnvKeys = make(map[string][]string)

// InitEnvironmentVariables handles the reading of .env files and
// setting up important environment variables necessary for kool
// to operate as expected.
func InitEnvironmentVariables(envStorage EnvStorage) {
	var (
		homeDir, workDir string
		err              error
	)

	homeDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not evaluate HOME directory - ", err)
	}
	if envStorage.Get("HOME") == "" {
		envStorage.Set("HOME", homeDir)
	}

	initUid(envStorage)

	// Always use os.Getwd() to determine the working directory rather than
	// trusting the PWD environment variable. PWD may be stale if a parent
	// process spawned kool with a different cwd without updating PWD.
	workDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Could not evaluate working directory - ", err)
	}
	envStorage.Set("PWD", workDir)
	envDirectory := canonicalDirectory(workDir)
	loadedEnvKeys[envDirectory] = nil

	for _, envFile := range envFiles {
		if _, err = os.Stat(envFile); os.IsNotExist(err) {
			continue
		}

		before := envKeySet(envStorage.All())
		err = envStorage.Load(envFile)
		if err != nil {
			log.Fatal("Failure loading environment file ", envFile, " error: '", err, "'")
		}
		for key := range envKeySet(envStorage.All()) {
			if !before[key] {
				loadedEnvKeys[envDirectory] = append(loadedEnvKeys[envDirectory], key)
			}
		}
	}

	config := loadKoolConfig(workDir)
	if config != nil && len(config.Workspaces) > 0 {
		envStorage.Set("KOOL_WORKSPACES_ENABLED", "true")
		initRift(envStorage, workDir)
		initGitWorktree(envStorage, workDir)
		initSourceProject(envStorage, workDir)
	}
	if config != nil && config.Proxy != nil {
		envStorage.Set("KOOL_PROXY_ENABLED", "true")
		initSourceProject(envStorage, workDir)
		initProxy(envStorage, workDir)
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

// LoadedEnvKeys returns variables introduced from a directory's environment files.
func LoadedEnvKeys(directory string) []string {
	return append([]string(nil), loadedEnvKeys[canonicalDirectory(directory)]...)
}

func canonicalDirectory(directory string) string {
	canonical, err := filepath.EvalSymlinks(directory)
	if err == nil {
		return canonical
	}
	canonical, err = filepath.Abs(directory)
	if err == nil {
		return canonical
	}
	return directory
}

func envKeySet(entries []string) map[string]bool {
	keys := make(map[string]bool, len(entries))
	for _, entry := range entries {
		keys[strings.SplitN(entry, "=", 2)[0]] = true
	}
	return keys
}

func loadKoolConfig(workDir string) *parser.KoolYaml {
	for _, name := range []string{"kool.yml", "kool.yaml"} {
		config, err := parser.ParseKoolYaml(filepath.Join(workDir, name))
		if err == nil {
			return config
		}
	}
	return nil
}

func initSourceProject(envStorage EnvStorage, workDir string) {
	if envStorage.Get("KOOL_WORKSPACE_SOURCE_PROJECT") != "" {
		return
	}
	project := envStorage.Get("COMPOSE_PROJECT_NAME")
	if project == "" {
		project = composeSourceProject(envStorage, workDir)
	}
	if project == "" {
		project = composeProjectName(filepath.Base(workDir))
	}
	envStorage.Set("KOOL_WORKSPACE_SOURCE_PROJECT", project)
}
