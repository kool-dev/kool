package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"kool-dev/kool/api"
	"kool-dev/kool/tgz"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// DeployFlags holds the flags for the start command
type DeployFlags struct {
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys your application usin Kool Dev",
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

	err = deploy.SendFile()

	if err != nil {
		execError("", err)
		os.Exit(1)
	}

	var finishes chan bool = make(chan bool)

	go func(deploy *api.Deploy, finishes chan bool) {
		for {
			err = deploy.GetStatus()

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

	// ignoring .gitignore
	if _, err = os.Stat(".gitignore"); err != os.ErrNotExist {
		var (
			file       *os.File
			ignoreBlob []byte
			gitIgnore  [][]byte
		)

		if file, err = os.Open(".gitignore"); err != nil {
			return
		}
		if ignoreBlob, err = ioutil.ReadAll(file); err != nil {
			return
		}
		gitIgnore = bytes.Split(ignoreBlob, []byte("\n"))
		gitIgnore = append(gitIgnore, []byte("/.git")) // ignoring .git itself
		tarball.SetIgnoreList(gitIgnore)
		gitIgnore = nil
	}

	cwd, _ = os.Getwd()
	filename, err = tarball.Compress(cwd)

	return
}
