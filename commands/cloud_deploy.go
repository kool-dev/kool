package commands

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/services/cloud"
	"kool-dev/kool/services/cloud/api"
	"kool-dev/kool/services/tgz"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	koolDeployEnv  = "kool.deploy.env"
	koolDeployFile = "kool.cloud.yml"
)

// KoolCloudDeployFlags holds the flags for the kool cloud deploy command
type KoolCloudDeployFlags struct {
	Token              string   // env: KOOL_API_TOKEN
	Timeout            uint     // env: KOOL_API_TIMEOUT
	WwwRedirect        bool     // env: KOOL_DEPLOY_WWW_REDIRECT
	DeployDomain       string   // env: KOOL_DEPLOY_DOMAIN
	DeployDomainExtras []string // env: KOOL_DEPLOY_DOMAIN_EXTRAS

	// Cluster            string // env: KOOL_DEPLOY_CLUSTER
	// env: KOOL_API_URL
}

// KoolDeploy holds handlers and functions for using Deploy API
type KoolDeploy struct {
	DefaultKoolService

	flags *KoolCloudDeployFlags
	env   environment.EnvStorage
	git   builder.Command
}

// NewDeployCommand initializes new kool deploy Cobra command
func NewDeployCommand(deploy *KoolDeploy) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a local application to a Kool Cloud environment",
		RunE:  DefaultCommandRunFunction(deploy),
		Args:  cobra.NoArgs,

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().StringVarP(&deploy.flags.Token, "token", "", "", "Token to authenticate with Kool Cloud API")
	cmd.Flags().StringVarP(&deploy.flags.DeployDomain, "domain", "", "", "Environment domain name to deploy to")
	cmd.Flags().UintVarP(&deploy.flags.Timeout, "timeout", "", 0, "Timeout in minutes for waiting the deployment to finish")
	cmd.Flags().StringArrayVarP(&deploy.flags.DeployDomainExtras, "domain-extra", "", []string{}, "List of extra domain aliases")
	cmd.Flags().BoolVarP(&deploy.flags.WwwRedirect, "www-redirect", "", false, "Redirect www to non-www domain")

	return
}

// NewKoolDeploy creates a new pointer with default KoolDeploy service
// dependencies.
func NewKoolDeploy() *KoolDeploy {
	return &KoolDeploy{
		*newDefaultKoolService(),
		&KoolCloudDeployFlags{},
		environment.NewEnvStorage(),
		builder.NewCommand("git"),
	}
}

