package commands

import (
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/utils"
	"kool-dev/kool/services/cloud"
	"kool-dev/kool/services/cloud/api"
	"kool-dev/kool/services/cloud/setup"
	"kool-dev/kool/services/tgz"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	koolDeployEnv = "kool.deploy.env"
)

// KoolCloudDeployFlags holds the flags for the kool cloud deploy command
type KoolCloudDeployFlags struct {
	// Token              string   // env: KOOL_API_TOKEN
	// DeployDomain       string   // env: KOOL_DEPLOY_DOMAIN
	Timeout            uint     // env: KOOL_API_TIMEOUT
	WwwRedirect        bool     // env: KOOL_DEPLOY_WWW_REDIRECT
	DeployDomainExtras []string // env: KOOL_DEPLOY_DOMAIN_EXTRAS

	// Cluster            string // env: KOOL_DEPLOY_CLUSTER
}

// KoolDeploy holds handlers and functions for using Deploy API
type KoolDeploy struct {
	DefaultKoolService

	cloud       *Cloud
	setupParser setup.CloudSetupParser
	flags       *KoolCloudDeployFlags
	env         environment.EnvStorage
	cloudConfig *cloud.DeployConfig
}

// NewDeployCommand initializes new kool deploy Cobra command
func NewDeployCommand(deploy *KoolDeploy) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a local application to a Kool.dev Cloud environment",
		RunE:  DefaultCommandRunFunction(deploy),
		Args:  cobra.NoArgs,

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().StringArrayVarP(&deploy.flags.DeployDomainExtras, "domain-extra", "", []string{}, "List of extra domain aliases")
	cmd.Flags().BoolVarP(&deploy.flags.WwwRedirect, "www-redirect", "", false, "Redirect www to non-www domain")
	cmd.Flags().UintVarP(&deploy.flags.Timeout, "timeout", "", 0, "Timeout in minutes for waiting the deployment to finish")

	return
}

// NewKoolDeploy creates a new pointer with default KoolDeploy service dependencies.
func NewKoolDeploy(cloud *Cloud) *KoolDeploy {
	env := environment.NewEnvStorage()
	return &KoolDeploy{
		*newDefaultKoolService(),
		cloud,
		setup.NewDefaultCloudSetupParser(env.Get("PWD")),
		&KoolCloudDeployFlags{},
		env,
		nil,
	}
}

