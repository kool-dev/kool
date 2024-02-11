package cloud

import (
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/api"
)

// Deployer service handles the deployment process
type Deployer struct {
	env environment.EnvStorage
	out shell.Shell
}

// NewDeployer creates a new handler for using the
// Kool Dev API for deploying your application.
func NewDeployer() *Deployer {
	return &Deployer{
		env: environment.NewEnvStorage(),
		out: shell.NewShell(),
	}
}

// SendFile integrates with the API to send the tarball
func (d *Deployer) CreateDeploy(tarballPath string) (resp *api.DeployCreateResponse, err error) {
	var create = api.NewDeployCreate()

	create.PostFile("deploy", tarballPath, "deploy.tgz")

	create.PostField("cluster", d.env.Get("KOOL_CLOUD_CLUSTER"))
	create.PostField("domain", d.env.Get("KOOL_DEPLOY_DOMAIN"))
	create.PostField("additional_domains", d.env.Get("KOOL_DEPLOY_DOMAIN_EXTRAS"))
	create.PostField("www_redirect", d.env.Get("KOOL_DEPLOY_WWW_REDIRECT"))

	resp, err = create.Run()

	return
}

func (d *Deployer) StartDeploy(created *api.DeployCreateResponse) (started *api.DeployStartResponse, err error) {
	var start = api.NewDeployStart(created)

	started, err = start.Run()
	return
}

func (d *Deployer) BuildError(created *api.DeployCreateResponse, gotErr error) (err error) {
	var buildErr = api.NewDeployError(created, gotErr)

	err = buildErr.Run()
	return
}
