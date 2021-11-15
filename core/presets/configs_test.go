package presets

import (
	"testing"
)

func TestPresetConfigHasTags(t *testing.T) {
	c := &PresetConfig{Tags: []string{"foo"}}

	if !c.HasTag("foo") {
		t.Errorf("should have tag 'foo'")
	} else if c.HasTag("bar") {
		t.Errorf("should NOT have tag 'bar'")
	}
}
