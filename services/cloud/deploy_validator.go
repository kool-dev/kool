package cloud

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// KoolDeployFile holds the config (meta + cloud) for a deploy
type DeployConfig struct {
	// the information about the folder/file parsed
	Meta struct {
		WorkingDir string
		Filename   string
	}

	Cloud *CloudConfig
}

// CloudConfig is the configuration for a deploy parsed from kool.cloud.yml
type CloudConfig struct {
	// version of the Kool.dev Cloud config file
	Version string `yaml:"version"`

	// services is a map of services to deploy
	Services map[string]*DeployConfigService `yaml:"services"`
}

// DeployConfigService is the configuration for a service to deploy
type DeployConfigService struct {
	Image  *string      `yaml:"image,omitempty"`
	Build  *interface{} `yaml:"build,omitempty"`
	Expose *int         `yaml:"expose,omitempty"`

	Public      interface{} `yaml:"public,omitempty"`
	Environment interface{} `yaml:"environment"`
}

// DeployConfigBuild is the configuration for a service to be built
type DeployConfigBuild struct {
	Context    *string                 `yaml:"context,omitempty"`
	Dockerfile *string                 `yaml:"dockerfile,omitempty"`
	Args       *map[string]interface{} `yaml:"args,omitempty"`
}

func ParseCloudConfig(workingDir string, koolDeployFile string) (deployConfig *DeployConfig, err error) {
	var (
		path    string = filepath.Join(workingDir, koolDeployFile)
		content []byte
	)

	if _, err = os.Stat(path); err != nil {
		// fallback to legacy naming convetion
		path = filepath.Join(workingDir, "kool.deploy.yml")

		if _, err = os.Stat(path); err != nil {
			err = fmt.Errorf("could not find required file '%s' on current working directory: %v", koolDeployFile, err)
			return
		}

		koolDeployFile = "kool.deploy.yml"
		return
	}

	if content, err = os.ReadFile(path); err != nil {
		return
	}

	deployConfig = &DeployConfig{
		Cloud: &CloudConfig{},
	}

	deployConfig.Meta.WorkingDir = workingDir
	deployConfig.Meta.Filename = koolDeployFile

	err = yaml.Unmarshal(content, deployConfig.Cloud)
	return
}

func ValidateConfig(deployConfig *DeployConfig) (err error) {
	for service, config := range deployConfig.Cloud.Services {
		// validates build file exists if defined
		if config.Build != nil {
			// we got something to build! check if it's a string
			if buildStr, buildIsString := (*config.Build).(string); buildIsString {
				// check if file exists
				var buildStat os.FileInfo
				if buildStat, err = os.Stat(filepath.Join(deployConfig.Meta.WorkingDir, buildStr)); err != nil {
					err = fmt.Errorf("service '%s' defines a build directory ('%s') that does not exist (%v)", service, buildStr, err)
					return
				}

				if !buildStat.IsDir() {
					err = fmt.Errorf("service '%s' build entry '%s' is not a directory (check v3 upgrade guide)", service, buildStr)
					return
				}
			}
		}

		// validates only one service can be public, and it must define a port
		if config.Public != nil {
			// being public, it must define the `port` entry as well
			if config.Expose == nil {
				err = fmt.Errorf("service (%s) is public, but it does not define the `export` entry", service)
				return
			}
		}
	}

	return
}
