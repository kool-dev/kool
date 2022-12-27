package cloud

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type DeployConfig struct {
	// services is a map of services to deploy
	Services map[string]*DeployConfigService `yaml:"services"`
}

// DeployConfigService is the configuration for a service to deploy
type DeployConfigService struct {
	Image  *string `yaml:"image"`
	Build  *string `yaml:"build"`
	Port   *string `yaml:"port"`
	Public []struct {
		Port *string `yaml:"port"`
		Path *string `yaml:"path"`
	} `yaml:"public"`
}

func ValidateKoolDeployFile(workingDir string, koolDeployFile string) (err error) {
	var (
		path    string
		content []byte

		deployConfig *DeployConfig = &DeployConfig{}
	)

	path = filepath.Join(workingDir, koolDeployFile)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = fmt.Errorf("could not find required file (%s) on current working directory", koolDeployFile)
		return
	} else if err != nil {
		return
	}

	if content, err = os.ReadFile(path); err != nil {
		return
	}

	if err = yaml.Unmarshal(content, deployConfig); err != nil {
		return
	}

	var gotPublicService = false
	for service, config := range deployConfig.Services {
		// validates build file exists if defined
		if config.Build != nil {
			// we got something to build! check that file
			if _, err = os.Stat(filepath.Join(workingDir, *config.Build)); os.IsNotExist(err) {
				err = fmt.Errorf("service (%s) defines a build file (%s) that does not exist", service, *config.Build)
				return
			} else if err != nil {
				return
			}
		}

		// validates only one service can be public, and it must define a port
		if config.Public != nil {
			// this is a public service
			if gotPublicService {
				// we can have only the one!
				err = fmt.Errorf("service (%s) is public, but another service is already public", service)
				return
			}

			gotPublicService = true

			// being public, it must define the `port` entry as well
			if config.Port == nil {
				err = fmt.Errorf("service (%s) is public, but it does not define the `port` entry", service)
				return
			}
		}
	}

	return
}
