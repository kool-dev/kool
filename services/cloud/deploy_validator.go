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
	Image *string `yaml:"image,omitempty"`
	Build *string `yaml:"build,omitempty"`
	Port  *int    `yaml:"port,omitempty"`

	Public []*DeployConfigPublicEntry `yaml:"public,omitempty"`
}

type DeployConfigPublicEntry struct {
	Port *int    `yaml:"port"`
	Path *string `yaml:"path,omitempty"`
}

func ValidateKoolDeployFile(workingDir string, koolDeployFile string) (err error) {
	var (
		path    string
		content []byte

		deployConfig *DeployConfig = &DeployConfig{}
	)

	path = filepath.Join(workingDir, koolDeployFile)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		// temporary failback to old file name
		path = filepath.Join(workingDir, "kool.deploy.yml")

		if _, err = os.Stat(path); os.IsNotExist(err) {
			err = fmt.Errorf("could not find required file (%s) on current working directory", "kool.cloud.yml")
			return
		} else if err != nil {
			return
		}
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
			// being public, it must define the `port` entry as well
			if config.Port == nil {
				err = fmt.Errorf("service (%s) is public, but it does not define the `port` entry", service)
				return
			}
		}
	}

	return
}
