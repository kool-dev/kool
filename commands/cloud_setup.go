package commands

import (
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud"
	"kool-dev/kool/services/compose"
	"strings"

	"github.com/spf13/cobra"
)

// KoolCloudSetup holds handlers and functions for setting up deployment configuration
type KoolCloudSetup struct {
	DefaultKoolService

	promptSelect shell.PromptSelect
	env          environment.EnvStorage
}

// NewSetupCommand initializes new kool deploy Cobra command
func NewSetupCommand(setup *KoolCloudSetup) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Set up local configuration files for deployment",
		RunE:  DefaultCommandRunFunction(setup),
		Args:  cobra.NoArgs,

		DisableFlagsInUseLine: true,
	}
}

// NewKoolCloudSetup factories new KoolCloudSetup instance pointer
func NewKoolCloudSetup() *KoolCloudSetup {
	return &KoolCloudSetup{
		*newDefaultKoolService(),

		shell.NewPromptSelect(),
		environment.NewEnvStorage(),
	}
}

// Execute runs the setup logic.
func (s *KoolCloudSetup) Execute(args []string) (err error) {
	var (
		composeConfig *compose.DockerComposeConfig
		serviceName   string

		deployConfig *cloud.DeployConfig = &cloud.DeployConfig{}
	)

	if !s.Shell().IsTerminal() {
		err = fmt.Errorf("setup command is not available in non-interactive mode")
		return
	}

	s.Shell().Info("Loading docker compose configuration...")

	if composeConfig, err = compose.ParseConsolidatedDockerComposeConfig(s.env.Get("PWD")); err != nil {
		return
	}

	s.Shell().Info("Docker compose configuration loaded. Starting interactive setup:")

	for serviceName = range composeConfig.Services {
		var answer string

		fmt.Printf("consolidated parsed '%s': %+v", serviceName, composeConfig.Services[serviceName])

		if answer, err = s.promptSelect.Ask(fmt.Sprintf("Do you want to deploy the service container '%s'?", serviceName), []string{"Yes", "No"}); err != nil {
			return
		}

		if answer == "No" {
			s.Shell().Warning(fmt.Sprintf("Not going to deploy service container '%s'", serviceName))
			continue
		}

		s.Shell().Info(fmt.Sprintf("Setting up service container '%s' for deployment", serviceName))
		deployConfig.Services[serviceName] = &cloud.DeployConfigService{}

		// handle image/build config
		if len(composeConfig.Services[serviceName].Volumes) > 0 {
			// if we have a build, keep it
		}

		// handle port/public config
		ports := composeConfig.Services[serviceName].Ports
		if len(ports) > 0 {
			potentialPorts := []string{}
			for i := range ports {
				mappedPorts := strings.Split(ports[i], ":")

				potentialPorts = append(potentialPorts, mappedPorts[len(mappedPorts)-1])
			}

			if len(potentialPorts) > 1 {
				if answer, err = s.promptSelect.Ask("Which port do you want to make public?", potentialPorts); err != nil {
					return
				}
			} else {
				answer = potentialPorts[0]
			}

			deployConfig.Services[serviceName].Port = new(string)
			*deployConfig.Services[serviceName].Port = answer
		}
	}

	return
}
