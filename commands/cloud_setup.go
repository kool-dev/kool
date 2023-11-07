package commands

import (
	"bytes"
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud"
	"kool-dev/kool/services/cloud/setup"
	"kool-dev/kool/services/compose"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	yaml3 "gopkg.in/yaml.v3"
)

// KoolCloudSetup holds handlers and functions for setting up deployment configuration
type KoolCloudSetup struct {
	DefaultKoolService

	setupParser  setup.CloudSetupParser
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
	env := environment.NewEnvStorage()
	return &KoolCloudSetup{
		*newDefaultKoolService(),

		setup.NewDefaultCloudSetupParser(env.Get("PWD")),
		shell.NewPromptSelect(),
		env,
	}
}

// Execute runs the setup logic.
func (s *KoolCloudSetup) Execute(args []string) (err error) {
	var (
		composeConfig *compose.DockerComposeConfig
		serviceName   string

		deployConfig *cloud.DeployConfig = &cloud.DeployConfig{
			Version:  "1.0",
			Services: make(map[string]*cloud.DeployConfigService),
		}

		postInstructions []func()
	)

	if !s.Shell().IsTerminal() {
		err = fmt.Errorf("setup command is not available in non-interactive mode")
		return
	}

	s.Shell().Warning("Kool.dev Cloud auto-setup is an experimental feature. Make sure to review all the generated configuration files before deploying.")

	s.Shell().Info("Loading docker compose configuration...")

	if composeConfig, err = compose.ParseConsolidatedDockerComposeConfig(s.env.Get("PWD")); err != nil {
		return
	}

	s.Shell().Info("Docker compose configuration loaded. Starting interactive setup:")

	var hasPublicPort bool = false

	for serviceName = range composeConfig.Services {
		var (
			confirmed bool
			isPublic  bool = false
			answer    string

			composeService = composeConfig.Services[serviceName]
		)

		if confirmed, err = s.promptSelect.Confirm("Do you want to deploy the service container '%s'?", serviceName); err != nil {
			return
		} else if !confirmed {
			s.Shell().Warning(fmt.Sprintf("SKIP - not deploying service container '%s'", serviceName))
			continue
		}

		s.Shell().Info(fmt.Sprintf("Setting up service container '%s' for deployment", serviceName))
		deployConfig.Services[serviceName] = &cloud.DeployConfigService{
			Environment: map[string]string{},
		}

		// services needs to have either a build config or refer to a pre-built image
		if composeService.Build == nil && composeService.Image == nil {
			err = fmt.Errorf("unable to deploy service '%s': it needs to define an image or spec to build one", serviceName)
			return
		}

		// handle image/build config
		if composeService.Build != nil {
			// validate the referenced file exists
			var buildFilePath string
			if ctx, isString := (*composeService.Build).(string); isString {
				// if it's a string, that should be the build path
				buildFilePath = filepath.Join(ctx, "Dockerfile")
			} else if buildConfig, isMap := (*composeService.Build).(map[string]interface{}); isMap {
				ctx, exists := buildConfig["context"].(string)
				if !exists || ctx == "" {
					ctx = "."
				}

				if customFilename, exists := buildConfig["dockerfile"].(string); exists {
					buildFilePath = filepath.Join(ctx, customFilename)
				} else {
					buildFilePath = filepath.Join(ctx, "Dockerfile")
				}
			}

			// now just make sure we can see/have this file
			if _, buildPathErr := os.Stat(buildFilePath); os.IsNotExist(buildPathErr) {
				err = fmt.Errorf("build config error: service '%s' points to non-existing Dockerfile '%s'", serviceName, buildFilePath)
				return
			}

			s.Shell().Info(fmt.Sprintf("Service container '%s' builds its image from '%s'", serviceName, buildFilePath))
		} else {
			s.Shell().Info(fmt.Sprintf("Service container '%s' uses image '%v'", serviceName, *composeService.Image))

			// no build config, so we'll need to build
			if len(composeService.Volumes) > 0 {
				if confirmed, err = s.promptSelect.Confirm("Do you want to create a new Dockerfile to build service '%s'?", serviceName); err != nil {
					return
				} else if confirmed {
					s.Shell().Info(fmt.Sprintf("Going to create Dockerfile for service '%s'", serviceName))

					// so here we should build the basic/simplest Dockerfile
					deployConfig.Services[serviceName].Build = new(string)
					*deployConfig.Services[serviceName].Build = "."

					if _, errStat := os.Stat("Dockerfile"); os.IsNotExist(errStat) {
						// we don't have a Dockerfile, let's make a basic one!
						var (
							dockerfile *os.File
							content    bytes.Buffer
						)

						if dockerfile, err = os.Create("Dockerfile"); err != nil {
							return
						}

						content.WriteString(fmt.Sprintf("FROM %s\n", (*composeService.Image).(string)))

						for _, vol := range composeService.Volumes {
							volParts := strings.Split(vol, ":")

							if !strings.HasPrefix(volParts[0], ".") && !strings.HasPrefix(volParts[0], "/") {
								s.Shell().Println(fmt.Sprintf("Skipping named volume '%s'", volParts[0]))
								continue
							}

							if confirmed, err = s.promptSelect.Confirm("Do you want to add folder '%s' onto '%s' in the Dockerfile for service '%s'?", volParts[0], volParts[1], serviceName); err != nil {
								return
							} else if confirmed {
								content.WriteString(fmt.Sprintf("\nCOPY %s %s\n", volParts[0], volParts[1]))
							}
						}

						if _, err = dockerfile.Write(content.Bytes()); err != nil {
							return
						}

						_ = dockerfile.Close()

						postInstructions = append(postInstructions, func() {
							s.Shell().Info(fmt.Sprintf("⇒ New Dockerfile was created to build service '%s' for deploy. Review and make sure it has all the required steps. ", serviceName))
						})
					}
				} else {
					postInstructions = append(postInstructions, func() {
						s.Shell().Info(fmt.Sprintf("⇒ Service '%s' uses volumes. Make sure to create the necessary Dockerfile and build it to deploy if necessary.", serviceName))
					})
				}
			}
		}

		// handle port/public config
		ports := composeService.Ports
		if len(ports) > 0 {
			s.Shell().Info(fmt.Sprintf("Service container '%s' exposes network ports", serviceName))

			potentialPorts := []string{}
			for i := range ports {
				mappedPorts := strings.Split(ports[i], ":")

				potentialPorts = append(potentialPorts, mappedPorts[len(mappedPorts)-1])
			}

			if !hasPublicPort {
				if confirmed, err = s.promptSelect.Confirm("Do you want to make service '%s' publicly accessible?", serviceName); err != nil {
					return
				} else if confirmed {
					hasPublicPort = true
					isPublic = true
				}
			}

			if len(potentialPorts) > 1 {
				if answer, err = s.promptSelect.Ask("Which port do you want to use for this service?", potentialPorts); err != nil {
					return
				}
			} else {
				answer = potentialPorts[0]
			}

			deployConfig.Services[serviceName].Port = new(int)
			*deployConfig.Services[serviceName].Port, _ = strconv.Atoi(answer)

			if isPublic {
				public := &cloud.DeployConfigPublicEntry{}
				public.Port = new(int)
				*public.Port = *deployConfig.Services[serviceName].Port

				deployConfig.Services[serviceName].Public = append(deployConfig.Services[serviceName].Public, public)
			}
		}
	}

	var yaml []byte
	if yaml, err = yaml3.Marshal(deployConfig); err != nil {
		return
	}

	if err = os.WriteFile(koolDeployFile, yaml, 0644); err != nil {
		return
	}

	s.Shell().Println("")

	for _, instruction := range postInstructions {
		instruction()
	}

	s.Shell().Println("")
	s.Shell().Println("")
	s.Shell().Success("Setup completed. Please review the generated configuration file before deploying.")
	s.Shell().Println("")

	return
}
