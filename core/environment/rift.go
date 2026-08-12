package environment

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var riftAncestors = func(workDir string) ([]byte, error) {
	return exec.Command("rift", "ancestors", workDir).Output()
}

func initRift(envStorage EnvStorage, workDir string) {
	workspace, found := findRiftWorkspace(workDir)
	if !found {
		return
	}

	output, err := riftAncestors(workDir)
	if err != nil {
		envStorage.Set("KOOL_WORKSPACE_ERROR", fmt.Sprintf("could not resolve Rift workspace ancestry: %v", err))
		return
	}

	var ancestors []string
	for _, ancestor := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if ancestor != "" {
			ancestors = append(ancestors, ancestor)
		}
	}
	source := workspace
	if len(ancestors) > 0 {
		source = ancestors[len(ancestors)-1]
	}

	initWorkspaceContext(envStorage, workDir, source, workspace, "rift", len(ancestors) > 0, false)
}

func findRiftWorkspace(workDir string) (string, bool) {
	for directory := filepath.Clean(workDir); ; directory = filepath.Dir(directory) {
		if _, err := os.Stat(filepath.Join(directory, ".rift")); err == nil {
			return directory, true
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			return "", false
		}
	}
}

func composeProjectName(name string) string {
	name = strings.ToLower(name)
	name = strings.Map(func(character rune) rune {
		if character >= 'a' && character <= 'z' || character >= '0' && character <= '9' || character == '-' || character == '_' {
			return character
		}
		return -1
	}, name)

	name = strings.TrimLeft(name, "-_")
	if name == "" {
		return "kool"
	}
	return name
}
