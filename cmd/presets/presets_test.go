package presets

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/parser"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestPresetsKoolFile(t *testing.T) {
	configs := GetConfigs()

	for preset, files := range GetAll() {
		var (
			parsed *parser.KoolYaml
			kool   []byte
		)

		if _, hasKool := files["kool.yml"]; !hasKool {
			config := new(PresetConfig)

			if err := yaml.Unmarshal([]byte(configs[preset]), &config); err != nil {
				t.Fatal(err)
			}

			if _, hasKoolQuestions := config.Questions["kool"]; !hasKoolQuestions {
				t.Errorf("preset %s must have either 'kool.yml' on preset files or 'kool' questions on 'preset-config.yml' to create a new 'kool.yml'", preset)
			}

			continue
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
