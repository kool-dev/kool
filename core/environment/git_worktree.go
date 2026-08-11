package environment

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

var gitWorktreeOutput = func(workDir string, args ...string) ([]byte, error) {
	arguments := append([]string{"-C", workDir}, args...)
	return exec.Command("git", arguments...).Output()
}

func initGitWorktree(envStorage EnvStorage, workDir string) {
	if envStorage.Get("KOOL_WORKSPACE_PROVIDER") != "" {
		return
	}

	rootOutput, err := gitWorktreeOutput(workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return
	}
	root := strings.TrimSpace(string(rootOutput))

	gitDirOutput, err := gitWorktreeOutput(workDir, "rev-parse", "--git-dir")
	if err != nil {
		setGitWorktreeError(envStorage, err)
		return
	}
	commonDirOutput, err := gitWorktreeOutput(workDir, "rev-parse", "--git-common-dir")
	if err != nil {
		setGitWorktreeError(envStorage, err)
		return
	}

	gitDir := absoluteGitPath(workDir, strings.TrimSpace(string(gitDirOutput)))
	commonDir := absoluteGitPath(workDir, strings.TrimSpace(string(commonDirOutput)))
	if gitDir == commonDir {
		return
	}

	listOutput, err := gitWorktreeOutput(workDir, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		setGitWorktreeError(envStorage, err)
		return
	}
	source := firstGitWorktree(listOutput)
	if source == "" {
		envStorage.Set("KOOL_WORKSPACE_ERROR", "could not resolve Git worktree source")
		return
	}

	initWorkspaceContext(envStorage, workDir, source, root, "worktree", true)
}

func absoluteGitPath(workDir, path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(workDir, path)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.Clean(absolute)
}

func firstGitWorktree(output []byte) string {
	for _, field := range strings.Split(string(output), "\x00") {
		if strings.HasPrefix(field, "worktree ") {
			return strings.TrimPrefix(field, "worktree ")
		}
	}
	return ""
}

func setGitWorktreeError(envStorage EnvStorage, err error) {
	envStorage.Set("KOOL_WORKSPACE_ERROR", fmt.Sprintf("could not resolve Git worktree: %v", err))
}
