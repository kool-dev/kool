package templates

import (
	"gopkg.in/yaml.v2"
)

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledParse       map[string]bool
	CalledGetServices bool
	CalledGetVolumes  bool
	CalledGetScripts  bool
	MockParseError    error
	MockGetServices   yaml.MapSlice
	MockGetVolumes    yaml.MapSlice
	MockGetScripts    map[string][]string
}

// Parse implements fake Parse behavior
func (f *FakeParser) Parse(content string) (err error) {
	if f.CalledParse == nil {
		f.CalledParse = make(map[string]bool)
	}

	f.CalledParse[content] = true
	err = f.MockParseError
	return
}

// GetServices implements fake GetServices behavior
func (f *FakeParser) GetServices() yaml.MapSlice {
	f.CalledGetServices = true
	return f.MockGetServices
}

// GetVolumes implements fake GetVolumes behavior
func (f *FakeParser) GetVolumes() yaml.MapSlice {
	f.CalledGetVolumes = true
	return f.MockGetVolumes
}

// GetScripts implements fake GetScripts behavior
func (f *FakeParser) GetScripts() map[string][]string {
	if f.MockGetScripts == nil {
		f.MockGetScripts = make(map[string][]string)
	}

	return f.MockGetScripts
}
