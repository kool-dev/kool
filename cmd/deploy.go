package cmd

import (
	"fmt"
	"kool-dev/kool/api"
	"kool-dev/kool/tgz"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// DeployFlags holds the flags for the start command
type DeployFlags struct {
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys your application using Kool Dev",
	Run:   runDeploy,
}

var deployFlags = &DeployFlags{}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	var (
		filename string
		deploy   *api.Deploy
		err      error
	)

	if url := os.Getenv("KOOL_API_URL"); url != "" {
		api.SetBaseURL(url)
	}

	fmt.Println("Create release file...")
	filename, err = createReleaseFile()

	if err != nil {
		execError("", err)
		os.Exit(1)
	}

	defer func(file string) {
		var err error
		if err = os.Remove(file); err != nil {
			fmt.Println("error trying to remove temporary tarball:", err)
		}
	}(filename)

	deploy = api.NewDeploy(filename)

	fmt.Println("Upload release file...")
	err = deploy.SendFile()

	if err != nil {
		execError("", err)
		os.Exit(1)
	}

	fmt.Println("Going to deploy...")

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
				execError("", err)
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
				fmt.Println("Deploy finished:", deploy.GetURL())
			} else {
				fmt.Println("Deploy failed.")
				os.Exit(1)
			}
			break
		}

	case <-time.After(time.Minute * 10):
		{
			fmt.Println("timeout waiting deploy to finish")
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
