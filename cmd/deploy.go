package cmd

import (
	"fmt"
	"kool-dev/kool/api"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"kool-dev/kool/tgz"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys your application using Kool Dev",
	Run:   runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	var (
		filename     string
		deploy       *api.Deploy
		err          error
		outputWriter shell.OutputWriter
	)

	outputWriter = shell.NewOutputWriter()

	if url := environment.NewEnvStorage().Get("KOOL_API_URL"); url != "" {
		api.SetBaseURL(url)
	}

	fmt.Println("Create release file...")
	filename, err = createReleaseFile()

	if err != nil {
		outputWriter.Error(err)
		os.Exit(1)
	}

	defer func(file string) {
		var err error
		if err = os.Remove(file); err != nil {
			outputWriter.Error(fmt.Errorf("error trying to remove temporary tarball: %v", err))
		}
	}(filename)

	deploy = api.NewDeploy(filename)

	fmt.Println("Upload release file...")
	err = deploy.SendFile()

	if err != nil {
		outputWriter.Error(err)
		os.Exit(1)
	}

	fmt.Println("Going to deploy...")

	timeout := 10 * time.Minute

	if min, err := strconv.Atoi(environment.NewEnvStorage().Get("KOOL_API_TIMEOUT")); err == nil {
		timeout = time.Duration(min) * time.Minute
	}

	var finishes chan bool = make(chan bool)

	go func(deploy *api.Deploy, finishes chan bool) {
		var lastStatus string
		for {
			err = deploy.GetStatus()

			if lastStatus != deploy.Status {
				lastStatus = deploy.Status
				fmt.Println("  > deploy:", lastStatus)
			}

			if err != nil {
				finishes <- false
				outputWriter.Error(err)
				break
			}

			if deploy.IsSuccessful() {
				finishes <- true
				break
			}

			time.Sleep(time.Second * 3)
		}
	}(deploy, finishes)

	var success bool
	select {
	case success = <-finishes:
		{
			if success {
				outputWriter.Success("Deploy finished: ", deploy.GetURL())
			} else {
				outputWriter.Error(fmt.Errorf("deploy failed"))
				os.Exit(1)
			}
			break
		}

	case <-time.After(timeout):
		{
			outputWriter.Error(fmt.Errorf("timeout waiting deploy to finish"))
			break
		}
	}
}

func createReleaseFile() (filename string, err error) {
	var (
		tarball *tgz.TarGz
		cwd     string
	)

	tarball, err = tgz.NewTemp()

	if err != nil {
		return
	}

	var hasGit bool = true
	if _, err = exec.LookPath("git"); err != nil {
		hasGit = false
	}

	if _, err = os.Stat(".git"); hasGit && !os.IsNotExist(err) {
		// we are in a Git environment
		var (
			output []byte
			files  []string
		)
		// Exclude list
		// git ls-files -d // delete files
		output, err = exec.Command("git", "ls-files", "-d").CombinedOutput()
		if err != nil {
			panic(fmt.Errorf("Failed listing deleted "))
		}
		tarball.SetIgnoreList(strings.Split(string(output), "\n"))

		// Include list
		// git ls-files -c
		output, err = exec.Command("git", "ls-files", "-c").CombinedOutput()
		if err != nil {
			panic(fmt.Errorf("Failed list Git cached files"))
		}
		files = append(files, strings.Split(string(output), "\n")...)
		// git ls-files -o --exclude-standard
		output, err = exec.Command("git", "ls-files", "-o", "--exclude-standard").CombinedOutput()
		if err != nil {
			panic(fmt.Errorf("Failed list Git untracked non-ignored files"))
		}
		files = append(files, strings.Split(string(output), "\n")...)

		filename, err = tarball.CompressFiles(files)
	} else {
		fmt.Println("Fallback to tarball full current working directory...")
		cwd, _ = os.Getwd()
		filename, err = tarball.CompressFolder(cwd)
	}

	return
}
