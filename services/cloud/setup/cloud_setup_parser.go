package setup

import (
	"os"
	"path/filepath"
)

const KoolDeployFile string = "kool.cloud.yml"

type CloudSetupParser interface {
	HasDeployConfig() bool
	ConfigFilePath() string
}

type DefaultCloudSetupParser struct {
	pwd string
}

func NewDefaultCloudSetupParser(pwd string) *DefaultCloudSetupParser {
	return &DefaultCloudSetupParser{pwd}
}

func (s *DefaultCloudSetupParser) ConfigFilePath() string {
	return filepath.Join(s.pwd, KoolDeployFile)
}

func (s *DefaultCloudSetupParser) HasDeployConfig() (has bool) {
	_, err := os.Stat(s.ConfigFilePath())

	return err == nil
}
