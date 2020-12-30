package compose

import (
	"gopkg.in/yaml.v2"
)

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledParse       map[string]bool
	CalledGetServices bool
	CalledSetService  map[string]bool
	CalledGetVolumes  bool
	CalledSetVolume   map[string]bool
	CalledString      bool
	MockParseError    error
	MockStringError   error
	MockGetServices   yaml.MapSlice
	MockGetVolumes    yaml.MapSlice
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

// SetService implements fake SetService behavior
func (f *FakeParser) SetService(service string, content interface{}) {
	if f.CalledSetService == nil {
		f.CalledSetService = make(map[string]bool)
	}

	f.CalledSetService[service] = true
}

// GetVolumes implements fake GetVolumes behavior
func (f *FakeParser) GetVolumes() yaml.MapSlice {
	f.CalledGetVolumes = true
	return f.MockGetVolumes
}

// SetVolume implements fake SetVolume behavior
func (f *FakeParser) SetVolume(volume string) {
	if f.CalledSetVolume == nil {
		f.CalledSetVolume = make(map[string]bool)
	}

	f.CalledSetVolume[volume] = true
}

// String implements fake String behavior
func (f *FakeParser) String() (content string, err error) {
	f.CalledString = true
	err = f.MockStringError
	return
}
