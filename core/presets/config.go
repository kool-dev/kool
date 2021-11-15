package presets

import (
	"kool-dev/kool/core/automate"
)

// PresetConfig preset config
type PresetConfig struct {
	Name   string                `yaml:"name"`
	Tags   []string              `yaml:"tags"`
	Create []*automate.ActionSet `yaml:"create"`
	Preset []*automate.ActionSet `yaml:"preset"`

	presetID string
}

func (c *PresetConfig) HasTag(tag string) bool {
	for _, t := range c.Tags {
		if tag == t {
			return true
		}
	}
	return false
}
