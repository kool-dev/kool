package presets

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestParsedConfigs(t *testing.T) {
	presets := make(map[string]bool)

	for preset := range GetAll() {
		presets[preset] = false
	}

	for preset, configValue := range GetConfigs() {
		if _, presetExists := presets[preset]; !presetExists {
			t.Errorf("preset %s does not exist", preset)
		}

		presets[preset] = true

		config := new(PresetConfig)

		if err := yaml.Unmarshal([]byte(configValue), &config); err != nil {
			t.Errorf("configuration for preset %s is invalid", preset)
		}

		if config.Language == "" {
			t.Errorf("language for preset %s is not configured", preset)
		}
	}

	for preset, isConfigured := range presets {
		if !isConfigured {
			t.Errorf("preset %s does not have a configuration file 'preset-config.yml'", preset)
		}
	}
}
