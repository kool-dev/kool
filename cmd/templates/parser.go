package templates

import (
	"strings"

	"gopkg.in/yaml.v2"
)

type yamlUnmarshalFnType func([]byte, interface{}) error

// TemplateFile holds preset template file data
type TemplateFile struct {
	Services yaml.MapSlice          `yaml:"services,omitempty"`
	Volumes  yaml.MapSlice          `yaml:"volumes,omitempty"`
	Scripts  map[string]interface{} `yaml:"scripts,omitempty"`
}

// DefaultParser holds data for template handling logic
type DefaultParser struct {
	template *TemplateFile
}

// Parser holds logic for handling templates
type Parser interface {
	Parse(string) error
	GetServices() yaml.MapSlice
	GetVolumes() yaml.MapSlice
	GetScripts() map[string][]string
}

var (
	yamlUnmarshalFn yamlUnmarshalFnType = yaml.Unmarshal
)

// NewParser creates new template parser
func NewParser() Parser {
	return &DefaultParser{}
}

// Parse parse template
func (t *DefaultParser) Parse(template string) (err error) {
	t.template = new(TemplateFile)
	err = yamlUnmarshalFn([]byte(template), &t.template)
	return
}

// GetServices get template services
func (t *DefaultParser) GetServices() yaml.MapSlice {
	return t.template.Services
}

// GetVolumes get template volumes
func (t *DefaultParser) GetVolumes() yaml.MapSlice {
	return t.template.Volumes
}

// GetScripts get template scripts
func (t *DefaultParser) GetScripts() (scriptsMap map[string][]string) {
	scriptsMap = make(map[string][]string, len(t.template.Scripts))

	for key, value := range t.template.Scripts {
		var scripts []string
		if line, isSingle := value.(string); isSingle {
			scripts = append(scripts, strings.TrimSpace(line))
		} else if lines, isList := value.([]interface{}); isList {
			for _, line := range lines {
				scripts = append(scripts, strings.TrimSpace(line.(string)))
			}
		}

		scriptsMap[key] = scripts
	}

	return
}
