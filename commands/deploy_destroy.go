package commands

import (
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/services/cloud/api"

	"github.com/spf13/cobra"
)

// KoolDeployDestroy holds handlers and functions for using Deploy API
type KoolDeployDestroy struct {
	DefaultKoolService

	env        environment.EnvStorage
	apiDestroy api.DestroyCall
}

// NewDeployDestroyCommand initializes new kool deploy Cobra command
func NewDeployDestroyCommand(destroy *KoolDeployDestroy) *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy an environment deployed to Kool Cloud",
		Args:  cobra.NoArgs,
		RunE:  DefaultCommandRunFunction(destroy),

		DisableFlagsInUseLine: true,
	}
}

// NewKoolDeployDestroy creates a new pointer with default KoolDeployDestroy service dependencies
func NewKoolDeployDestroy() *KoolDeployDestroy {
	return &KoolDeployDestroy{
		*newDefaultKoolService(),
		environment.NewEnvStorage(),
		api.NewDefaultDestroyCall(),
	}
}

// Execute runs the deploy destroy logic - integrating with Deploy API
func (d *KoolDeployDestroy) Execute(args []string) (err error) {
	var (
		domain string
		resp   *api.DestroyResponse
	)

	if url := d.env.Get("KOOL_API_URL"); url != "" {
		api.SetBaseURL(url)
	}

	if domain = d.env.Get("KOOL_DEPLOY_DOMAIN"); domain == "" {
		err = fmt.Errorf("missing deploy domain (env KOOL_DEPLOY_DOMAIN)")
		return
	}

	d.apiDestroy.Query().Set("domain", domain)

	if resp, err = d.apiDestroy.Call(); err != nil {
		return
	}

	d.Success(fmt.Sprintf("Environment (ID: %d) scheduled for deleting.", resp.Environment.ID))

	return
}
