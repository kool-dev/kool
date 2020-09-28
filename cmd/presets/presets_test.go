package presets

import (
	"gopkg.in/yaml.v2"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/parser"
	"testing"
)

func TestPresetsKoolFile(t *testing.T) {
	for preset, files := range presets {
		var (
			parsed *parser.KoolYaml
			kool   []byte
		)

		if _, hasKool := files["kool.yml"]; !hasKool {
			t.Errorf("kool.yml is missing from %s preset", preset)
		}

		if kool = []byte(files["kool.yml"]); len(kool) == 0 {
			t.Errorf("failed on reading kool.yml from %s preset", preset)
		}

		parsed = new(parser.KoolYaml)

		if err := yaml.Unmarshal(kool, parsed); err != nil {
			t.Errorf("failed on parsing kool.yml from %s preset", preset)
		}

		for _, script := range parsed.Scripts {
			var (
				isSingle, isList bool
				lines            []interface{}
				line             string
			)

			if line, isSingle = script.(string); isSingle {
				if _, err := builder.ParseCommand(line); err != nil {
					t.Errorf("kool.yml line ('%s') could not be parsed as a command from %s preset", line, preset)
				}
			} else if lines, isList = script.([]interface{}); isList {
				for _, i := range lines {
					if _, err := builder.ParseCommand(i.(string)); err != nil {
						t.Errorf("kool.yml line ('%s') could not be parsed as a command from %s preset", i, preset)
					}
				}
			} else {
				t.Errorf("kool.yml content has invalid scripts from %s preset", preset)
			}
		}
	}
}