// Execute runs the deploy logic.
func (d *KoolDeploy) Execute(args []string) (err error) {
	var (
		filename      string
		deployCreated *api.DeployCreateResponse

		deployer = cloud.NewDeployer()
	)

	d.Shell().Info("Load and validate config files...")

	if err = d.loadAndValidateConfig(); err != nil {
		return
	}

	d.Shell().Info("Create release file...")
	if filename, err = d.createReleaseFile(); err != nil {
		return
	}
	defer d.cleanupReleaseFile(filename)

	s := utils.MakeFastLoading("Creating deploy...", "Deploy created.", d.Shell().OutStream())
	if deployCreated, err = deployer.CreateDeploy(filename); err != nil {
		return
	}
	s.Stop()

	d.Shell().Info("Building images...")
	for svcName, svc := range d.cloudConfig.Cloud.Services {
		if svc.Build != nil {
			d.Shell().Info("  > Build deploy image for service: ", svcName)

			if err = cloud.BuildPushImageForDeploy(svcName, svc, deployCreated); err != nil {
				if reportErr := deployer.BuildError(deployCreated, err); reportErr != nil {
					d.Shell().Error(fmt.Errorf("error reporting build error: %v", reportErr))
				}

				return
			}

			d.Shell().Info("  > Image for service: ", svcName, " built & pushed successfully.")
		}
	}

	if deployCreated.LogsUrl != "" {
		d.Shell().Info(strings.Repeat("-", 40))
		d.Shell().Info("Logs available at: ", deployCreated.LogsUrl)
		d.Shell().Info(strings.Repeat("-", 40))
	}

	s = utils.MakeFastLoading("Start deploying...", "Deploy started.", d.Shell().OutStream())
	if _, err = deployer.StartDeploy(deployCreated); err != nil {
		return
	}
	s.Stop()

	timeout := 15 * time.Minute
	if d.flags.Timeout > 0 {
		timeout = time.Duration(d.flags.Timeout) * time.Minute
	} else if min, err := strconv.Atoi(d.env.Get("KOOL_API_TIMEOUT")); err == nil {
		timeout = time.Duration(min) * time.Minute
	}

	var finishes chan bool = make(chan bool)

	go func(deployCreated *api.DeployCreateResponse, finishes chan bool) {
		var (
			status     *api.DeployStatusResponse
			lastStatus string
			err        error
		)

		s = utils.MakeSlowLoading("Waiting deploy to finish...", "Deploy finished.", d.Shell().OutStream())

		for {
			if status, err = api.NewDeployStatus(deployCreated).Run(); err != nil {
				return
			}

			if lastStatus != status.Status {
				lastStatus = status.Status
				d.Shell().Println("  > deploy:", lastStatus)
			}

			if err != nil {
				s.Stop()
				finishes <- false
				d.Shell().Error(err)
				break
			}

			if status.Status == "success" || status.Status == "failed" {
				s.Stop()
				finishes <- status.Status == "success"
				break
			}

			time.Sleep(time.Second * 2)
		}
	}(deployCreated, finishes)

	var success bool
	select {
	case success = <-finishes:
		{
			if success {
				d.Shell().Success("Deploy finished: ", deployCreated.Deploy.Environment.Name)
				d.Shell().Success("")
				d.Shell().Success("Access your environment at: ", deployCreated.Deploy.Url)
				d.Shell().Success("")
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
	)

	tarball, err = tgz.NewTemp()

	if err != nil {
		return
	}

	var allFiles []string

	// new behavior - tarball only the required files
	var possibleKoolDeployYmlFiles []string = []string{
		"kool.deploy.yml",
		"kool.deploy.yaml",
		"kool.cloud.yml",
		"kool.cloud.yaml",
	}

	for _, file := range possibleKoolDeployYmlFiles {
		if !strings.HasPrefix(file, "/") {
			file = filepath.Join(d.env.Get("PWD"), file)
		}

		if _, err = os.Stat(file); err == nil {
			allFiles = append(allFiles, file)
		}
	}

	if len(allFiles) == 0 {
		err = fmt.Errorf("no kool.cloud.yml config files found")
		return
	}

	if cf := d.env.Get("COMPOSE_FILE"); cf != "" {
		allFiles = append(allFiles, strings.Split(cf, ":")...)
	} else {
		allFiles = append(allFiles, filepath.Join(d.env.Get("PWD"), "docker-compose.yml"))
	}

	// we need to include the .env files as well
	allFiles = append(allFiles, d.cloudConfig.GetEnvFiles()...)

	d.shell.Println("Compressing files:")
	for _, file := range allFiles {
		d.shell.Println("  -", file)
	}

	tarball.SourceDir(d.env.Get("PWD"))

	filename, err = tarball.CompressFiles(allFiles)

	if err == nil {
		d.shell.Println("Files compression done.")
	}

	return
}

func (d *KoolDeploy) loadAndValidateConfig() (err error) {
	if d.cloudConfig, err = cloud.ParseCloudConfig(d.env.Get("PWD"), setup.KoolDeployFile); err != nil {
		return
	}

	if err = cloud.ValidateConfig(d.cloudConfig); err != nil {
		return
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

func (d *KoolDeploy) cleanupReleaseFile(filename string) {
	if err := os.Remove(filename); err != nil {
		d.Shell().Error(fmt.Errorf("error trying to remove temporary tarball: %v", err))
	}
}
