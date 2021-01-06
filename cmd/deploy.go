package cmd

import (
	"fmt"
	"kool-dev/kool/api"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"kool-dev/kool/tgz"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	koolDeployEnv = "kool.deploy.env"
)

// KoolDeploy holds handlers and functions for using Deploy API
type KoolDeploy struct {
	DefaultKoolService

	envStorage environment.EnvStorage
	git        builder.Command
}

// NewDeployCommand initializes new kool deploy Cobra command
func NewDeployCommand(deploy *KoolDeploy) *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Deploys your application using Kool Dev",
		Run:   DefaultCommandRunFunction(deploy),
	}
}

// NewKoolDeploy creates a new pointer with default KoolDeploy service
// dependencies.
func NewKoolDeploy() *KoolDeploy {
	defaultKoolService := newDefaultKoolService()
	return &KoolDeploy{
		*defaultKoolService,
		environment.NewEnvStorage(),

		builder.NewCommand("git"),
	}
}

func init() {
	rootCmd.AddCommand(NewDeployCommand(NewKoolDeploy()))
}

// Execute runs the deploy logic.
func (d *KoolDeploy) Execute(args []string) (err error) {
	var (
		filename string
		deploy   *api.Deploy
	)

	if url := d.envStorage.Get("KOOL_API_URL"); url != "" {
		api.SetBaseURL(url)
	}

	d.Println("Create release file...")
	filename, err = d.createReleaseFile()

	if err != nil {
		return
	}

	defer func(file string) {
		var err error
		if err = os.Remove(file); err != nil {
			d.Error(fmt.Errorf("error trying to remove temporary tarball: %v", err))
		}
	}(filename)

	deploy = api.NewDeploy(filename)

	d.Println("Upload release file...")
	err = deploy.SendFile()

	if err != nil {
		return
	}

	d.Println("Going to deploy...")

	timeout := 10 * time.Minute

	if min, err := strconv.Atoi(d.envStorage.Get("KOOL_API_TIMEOUT")); err == nil {
		timeout = time.Duration(min) * time.Minute
	}

	var finishes chan bool = make(chan bool)

	go func(deploy *api.Deploy, finishes chan bool) {
		var (
			lastStatus string
			err        error
		)

		for {
			err = deploy.GetStatus()

			if lastStatus != deploy.Status {
				lastStatus = deploy.Status
				d.Println("  > deploy:", lastStatus)
			}

			if err != nil {
				finishes <- false
				d.Error(err)
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
				d.Success("Deploy finished: ", deploy.GetURL())
			} else {
				err = fmt.Errorf("deploy failed")
				return
			}
			break
		}

	case <-time.After(timeout):
		{
			err = fmt.Errorf("timeout waiting deploy to finish")
			break
		}
	}

	return
}

func (d *KoolDeploy) createReleaseFile() (filename string, err error) {
	var (
		tarball *tgz.TarGz
		cwd     string
	)

	tarball, err = tgz.NewTemp()

	if err != nil {
		return
	}

	var hasGit bool = true
	if errGit := d.LookPath(d.git); errGit != nil {
		hasGit = false
	}

	if _, errGit := os.Stat(".git"); !hasGit || os.IsNotExist(errGit) {
		// not a GIT repo/environment!
		d.Println("Fallback to tarball full current working directory...")
		cwd, _ = os.Getwd()
		filename, err = tarball.CompressFolder(cwd)
		return
	}

	// we are in a GIT environment!
	var (
		output string
		files  []string
	)
	// Exclude list - git ls-files -d
	output, err = d.Exec(d.git, "ls-files", "-d")
	if err != nil {
		d.Error(err)
		err = fmt.Errorf("failed listing GIT deleted files")
		return
	}
	tarball.SetIgnoreList(strings.Split(output, "\n"))

	// Include list - git ls-files -c
	output, err = d.Exec(d.git, "ls-files", "-c")
	if err != nil {
		d.Error(err)
		err = fmt.Errorf("failed list GIT cached files")
		return
	}
	files = append(files, strings.Split(output, "\n")...)
	// git ls-files -o --exclude-standard
	output, err = d.Exec(d.git, "ls-files", "-o", "--exclude-standard")
	if err != nil {
		d.Error(err)
		err = fmt.Errorf("failed list Git untracked non-ignored files")
		return
	}
	files = append(files, strings.Split(output, "\n")...)

	filename, err = tarball.CompressFiles(d.handleDeployEnv(files))
	return
}

// handleDeployEnv tackles a special case on kool.deploy.env file.
// This file can or cannot be versioned (good practice not to, since
// it may include sensitive data). In the case of it being ignored
// from GIT, we still are required to send it - it is required for
// the Deploy API.
func (d *KoolDeploy) handleDeployEnv(files []string) []string {
	if _, envErr := os.Stat(koolDeployEnv); envErr == os.ErrNotExist {
		d.Error(fmt.Errorf("missing required file: %v", koolDeployEnv))
		return files
	}

	var isAlreadyIncluded bool = false
	for _, file := range files {
		if file == koolDeployEnv {
			isAlreadyIncluded = true
			break
		}
	}

	if !isAlreadyIncluded {
		files = append(files, koolDeployEnv)
	}

	return files
}
