package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"kool-dev/kool/api"
	"kool-dev/kool/tgz"
	"os"

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
		filename, deployID string
		deploy             *api.Deploy
		err                error
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
		if err := os.Remove(file); err != nil {
			fmt.Println("error trying to remove temporary tarball:", err)
		}
	}(filename)

	deploy = api.NewDeploy(filename)

	deployID, err = deploy.SendFile()

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
	// api.Deploy()
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
	if _, err := os.Stat(".gitignore"); err != os.ErrNotExist {
		var (
			file       *os.File
			ignoreBlob []byte
			gitIgnore  [][]byte
		)

		file, err = os.Open(".gitignore")
		ignoreBlob, err = ioutil.ReadAll(file)
		gitIgnore = bytes.Split(ignoreBlob, []byte("\n"))
		gitIgnore = append(gitIgnore, []byte("/.git")) // ignoring .git itself
		tarball.SetIgnoreList(gitIgnore)
		gitIgnore = nil
	}

	cwd, _ = os.Getwd()
	filename, err = tarball.Compress(cwd)

	return
}
