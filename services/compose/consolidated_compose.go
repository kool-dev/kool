package compose

import (
	"fmt"
	"kool-dev/kool/services/yamler"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
	yaml3 "gopkg.in/yaml.v3"

	"github.com/compose-spec/compose-go/template"
)

type DockerComposeConfig struct {
	Services map[string]*struct {
		Image   *interface{} `yaml:"image"`
		Build   *interface{} `yaml:"build"`
		Ports   []string     `yaml:"ports"`
		Volumes []string     `yaml:"volumes"`
	} `yaml:"services"`
}

func ParseConsolidatedDockerComposeConfig(workingDir string) (dockerComposeConfig *DockerComposeConfig, err error) {
	var (
		composerFiles []string
		content       []byte

		consolidatedDockerCompose *yaml3.Node = &yaml3.Node{}

		merger *yamler.DefaultMerger = &yamler.DefaultMerger{}
	)

	dockerComposeConfig = &DockerComposeConfig{}

	if composerFiles, err = getDockerComposeFiles(workingDir); err != nil {
		return
	}

	for _, file := range composerFiles {
		if content, err = os.ReadFile(file); err != nil {
			return
		}

		node := &yaml3.Node{}
		if err = yaml3.Unmarshal(content, node); err != nil {
			return
		}

		if err = merger.Merge(node, consolidatedDockerCompose); err != nil {
			return
		}
	}

	if content, err = yaml3.Marshal(consolidatedDockerCompose); err != nil {
		return
	}

	contentString := string(content)
	if contentString, err = template.Substitute(contentString, os.LookupEnv); err != nil {
		return
	}

	err = yaml.Unmarshal([]byte(contentString), dockerComposeConfig)
	return
}

// getDockerComposeFiles returns a list of docker-compose files (absolute paths)
func getDockerComposeFiles(workingDir string) (files []string, err error) {
	composerFile := os.Getenv("COMPOSE_FILE")

	if composerFile != "" {
		files = strings.Split(composerFile, ":")

		for i := range files {
			file := filepath.Join(workingDir, files[i])
			if _, err = os.Stat(file); os.IsNotExist(err) {
				err = fmt.Errorf("could not find required file (%s) on current working directory (referenced by COMPOSE_FILE)", file)
				return
			} else if err != nil {
				return
			}
		}

		return
	}

	// fallback default file
	file := filepath.Join(workingDir, "docker-compose.yml")

	if _, err = os.Stat(file); os.IsNotExist(err) {
		err = fmt.Errorf("could not find required file 'docker-compose.yml' on current working directory")
		return
	} else if err == nil {
		files = append(files, file)
	}

	return
}
