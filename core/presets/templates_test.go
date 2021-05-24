package presets

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestParsedTemplates(t *testing.T) {
	for service, templates := range GetTemplates() {
		for templateKey, template := range templates {
			var yamlData yaml.MapSlice

			if !strings.Contains(templateKey, ".yml") {
				t.Error("template is not a yml file")
			}

			if err := yaml.Unmarshal([]byte(template), &yamlData); err != nil {
				t.Errorf("%s/%s is an invalid template", service, templateKey)
			}
		}
	}
}
