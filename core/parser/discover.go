package parser

import (
	"os"
	"path/filepath"
)

// koolYamlFileNames lists the accepted kool config file names, in lookup order.
var koolYamlFileNames = []string{"kool.yml", "kool.yaml"}

// FindKoolYaml returns the path of the kool config file within the given
// directory, preferring kool.yml over kool.yaml. It returns ErrKoolYmlNotFound
// when neither file exists.
func FindKoolYaml(dir string) (file string, err error) {
	for _, name := range koolYamlFileNames {
		candidate := filepath.Join(dir, name)
		if _, err = os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", ErrKoolYmlNotFound
}

// LoadKoolYaml finds and decodes the kool config file within the given
// directory. It returns ErrKoolYmlNotFound when no config file exists; any
// other error means a file was found but could not be decoded, so callers can
// tell a malformed config apart from a missing one.
func LoadKoolYaml(dir string) (parsed *KoolYaml, err error) {
	var file string

	if file, err = FindKoolYaml(dir); err != nil {
		return
	}

	return ParseKoolYaml(file)
}
