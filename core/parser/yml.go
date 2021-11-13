package parser

import (
	"fmt"
	"io"
	"kool-dev/kool/core/builder"
	"os"

	"github.com/agnivade/levenshtein"
	"gopkg.in/yaml.v2"
)

// SimilarThreshold represents the minimal Levenshteindistance of two
// script names for them to be considered similarss
const SimilarThreshold int = 2

type yamlMarshalFnType func(interface{}) ([]byte, error)

// KoolYaml holds the structure for parsing the custom commands file
type KoolYaml struct {
	Scripts map[string]interface{} `yaml:"scripts"`
}

// KoolYamlParser holds logic for handling kool yaml
type KoolYamlParser interface {
	Parse(string) error
	HasScript(string) bool
	ParseCommands(string) ([]builder.Command, error)
	SetScript(string, []string)
	String() (string, error)
}

var yamlMarshalFn yamlMarshalFnType = yaml.Marshal

// ParseKoolYaml decodes the target kool.yml into its
// the expected KoolYaml representation.
func ParseKoolYaml(filePath string) (parsed *KoolYaml, err error) {
	var (
		file *os.File
		raw  []byte
	)

	if file, err = os.OpenFile(filePath, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}

	defer file.Close()

	if raw, err = io.ReadAll(file); err != nil {
		return
	}

	parsed = new(KoolYaml)
	err = yaml.Unmarshal(raw, parsed)

	return
}

// Parse decodes the target kool.yml
func (y *KoolYaml) Parse(filePath string) (err error) {
	var parsed *KoolYaml
	if parsed, err = ParseKoolYaml(filePath); err != nil {
		return
	}

	y.Scripts = parsed.Scripts
	return
}

// HasScript tells if the given script exists on this parsed YAML.
func (y *KoolYaml) HasScript(script string) (has bool) {
	if y.Scripts != nil {
		_, has = y.Scripts[script]
	}
	return
}

// GetSimilars checks for scripts with similar name.
func (y *KoolYaml) GetSimilars(script string) (has bool, similars []string) {
	var name string
	for name = range y.Scripts {
		if levenshtein.ComputeDistance(name, script) < SimilarThreshold {
			has = true
			similars = append(similars, name)
		}
	}
	return
}

// ParseCommands parsed the given script from kool.yml file into a list
// of commands parsed.
func (y *KoolYaml) ParseCommands(script string) (commands []builder.Command, err error) {
	var (
		isSingle bool
		isList   bool
		line     string
		lines    []interface{}
		command  *builder.DefaultCommand
	)

	if line, isSingle = y.Scripts[script].(string); isSingle {
		if command, err = builder.ParseCommand(line); err != nil {
			return
		}

		commands = append(commands, command)
	} else if lines, isList = y.Scripts[script].([]interface{}); isList {
		for _, i := range lines {
			if command, err = builder.ParseCommand(i.(string)); err != nil {
				return
			}

			commands = append(commands, command)
		}
	} else {
		err = fmt.Errorf("failed parsing script '%s': expected string or array of strings", script)
	}
	return
}

// SetScript set script into kool yaml
func (y *KoolYaml) SetScript(key string, commands []string) {
	if len(commands) == 0 {
		return
	}

	if y.Scripts == nil {
		y.Scripts = make(map[string]interface{})
	}

	if len(commands) == 1 {
		y.Scripts[key] = commands[0]
	} else {
		var scripts []interface{}

		for _, c := range commands {
			scripts = append(scripts, c)
		}
		y.Scripts[key] = scripts
	}
}

// String returns docker-compose as string
func (y *KoolYaml) String() (content string, err error) {
	var parsedBytes []byte

	if parsedBytes, err = yamlMarshalFn(y); err != nil {
		return
	}

	content = string(parsedBytes)
	return
}