// Execute runs the deploy logic.
func (d *KoolDeploy) Execute(args []string) (err error) {
	var (
		filename string
		deploy   *api.Deploy
	)

	if err = d.validate(); err != nil {
		return
	}

	if url := d.env.Get("KOOL_API_URL"); url != "" {
		api.SetBaseURL(url)
	}

	d.Shell().Info("Create release file...")
	if filename, err = d.createReleaseFile(); err != nil {
		return
	}

	defer func(file string) {
		var err error
		if err = os.Remove(file); err != nil {
			d.Shell().Error(fmt.Errorf("error trying to remove temporary tarball: %v", err))
		}
	}(filename)

	deploy = api.NewDeploy(filename)

	d.Shell().Info("Upload release file...")
	if err = deploy.SendFile(); err != nil {
		return
	}

	d.Shell().Println("Going to deploy...")

	timeout := 15 * time.Minute

	if d.flags.Timeout > 0 {
		timeout = time.Duration(d.flags.Timeout) * time.Minute
	} else if min, err := strconv.Atoi(d.env.Get("KOOL_API_TIMEOUT")); err == nil {
		timeout = time.Duration(min) * time.Minute
	}

	var finishes chan bool = make(chan bool)

	go func(deploy *api.Deploy, finishes chan bool) {
		var (
			lastStatus string
			err        error
		)

		for {
			err = deploy.FetchLatestStatus()

			if lastStatus != deploy.Status.Status {
				lastStatus = deploy.Status.Status
				d.Shell().Println("  > deploy:", lastStatus)
			}

			if err != nil {
				finishes <- false
				d.Shell().Error(err)
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
				d.Shell().Success("Deploy finished: ", deploy.GetURL())
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
	if errGit := d.Shell().LookPath(d.git); errGit != nil {
		hasGit = false
	}

	if _, errGit := os.Stat(".git"); !hasGit || os.IsNotExist(errGit) {
		// not a GIT repo/environment!
		d.Shell().Println("Fallback to tarball full current working directory...")
		cwd, _ = os.Getwd()
		filename, err = tarball.CompressFolder(cwd)
		return
	}

	// we are in a GIT environment!
	var (
		files, allFiles []string

		gitListingFilesFlags = [][]string{
			// Include commited files - git ls-files -c
			{"-c"},

			// Untracked files - git ls-files -o --exclude-standard
			{"-o", "--exclude-standard"},
		}
	)

	// Exclude list - git ls-files -d
	if files, err = d.parseFilesListFromGIT([]string{"-d"}); err != nil {
		return
	}
	tarball.SetIgnoreList(files)

	for _, lsArgs := range gitListingFilesFlags {
		if files, err = d.parseFilesListFromGIT(lsArgs); err != nil {
			return
		}

		allFiles = append(allFiles, files...)
	}
	filename, err = tarball.CompressFiles(d.handleDeployEnv(allFiles))
	return
}

func (d *KoolDeploy) parseFilesListFromGIT(args []string) (files []string, err error) {
	var (
		output, file string
	)

	output, err = d.Shell().Exec(d.git, append([]string{"ls-files", "-z"}, args...)...)
	if err != nil {
		err = fmt.Errorf("failed listing GIT files: %s", err.Error())
		return
	}

	// -z parameter returns the utf-8 file names separated by 0 bytes
	for _, file = range strings.Split(output, string(rune(0x00))) {
		if file == "" {
			continue
		}

		files = append(files, file)
	}
	return
}

// handleDeployEnv tackles a special case on kool.deploy.env file.
// This file can or cannot be versioned (good practice not to, since
// it may include sensitive data). In the case of it being ignored
// from GIT, we still are required to send it - it is required for
// the Deploy API.
func (d *KoolDeploy) handleDeployEnv(files []string) []string {
	path := filepath.Join(d.env.Get("PWD"), koolDeployEnv)
	if _, envErr := os.Stat(path); os.IsNotExist(envErr) {
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

func (d *KoolDeploy) validate() (err error) {
	if err = cloud.ValidateKoolDeployFile(d.env.Get("PWD"), koolDeployFile); err != nil {
		return
	}

	// if no domain is set, we try to get it from the environment
	if d.flags.DeployDomain == "" && d.env.Get("KOOL_DEPLOY_DOMAIN") == "" {
		err = fmt.Errorf("Missing deploy domain. Please set it via --domain or KOOL_DEPLOY_DOMAIN environment variable.")
		return
	} else if d.flags.DeployDomain != "" {
		// shares the flag via environment variable
		d.env.Set("KOOL_DEPLOY_DOMAIN", d.flags.DeployDomain)
	}

	// if no token is set, we try to get it from the environment
	if d.flags.Token == "" && d.env.Get("KOOL_API_TOKEN") == "" {
		err = fmt.Errorf("Missing Kool Cloud API token. Please set it via --token or KOOL_API_TOKEN environment variable.")
		return
	} else if d.flags.Token != "" {
		d.env.Set("KOOL_API_TOKEN", d.flags.Token)
	}

	// share the www-redirection flag via environment variable
	if d.flags.WwwRedirect {
		d.env.Set("KOOL_DEPLOY_WWW_REDIRECT", "true")
	}

	// share the domain extras via environment variable
	if len(d.flags.DeployDomainExtras) > 0 {
		d.env.Set("KOOL_DEPLOY_DOMAIN_EXTRAS", strings.Join(d.flags.DeployDomainExtras, ","))
	}

	// TODO: make a call to the cloud API to validate the config
	// - validate the token is valid
	// - validate the domain is valid / the token gives access to it
	// - validate the domain extras are valid / the token gives access to them

	return
}
